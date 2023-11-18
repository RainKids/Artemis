package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"strings"
)

type Mget struct {
	Index   string
	ID      string
	Routing string
}
type queryOption struct {
	//为了确保排序字段有序性，这里使用切片（map是无序的，会导致实际字段排序顺序不符合预期）
	Orders               []map[string]bool
	Highlight            *elastic.Highlight
	Profile              bool
	EnableDSL            bool
	ExcludeFields        []string
	IncludeFields        []string
	SlowQueryMillisecond int64
	Preference           string
	FetchSource          *bool
}
type QueryOption func(queryOption *queryOption)

const DefaultPreference = "_local"

func WithOrders(orders []map[string]bool) QueryOption {
	return func(opt *queryOption) {
		opt.Orders = orders
	}
}
func WithHighlight(highlight *elastic.Highlight) QueryOption {
	return func(opt *queryOption) {
		opt.Highlight = highlight
	}
}

func WithProfile(profile bool) QueryOption {
	return func(opt *queryOption) {
		opt.Profile = profile
	}
}

func WithEnableDSL(enableDSL bool) QueryOption {
	return func(opt *queryOption) {
		opt.EnableDSL = enableDSL
	}
}
func WithIncludeFields(includeFields []string) QueryOption {
	return func(opt *queryOption) {
		opt.IncludeFields = includeFields
	}
}

func WithExcludeFields(excludeFields []string) QueryOption {
	return func(opt *queryOption) {
		opt.ExcludeFields = excludeFields
	}
}

func WithSlowQueryMillisecond(slowQueryLogMillisecond int64) QueryOption {
	return func(opt *queryOption) {
		opt.SlowQueryMillisecond = slowQueryLogMillisecond
	}
}

func WithPreference(preference string) QueryOption {
	return func(opt *queryOption) {
		opt.Preference = preference
	}
}

func WithFetchSource(fetchSource bool) QueryOption {
	return func(opt *queryOption) {
		opt.FetchSource = &fetchSource
	}
}

func (c *Client) Get(ctx context.Context, indexName, id, routing string) (*elastic.GetResult, error) {

	//由于副本分片也能提供数据查询，所以当查询请求能从本地分片获取数据时，就不需要转发到其他节点获取数据了，
	//这样可以提高查询缓存命中率，减少跨节点的查询请求，
	//sdk的默认策略是随机获取
	if len(id) < 0 {
		return nil, errors.New("无效的id")
	}
	getService := c.Client.Get().Index(indexName).Id(id).Preference(DefaultPreference)
	if len(routing) > 0 {
		getService.Routing(routing)
	}
	return getService.Do(ctx)
}

func (c *Client) Mget(ctx context.Context, mgetItems []Mget) (*elastic.MgetResponse, error) {
	multiGetService := c.Client.Mget().Preference(DefaultPreference)
	multiGetItems := make([]*elastic.MultiGetItem, 0)
	for _, item := range mgetItems {
		multiGetItem := &elastic.MultiGetItem{}
		multiGetItem.Index(item.Index)
		multiGetItem.Id(item.ID)
		if len(item.Routing) > 0 {
			multiGetItem.Routing(item.Routing)
		}
		multiGetItems = append(multiGetItems, multiGetItem)

	}
	return multiGetService.Add(multiGetItems...).Do(ctx)
}
func (c *Client) Query(ctx context.Context, indexName string, routings []string, query elastic.Query, from, size int, options ...QueryOption) (*elastic.SearchResult, error) {
	queryOpt := &queryOption{}
	for _, f := range options {
		if f != nil {
			f(queryOpt)
		}
	}

	//设置Source
	fetchSource := true
	if queryOpt.FetchSource != nil && !*queryOpt.FetchSource {
		fetchSource = false
	}
	fetchSourceContext := elastic.NewFetchSourceContext(fetchSource)
	if len(queryOpt.IncludeFields) > 0 {
		fetchSourceContext.Include(queryOpt.IncludeFields...)
	}
	if len(queryOpt.ExcludeFields) > 0 {
		fetchSourceContext.Exclude(queryOpt.ExcludeFields...)
	}

	//构造查询条件
	searchSource := elastic.NewSearchSource()
	searchSource = searchSource.FetchSourceContext(fetchSourceContext).Query(query).From(from).Size(size)
	if len(queryOpt.Orders) > 0 {
		for _, orderM := range queryOpt.Orders {
			for field, order := range orderM {
				searchSource.Sort(field, order)
			}
		}
	}
	if queryOpt.Highlight != nil {
		searchSource.Highlight(queryOpt.Highlight)
	}
	searchSource.Profile(queryOpt.Profile)

	searchService := c.Client.Search(indexName).SearchSource(searchSource).IgnoreUnavailable(true).Preference(DefaultPreference)
	if len(routings) > 0 {
		searchService.Routing(routings...)
	}
	if len(queryOpt.Preference) > 0 {
		searchService.Preference(queryOpt.Preference)
	} else {
		searchService.Preference(DefaultPreference)
	}

	res, err := searchService.Do(ctx)
	//获取查询语句
	src, _ := searchSource.Source()
	data, _ := json.Marshal(src)
	rs := strings.Join(routings, ",")
	if c.DebugMode || c.QueryLogEnable || queryOpt.EnableDSL {
		c.logger.Info(fmt.Sprintf("DSL : %s, routing: %s", string(data), rs))
	}
	if queryOpt.SlowQueryMillisecond > 0 && res != nil && res.TookInMillis >= queryOpt.SlowQueryMillisecond {
		c.logger.Info(fmt.Sprintf("slow query DSL : %s, routing: %s", string(data), rs))
	}
	return res, err

}
func (c *Client) ScrollQuery(ctx context.Context, index []string, typeStr string, query elastic.Query, size int, routings []string, callback func(res *elastic.SearchResult, err error), options ...QueryOption) {
	queryOpt := &queryOption{}
	for _, f := range options {
		if f != nil {
			f(queryOpt)
		}
	}
	fetchSource := true
	if queryOpt.FetchSource != nil && !*queryOpt.FetchSource {
		fetchSource = false
	}
	fetchSourceContext := elastic.NewFetchSourceContext(fetchSource)
	searchSource := elastic.NewSearchSource()
	searchSource = searchSource.FetchSourceContext(fetchSourceContext).Query(query)

	if len(queryOpt.Orders) > 0 {
		for _, orderM := range queryOpt.Orders {
			for field, order := range orderM {
				searchSource.Sort(field, order)
			}
		}
	}
	if queryOpt.Highlight != nil {
		searchSource.Highlight(queryOpt.Highlight)
	}
	searchSource.Profile(queryOpt.Profile)
	src, _ := searchSource.Source()
	data, _ := json.Marshal(src)
	rs := strings.Join(routings, ",")
	if c.DebugMode || c.QueryLogEnable || queryOpt.EnableDSL {
		c.logger.Info(fmt.Sprintf("DSL : %s, routing: %s", string(data), rs))
	}
	scrollService := c.Client.Scroll(index...).SearchSource(searchSource).Size(size).Preference(DefaultPreference)
	if len(routings) > 0 {
		scrollService.Routing(routings...)
	}
	if len(queryOpt.Preference) > 0 {
		scrollService.Preference(queryOpt.Preference)
	} else {
		scrollService.Preference(DefaultPreference)
	}
	//scroll保存在ES集群中的上下文信息会占用大量内存资源，虽然会在一段时间后自动清理，当我们知道scroll结束后,
	//需要手动调用clear释放资源
	defer scrollService.Clear(ctx)
	for {
		res, err := scrollService.Do(ctx)
		if err == io.EOF {
			break
		}
		if queryOpt.SlowQueryMillisecond > 0 && res != nil && res.TookInMillis >= queryOpt.SlowQueryMillisecond {
			c.logger.Info(fmt.Sprintf("slow query DSL : %s, routing: %s", string(data), rs))
		}
		if res == nil {
			c.logger.Info("nil results !")
			break
		}
		if res.Hits == nil {
			c.logger.Info("expected results.Hits != nil; got nil")
		}

		if len(res.Hits.Hits) == 0 {
			break
		}
		callback(res, err)
	}
}

// Query 搜索
func (c *Client) SimpleQuery(index, queryString string, ctx context.Context) *elastic.SearchResult {
	var (
		res *elastic.SearchResult
		err error
	)
	//取所有
	res, err = EsConn.Client.Search(index).Do(ctx)
	if len(queryString) > 0 {
		//字段相等
		q := elastic.NewQueryStringQuery(queryString)
		res, err = EsConn.Client.Search(index).
			Query(q).
			Do(context.Background())
	}
	if err != nil {
		c.logger.Error("Query error", zap.Error(err))
		return nil
	}

	return res
}

// List 分页列表
func (c *Client) List(size, page int, index, sortType, sortField, queryString string) *elastic.SearchResult {
	var res *elastic.SearchResult
	var err error
	q := elastic.NewQueryStringQuery(queryString)

	// 排序类型 desc asc es 中只使用 bool 值  true or false
	ascending := true
	if sortType == "desc" {
		ascending = false
	}

	if size < 0 || page < 0 {
		c.logger.Error("param error", zap.Error(err))
		return nil
	}
	if len(queryString) > 0 {
		res, err = EsConn.Client.Search(index).
			Query(q).
			Size(size).
			From((page)*size).
			Sort(sortField, ascending).
			Timeout(Config.Timeout).
			Do(context.Background())

	} else {
		res, err = c.Client.Search("index").
			Size(size).
			From((page)*size).
			Sort(sortField, ascending).
			Timeout(Config.Timeout).
			Do(context.Background())
	}

	if err != nil {
		c.logger.Error("func list error", zap.Error(err))
		return nil
	}
	return res

}

// Aggregation 聚合 平均
func (cli *Client) Aggregation(Params map[string]string) *elastic.SearchResult {
	var res *elastic.SearchResult
	var err error

	//需要聚合的指标 求平均
	avg := elastic.NewAvgAggregation().Field(Params["avg"])
	//单位时间和指定字段
	aggs := elastic.NewDateHistogramAggregation().
		Field(Params["field"]).
		SubAggregation(Params["agg_name"], avg)

	res, err = EsConn.Client.Search(Params["index"]).
		Size(0).
		Aggregation(Params["aggregation_name"], aggs).
		Timeout(Config.Timeout).
		Do(context.Background())

	if err != nil {
		cli.logger.Error("func Aggregation error", zap.Error(err))
		return nil
	}

	return res
}

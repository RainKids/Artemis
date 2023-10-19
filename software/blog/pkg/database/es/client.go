package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/wire"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
	"sync"
	"time"
)

var EsConn *Client
var Config Options

var EStdInfoLogger stdLogger
var EStdErrorLogger stdLogger

func init() {
	EStdInfoLogger = log.New(os.Stdout, "[es] ", log.LstdFlags|log.Lshortfile)
	EStdErrorLogger = log.New(os.Stderr, "[es] ", log.LstdFlags|log.Lshortfile)
}

type stdLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type Client struct {
	Name           string
	Urls           []string
	QueryLogEnable bool
	Client         *elastic.Client
	Bulk           *Bulk
	BulkProcessor  *elastic.BulkProcessor
	DebugMode      bool
	//本地缓存已经创建的索引，用于加速索引是否存在的判断
	CachedIndices sync.Map
	lock          sync.Mutex
	logger        *zap.Logger
	Username      string
	password      string
}

type Bulk struct {
	Name          string
	Workers       int
	FlushInterval time.Duration
	ActionSize    int //每批提交的文档数
	RequestSize   int //每批提交的文档大小
	AfterFunc     elastic.BulkAfterFunc
	Ctx           context.Context
}

// Option ES配置
type Options struct {
	URL                        []string      `yaml:"url"`
	HealthCheck                time.Duration `yaml:"healthCheck"`
	Sniff                      bool          `yaml:"sniff"`
	Gzip                       bool          `yaml:"gzip"`
	Timeout                    string        `yaml:"timeout"`
	QueryLogEnable             bool          `yaml:"queryLogEnable"`
	GlobalSlowQueryMillisecond int64         `yaml:"globalSlowQueryMillisecond"`
	Bulk                       *Bulk         `yaml:"bulk"`
	DebugMode                  bool          `yaml:"debugMode"`
	Scheme                     string        `yaml:"scheme"`
}

// NewOptions for ES
func NewOptions(v *viper.Viper, logger *zap.Logger) *Options {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("es", o); err != nil {
		logger.Error("unmarshal es option error", zap.Error(err))
		return nil
	}

	logger.Info("load es options success", zap.Any("es options", o))
	return o
}

// New 初始化ES连接信息
func New(o *Options, logger *zap.Logger) (esConn *elastic.Client) {

	EsConn = &Client{
		Urls:           o.URL,
		QueryLogEnable: o.QueryLogEnable,
		//Username:       username,
		//password:       password,
		Bulk:          DefaultBulk(o.Bulk.Workers, o.Bulk.ActionSize, o.Bulk.RequestSize, o.Bulk.FlushInterval),
		CachedIndices: sync.Map{},
		lock:          sync.Mutex{},
	}

	esOptions := getBaseOptions(o.HealthCheck, o.Sniff, o.Gzip, o.URL...)
	err := EsConn.newClient(esOptions)
	EsConn.Client = esConn
	if err != nil {
		logger.Error("es NewClient create error", zap.Error(err))
		return
	}
	var esConnInfo *elastic.PingResult
	var code int
	for _, url := range o.URL {
		esConnInfo, code, err = esConn.Ping(url).Do(context.Background())
		if err != nil {
			logger.Error("Ping esConn error", zap.Error(err))
			return
		}
	}
	logger.Info("ES returned with code and version",
		zap.Any("code", code),
		zap.Any("esConn", esConn),
		zap.Any("version", esConnInfo.Version.Number),
	)
	return esConn
}

func (c *Client) newClient(options []elastic.ClientOptionFunc) error {
	client, err := elastic.NewClient(options...)
	if err != nil {
		return err
	}
	c.Client = client

	if c.Bulk.Name == "" {
		c.Bulk.Name = c.Name
	}

	if c.Bulk.Workers <= 0 {
		c.Bulk.Workers = 1
	}

	//参数合理性校验

	if c.Bulk.RequestSize > 100*1024*1024 {
		EStdInfoLogger.Print("Bulk RequestSize must be smaller than 100MB; it will be ignored.")
		c.Bulk.RequestSize = 100 * 1024 * 1024
	}

	if c.Bulk.ActionSize >= 10000 {
		EStdInfoLogger.Print("Bulk ActionSize must be smaller than 10000; it will be ignored.")
		c.Bulk.ActionSize = 10000
	}

	if c.Bulk.FlushInterval >= 60 {
		EStdInfoLogger.Print("Bulk FlushInterval must be smaller than 60s; it will be ignored.")
		c.Bulk.FlushInterval = time.Second * 60
	}
	if c.Bulk.AfterFunc == nil {
		c.Bulk.AfterFunc = defaultBulkFunc
	}
	if c.Bulk.Ctx == nil {
		c.Bulk.Ctx = context.Background()
	}

	c.BulkProcessor, err = c.Client.BulkProcessor().
		Name(c.Bulk.Name).
		Workers(c.Bulk.Workers).
		BulkActions(c.Bulk.ActionSize).
		BulkSize(c.Bulk.RequestSize).
		FlushInterval(c.Bulk.FlushInterval).
		Stats(true).
		After(c.Bulk.AfterFunc).
		Do(c.Bulk.Ctx)
	if err != nil {
		EStdErrorLogger.Print("init bulkProcessor error ", err)
	}
	return nil
}

func getBaseOptions(healthCheck time.Duration, sniff, Gzip bool, urls ...string) []elastic.ClientOptionFunc {
	options := make([]elastic.ClientOptionFunc, 0)
	options = append(options, elastic.SetURL(urls...))
	//options = append(options, elastic.SetBasicAuth(username, password))
	options = append(options, elastic.SetHealthcheckTimeoutStartup(healthCheck*time.Second))
	//开启Sniff，SDK会定期(默认15分钟一次)嗅探集群中全部节点，将全部节点都加入到连接列表中，
	//后续新增的节点也会自动加入到可连接列表，但实际生产中我们可能会设置专门的协调节点，所以默认不开启嗅探
	options = append(options, elastic.SetSniff(sniff))
	options = append(options, elastic.SetErrorLog(EStdInfoLogger))
	options = append(options, elastic.SetGzip(Gzip))
	options = append(options, elastic.SetInfoLog(EStdInfoLogger))
	return options
}

func DefaultBulk(workers, actionSize, requestSize int, flushInterval time.Duration) *Bulk {
	return &Bulk{
		Workers:       workers,
		FlushInterval: flushInterval,
		ActionSize:    actionSize,
		RequestSize:   requestSize, // 5 MB,
		AfterFunc:     defaultBulkFunc,
		Ctx:           context.Background(),
	}
}

func defaultBulkFunc(executionId int64, requests []elastic.BulkableRequest, response *elastic.BulkResponse, err error) {
	if err != nil || (response != nil && response.Errors) {
		res, _ := json.Marshal(response)
		EStdInfoLogger.Printf("executionId: %d ;requests : %v; response : %s ; err : %+v", executionId, requests, res, err)
	}

}

// GetByID 通过ID查找
func (cli *Client) GetByID(Params map[string]string) *elastic.GetResult {
	var (
		res *elastic.GetResult
		err error
	)
	if len(Params["id"]) < 0 {
		fmt.Printf("param error")
		return res
	}

	res, err = EsConn.Client.Get().
		Index(Params["index"]).
		Id(Params["id"]).
		Do(context.Background())

	if err != nil {
		cli.logger.Error("GetByID error", zap.Error(err))
		return nil
	}

	return res
}

func (c *Client) AddIndexCache(indexName ...string) {
	for _, index := range indexName {
		c.CachedIndices.Store(index, true)
	}
}
func (c *Client) DeleteIndexCache(indexName ...string) {
	for _, index := range indexName {
		c.CachedIndices.Delete(index)
	}
}
func (c *Client) Close() error {
	return c.BulkProcessor.Close()
}

// ProviderSet inject es settings
var ProviderSet = wire.NewSet(New, NewOptions)

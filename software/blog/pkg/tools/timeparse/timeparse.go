package timeparse

import (
	"blog/global"
	"bytes"
	"context"
	"database/sql/driver"
	"errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"strconv"
	"strings"
	xtime "time"
)

var cst *xtime.Location

const CSTLayout = global.TimLayOut

type Time int64

func init() {
	var err error
	if cst, err = xtime.LoadLocation("Asia/Shanghai"); err != nil {
		panic(err)
	}
}

// RFC3339ToCSTLayout convert rfc3339 value to china standard time layout
// 2020-11-08T08:18:46+08:00 => 2020-11-08 08:18:46
func RFC3339ToCSTLayout(value string) (string, error) {
	ts, err := xtime.Parse(xtime.RFC3339, value)
	if err != nil {
		return "", err
	}
	return ts.In(cst).Format(CSTLayout), nil
}

// CSTLayoutString 格式化时间
// 返回 "2006-01-02 15:04:05" 格式的时间
func CSTLayoutString() string {
	ts := xtime.Now()
	return ts.In(cst).Format(CSTLayout)
}

// CSTLayoutStringToUnix 返回 unix 时间戳
// 2020-01-24 21:11:11 => 1579871471
func CSTLayoutStringToUnix(cstLayoutString string) (int64, error) {
	stamp, err := xtime.ParseInLocation(CSTLayout, cstLayoutString, cst)
	if err != nil {
		return 0, err
	}
	return stamp.Unix(), nil
}

func CSTLayoutStringToTime(cstLayoutString string) (*xtime.Time, error) {
	stamp, err := xtime.ParseInLocation(CSTLayout, cstLayoutString, cst)
	if err != nil {
		return nil, err
	}
	return &stamp, nil
}

// GeneratorYesterdayPeriod 生成昨日时间区间
// 2006-01-02 00:00:00 ----- 2006-01-03 00:00:00
func GeneratorYesterdayPeriod() (xtime.Time, xtime.Time) {
	now := xtime.Now().In(xtime.UTC)
	begin := xtime.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
	end := xtime.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return begin, end
}

// TimeInterval 时间间隔
type TimeInterval struct {
	Begin xtime.Time
	End   xtime.Time
}

// TimeProcessInterval 时间处理成时间区间
func TimeProcessInterval(t ...xtime.Time) []TimeInterval {
	tis := make([]TimeInterval, 0)
	for _, v := range t {
		tm1 := xtime.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Location())
		tm2 := tm1.AddDate(0, 0, 1)
		ti := TimeInterval{
			Begin: tm1,
			End:   tm2,
		}
		tis = append(tis, ti)
	}
	return tis
}

// TimeProcessReportParams 时间方式验证
func TimeProcessReportParams(style string, ti []string) ([]xtime.Time, error) {
	if style != "day" && style != "week" && style != "month" {
		return nil, errors.New("unsupported_metric_calculation")
	}
	if len(ti) != 2 {
		return nil, errors.New("date_format_error")
	}
	beginStr := ti[0]
	beginTime, err := xtime.Parse("2006-01-02", beginStr)
	if err != nil {
		return nil, errors.New("date_format_error")
	}
	endStr := ti[1]
	endTime, err := xtime.Parse("2006-01-02", endStr)
	if err != nil {
		return nil, errors.New("date_format_error")
	}
	if endTime.Before(beginTime) {
		return nil, errors.New("date_format_error")
	}
	// 按天选择不能超过90天
	if style == "day" && beginTime.AddDate(0, 0, 90).Before(endTime) {
		return nil, errors.New("date_format_error")
	}
	// 按周选择不能超过50周
	if style == "week" && beginTime.AddDate(0, 0, 50*7).Before(endTime) {
		return nil, errors.New("date_format_error")
	}
	// 按月选择不能超过10年
	if style == "month" && beginTime.AddDate(10, 0, 0).Before(endTime) {
		return nil, errors.New("date_format_error")
	}
	// 按周选择 第一天必须是周一 最后一天必须是周日
	if style == "week" && (beginTime.Weekday() != xtime.Weekday(1) || endTime.Weekday() != xtime.Weekday(0)) {
		return nil, errors.New("date_format_error")
	}
	// 按月选择 第一天必须是一号 最后一天必须是月末
	if style == "month" && (beginTime.Day() != 1 || endTime.AddDate(0, 0, 1).Day() != 1) {
		return nil, errors.New("date_format_error")
	}
	return []xtime.Time{beginTime, endTime.Add(xtime.Hour * 24)}, nil
}

// StrTime 格式化时间
func StrTime(times xtime.Time) string {
	var res string
	datetime := "2006-01-02 15:04:05" //待转化为时间戳的字符串
	time1, _ := xtime.ParseInLocation(datetime, times.Format(datetime), xtime.Local)
	atime := time1.Unix()
	if atime < 0 {
		return res
	}
	var byTime = []int64{365 * 24 * 60 * 60, 24 * 60 * 60, 60 * 60, 60, 1}
	var unit = []string{"年前", "天前", "小时前", "分钟前", "秒钟前"}
	now := xtime.Now().Unix()
	ct := now - atime
	if ct < 0 {
		return "刚刚"
	}

	if ct > 30*24*60*60 {
		return times.Format(datetime)
	}
	for i := 0; i < len(byTime); i++ {
		if ct < byTime[i] {
			continue
		}
		var temp = math.Floor(float64(ct / byTime[i]))
		ct = ct % byTime[i]
		if temp > 0 {
			var tempStr string
			tempStr = strconv.FormatFloat(temp, 'f', -1, 64)
			res = MergeString(tempStr, unit[i]) //此处调用了一个我自己封装的字符串拼接的函数（你也可以自己实现）
		}
		break
	}
	return res
}

// StrDate 格式化日期
func StrDate(times xtime.Time) string {
	var res string
	datetime := "2006-01-02" //待转化为时间戳的字符串
	time1, _ := xtime.ParseInLocation(datetime, times.Format(datetime), xtime.Local)
	atime := time1.Unix()
	if atime < 0 {
		return res
	}
	return times.Format(datetime)
}

// MergeString @des 拼接字符串  args ...string 要被拼接的字符串序列 return string
func MergeString(args ...string) string {
	buffer := bytes.Buffer{}
	for i := 0; i < len(args); i++ {
		buffer.WriteString(args[i])
	}
	return buffer.String()
}

// GetNowMonthDays 获取当前月已过的天数(包含当日)
func GetNowMonthDays() int {
	now := xtime.Now()
	day := now.Day()
	return day
}

func (jt *Time) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case xtime.Time:
		*jt = Time(sc.Unix())
	case string:
		var i int64
		i, err = strconv.ParseInt(sc, 10, 64)
		*jt = Time(i)
	}
	return
}

// Value get time value.
func (jt Time) Value() (driver.Value, error) {
	return xtime.Unix(int64(jt), 0), nil
}

// Time get time.
func (jt Time) Time() xtime.Time {
	return xtime.Unix(int64(jt), 0)
}

// Duration be used toml unmarshal string time, like 1s, 500ms.
type Duration xtime.Duration

// UnmarshalText unmarshal text to duration.
func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := xtime.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}

// Shrink will decrease the duration by comparing with context's timeout duration
// and return new timeout\context\CancelFunc.
func (d Duration) Shrink(c context.Context) (Duration, context.Context, context.CancelFunc) {
	if deadline, ok := c.Deadline(); ok {
		if ctimeout := xtime.Until(deadline); ctimeout < xtime.Duration(d) {
			// deliver small timeout
			return Duration(ctimeout), c, func() {}
		}
	}
	ctx, cancel := context.WithTimeout(c, xtime.Duration(d))
	return d, ctx, cancel
}

func ParseWithLocation(name string, timeStr string) (xtime.Time, string, error) {
	locationName := name
	if l, err := xtime.LoadLocation(locationName); err != nil {
		return xtime.Time{}, "", err
	} else {
		//转成带时区的时间
		lt, _ := xtime.ParseInLocation(CSTLayout, timeStr, l)
		//直接转成相对时间
		return lt, lt.Format(CSTLayout), nil
	}
}
func ParseWithLocationToUtc(name string, timeStr string) (string, error) {
	locationName := name
	if l, err := xtime.LoadLocation(locationName); err != nil {
		return "", err
	} else {
		//转成带时区的时间
		lt, _ := xtime.ParseInLocation(CSTLayout, timeStr, l)
		//直接转成相对时间
		//fmt.Println(time.Now().In(l).Format(TIME_LAYOUT))
		loc, err := xtime.LoadLocation("UTC")
		if err != nil {
			return "", err
		}

		lt = lt.In(loc)
		//fmt.Println("--->", lt.Format(TIME_LAYOUT))
		return lt.Format(CSTLayout), nil
	}
}

func TimeStrandFormat(timeStr string) string {
	time_array := strings.Split(timeStr, "-")
	for i := 1; i < 3; i++ {
		if len(time_array[i]) < 2 {
			time_array[i] = "0" + time_array[i]
		}
	}
	return strings.Join(time_array, "-")
}

func GoTimeToPbTime(t xtime.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func PbTimeToGoTime(pbTime *timestamppb.Timestamp) xtime.Time {
	return pbTime.AsTime()
}

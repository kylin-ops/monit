package http_server

import (
	"io"
	"monit/settings"
	"net/http"
	"time"
)

var httpConfig = settings.Config.HttpServer

func httpHandler(w http.ResponseWriter, r *http.Request) {
	var logParser = LogParser{UriAccessCount: map[string]*UrlCount{}, IpAccessCount: map[string]int{}}
	query := r.URL.Query()
	s, ok1 := query["stime"]
	e, ok2 := query["etime"]
	countType, ok3 := query["type"]
	// 获取统计的开始结束时间
	if !ok1 && !ok2 {
		t, _ := time.ParseDuration("-24h")
		logParser.StartTime = time.Now().Add(t)
		logParser.EndTime = time.Now()
	} else if ok1 && ok2 {
		stime, err := time.Parse("2006-01-02T15:04:05", s[0])
		if err != nil {
			io.WriteString(w, "url stime 参数格式错误, 正确格式:2006-01-02T15:04:05")
			return
		}
		etime, err := time.Parse("2006-01-02T15:04:05", e[0])
		if err != nil {
			io.WriteString(w, "url etime 参数格式错误, 正确格式:2006-01-02T15:04:05")
			return
		}
		if etime.Before(stime) {
			io.WriteString(w, "stime的值不能大于etime的值")
			return
		}
		logParser.StartTime = stime
		logParser.EndTime = etime
	}
	if !ok3 {
		countType = []string{"uri"}
	}
	switch countType[0] {
	case "uri":
		urlData, err := logParser.UrlLog(httpConfig.NgxLogPath)
		if err != nil {
			io.WriteString(w, "读取nginx文件错误,错误内容:"+err.Error())
			return
		}
		io.WriteString(w, urlData)
	case "ip":
		ipData, err := logParser.IpLog(httpConfig.NgxLogPath)
		if err != nil {
			io.WriteString(w, "读取nginx文件错误,错误内容:"+err.Error())
			return
		}
		io.WriteString(w, ipData)
	default:
		io.WriteString(w, "url查询的类型参数只能为ip或uri")
	}
}

func Runner() {
	http.HandleFunc("/", httpHandler)
	err := http.ListenAndServe(httpConfig.Address, nil)
	panic(err)
}

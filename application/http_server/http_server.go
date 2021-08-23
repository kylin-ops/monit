package http_server

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type UrlCount struct {
	Count           int
	ResponseTime    float64
	MaxResponseTime float64
}

type LogParser struct {
	StartTime      time.Time
	EndTime        time.Time
	UriAccessCount map[string]*UrlCount
	IpAccessCount  map[string]int
}

func lineParser(line string) (ip, uri string, accessTime time.Time, responseTime float64, err error) {
	reg := regexp.MustCompile("\\s+$")
	line = reg.ReplaceAllString(line, "")
	items := strings.Split(line, " ")
	if len(items) > 7 {
		ip = items[0]
		uri = items[6]
		method := items[5][1:]
		uri = method + ":" + uri
		accessTime, err = time.Parse("[02/Jan/2006:15:04:05", items[3])
		responseTime, _ = strconv.ParseFloat(items[len(items)-1], 64)
	}
	return ip, uri, accessTime, responseTime, err
}

func floatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

func (l *LogParser) NgxLogHandler(logPath, countType, excludeFile string) error {
	f, err := os.Open(logPath)
	rgx := regexp.MustCompile(excludeFile)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, "[error]") {
			continue
		}
		ip, uri, accessTime, responseTime, err := lineParser(line)
		uri = strings.Split(uri, "?")[0]
		if err != nil {
			continue
		}
		if rgx.MatchString(uri) {
			continue
		}
		if accessTime.Before(l.EndTime) && accessTime.After(l.StartTime) {
			if countType == "url" {
				if _, ok := l.UriAccessCount[uri]; ok {
					l.UriAccessCount[uri].Count += 1
					l.UriAccessCount[uri].ResponseTime += responseTime
					if l.UriAccessCount[uri].MaxResponseTime < responseTime {
						l.UriAccessCount[uri].MaxResponseTime = responseTime
					}
				} else {
					l.UriAccessCount[uri] = &UrlCount{}
				}
			} else if countType == "ip" {
				l.IpAccessCount[ip] += 1
			}
		}
		if accessTime.After(l.EndTime) {
			break
		}
	}
	return nil
}

func (l *LogParser) UrlLog(logPath string) (content string, err error) {
	content = fmt.Sprintf("%-9s : %-9s : %-9s : %s\n", "访问次数", "最大响应时间", "平均响应时间", "URL")
	if err := l.NgxLogHandler(logPath, "url", httpConfig.ExcludeFile); err != nil {
		return content, err
	}
	for k, v := range l.UriAccessCount {
		content += fmt.Sprintf("%-9d : %-9f : %-9f : %s\n", v.Count, v.MaxResponseTime, floatRound(v.ResponseTime/float64(v.Count), 4), k)
	}
	return content, nil
}

func (l *LogParser) IpLog(logPath string) (content string, err error) {
	content = fmt.Sprintf("%-9s : %s\n", "IP地址", "访问次数")
	if err := l.NgxLogHandler(logPath, "url", httpConfig.ExcludeFile); err != nil {
		return content, err
	}
	for k, v := range l.IpAccessCount {
		content += fmt.Sprintf("%-9d : %s\n", v, k)
	}
	return content, nil
}

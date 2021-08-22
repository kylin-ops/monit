package http_server

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type LogParser struct {
	StartTime      time.Time
	EndTime        time.Time
	UriAccessCount map[string]int
	IpAccessCount  map[string]int
}

func lineParser(line string) (ip, uri string, accessTime time.Time, err error) {
	items := strings.Split(line, " ")
	if len(items) > 7 {
		ip = items[0]
		uri = items[6]
		method := items[5][1:]
		uri = method + ":" + uri
		accessTime, err = time.Parse("[02/Jan/2006:15:04:05", items[3])
	}
	return ip, uri, accessTime, err
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
		ip, uri, accesstime, err := lineParser(line)
		if err != nil {
			continue
		}
		if rgx.MatchString(uri) {
			continue
		}
		if accesstime.Before(l.EndTime) && accesstime.After(l.StartTime) {
			if countType == "url" {
				l.UriAccessCount[uri] += 1
			} else if countType == "ip" {
				l.IpAccessCount[ip] += 1
			}
		}
		if accesstime.After(l.EndTime) {
			break
		}
	}
	return nil
}

func (l *LogParser) UrlLog(logPath string) (content string, err error) {
	content = fmt.Sprintf("%-9s : %s\n", "URL", "访问次数")
	if err := l.NgxLogHandler(logPath, "url", httpConfig.ExcludeFile); err != nil {
		return content, err
	}
	for k, v := range l.UriAccessCount {
		content += fmt.Sprintf("%-9d : %s\n", v, k)
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

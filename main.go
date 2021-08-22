package main

import (
	"monit/application/http_detector"
	"monit/application/http_server"
	"monit/settings"
	"sync"
)

var logger = settings.Logger

func main() {
	logger.Infoln("启动服务")
	wg := sync.WaitGroup{}
	wg.Add(1)
	http_detector.Runner()
	http_server.Runner()
	wg.Wait()
}

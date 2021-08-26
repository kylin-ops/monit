package settings

import (
	"encoding/json"
	"github.com/kylin-ops/logger"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

type log struct {
	Level     string `json:"level"`
	Path      string `json:"path"`
	RollTime  int    `json:"roll_time"`
	Count     int    `json:"count"`
	IsConsole bool   `json:"is_console"`
}

type aliyunSms struct {
	RegionId        string `json:"region_id"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	SignName        string `json:"sign_name"`
	TemplateCode    string `json:"template_code"`
	Phones          string `json:"phones"`
}

type Detectors struct {
	Url    string `json:"url"`
	Method string `json:"method"`
}

type Login struct {
	Detectors
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type httpDetectors struct {
	Interval       int         `json:"interval"`
	FailedRetry    int         `json:"failed_retry"`
	AlertPreHour   int         `json:"alert_pre_hour"`
	StartAlertTime string      `json:"start_alert_time"`
	EndAlertTime   string      `json:"end_alert_time"`
	Login          Login       `json:"login"`
	Detectors      []Detectors `json:"detectors"`
}

type httpServer struct {
	Address     string `json:"address"`
	NgxLogPath  string `json:"ngx_log_path"`
	ExcludeFile string `json:"exclude_file"`
}

type config struct {
	HttpServer    httpServer    `json:"http_server"`
	Logger        log           `json:"logger"`
	AliyunSms     aliyunSms     `json:"aliyun_sms"`
	HttpDetectors httpDetectors `json:"http_detectors"`
}

func initConfig() *config {
	var conf config
	f, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(f, &conf); err != nil {
		panic(err)
	}
	return &conf
}

func initLogger() *logrus.Logger {
	logg, err := logger.NewLogger(&logger.Options{
		Level:     Config.Logger.Level,
		Path:      Config.Logger.Path,
		RollTime:  time.Duration(Config.Logger.RollTime) * time.Hour,
		LogCount:  Config.Logger.Count,
		IsConsole: Config.Logger.IsConsole,
	})
	if err != nil {
		panic(err)
	}
	return logg
}

func initSms() *AliyunSms {
	return &AliyunSms{
		RegionId:        Config.AliyunSms.RegionId,
		AccessKeyId:     Config.AliyunSms.AccessKeyId,
		AccessKeySecret: Config.AliyunSms.AccessKeySecret,
		SignName:        Config.AliyunSms.SignName,
		TemplateCode:    Config.AliyunSms.TemplateCode,
		Phones:          Config.AliyunSms.Phones,
	}
}

var (
	Config = initConfig()
	Logger = initLogger()
	Sms    = initSms()
)

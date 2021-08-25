package http_detector

import (
	"errors"
	"fmt"
	"github.com/kylin-ops/grequests"
	"github.com/kylin-ops/timer"
	"monit/settings"
	"strconv"
	"strings"
	"time"
)

var Timer = timer.NewTimer()
var logger = settings.Logger
var sms = settings.Sms

var smsTemplate = `
级别:【告警】
URl:%s
事件:%s
`

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

type HttpDetector struct {
	Interval       int         `json:"interval"`
	FailedRetry    int         `json:"failed_retry"`
	AlertPreHour   int         `json:"alert_pre_hour"`
	StartAlertTime string      `json:"start_alert_time"`
	EndAlertTime   string      `json:"end_alert_time"`
	Token          string      `json:"-"`
	Login          Login       `json:"login"`
	Detectors      []Detectors `json:"detectors"`
}

// 登录获取token
func (h *HttpDetector) HttpLogin() error {
	var token, tokenType string
	var o bool
	req, err := grequests.Get(h.Login.Url, &grequests.RequestOptions{
		Header: map[string]string{"Authorization": "Basic d2ViYXBwOndlYmFwcA=="},
		Params: map[string]string{
			"grant_type": "password",
			"username":   h.Login.Username,
			"password":   h.Login.Password,
		},
		Timeout: time.Second * 2,
	})
	if err != nil {
		return err
	}
	defer req.Close()
	if req.StatusCode() != 200 {
		return errors.New("登录测试失败,响应码:%s" + strconv.Itoa(req.StatusCode()))
	}
	var data map[string]interface{}
	if err = req.Json(&data); err != nil {
		return err
	}
	if t, ok := data["access_token"]; !ok {
		return errors.New("登录响应数据里没有\"access_token\"")
	} else {
		if token, o = t.(string); !o {
			return errors.New("登录响应数据里\"access_token\"数据不是字符串")
		}
	}
	if tt, ok := data["token_type"]; !ok {
		return errors.New("登录响应数据里没有\"token_type\"")
	} else {
		if tokenType, o = tt.(string); !o {
			return errors.New("登录响应数据里\"token_type\"数据不是字符串")
		}
	}
	h.Token = tokenType + " " + token
	return nil
}

func (h *HttpDetector) timeRange() (timeRange map[string]bool) {
	timeRange = map[string]bool{}
	st, err := time.ParseInLocation("15:04:05", h.StartAlertTime, time.Local)
	if err != nil {
		panic(err)
	}
	et, err := time.ParseInLocation("15:04:05", h.EndAlertTime, time.Local)
	if err != nil {
		panic(err)
	}
	hour, _ := time.ParseDuration("1h")
	timeRange[st.Format("15")] = true
	timeRange[et.Format("15")] = true
	for {
		st = st.Add(hour)
		if st.Format("15") == et.Format("15") {
			break
		}
		timeRange[st.Format("15")] = true
	}
	return timeRange
}

func (h *HttpDetector) DetectExec(method, url string) error {
	var req *grequests.Response
	var options = &grequests.RequestOptions{Header: map[string]string{
		"Authorization": h.Token,
	}}
	var err error
	switch strings.ToUpper(method) {
	case "GET":
		req, err = grequests.Get(url, options)
	case "POST":
		req, err = grequests.Post(url, options)
	case "PUT":
		req, err = grequests.Put(url, options)
	case "DELETE":
		req, err = grequests.Delete(url, options)
	case "PATCH":
		req, err = grequests.Patch(url, options)
	}
	if err != nil {
		return err
	}
	defer req.Close()
	if req.StatusCode() > 399 && req.StatusCode() < 200 {
		return nil
	}
	return errors.New(fmt.Sprintf("探测失败，服务响应码是%d", req.StatusCode()))
}

// 按照规则执行
func (h *HttpDetector) HttpDetect() {
	hour := time.Now().Format("15")
	alertCount := 0
	err := Timer.Add("http_check", time.Second*time.Duration(h.Interval), func() {
		timeRange := h.timeRange()
		var err error
		for i := 1; i < h.FailedRetry+1; i++ {
			err = h.HttpLogin()
			if err != nil {
				logger.Infof(fmt.Sprintf("%s: 第%d登录失败, %d秒后重试", h.Login.Url, i, h.FailedRetry))
				time.Sleep(time.Duration(h.FailedRetry) * time.Second)
				continue
			}
			break
		}
		for i := 1; i < h.FailedRetry+1; i++ {
			for _, addr := range h.Detectors {
				err = h.DetectExec(addr.Method, addr.Url)

			}
		}

	})
	if err != nil {
		panic(err)
	}
}

package http_detector

import (
	"errors"
	"fmt"
	"github.com/kylin-ops/grequests"
	"github.com/kylin-ops/timer"
	"monit/settings"
	"strconv"
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

type HttpDetector struct {
	Url            string
	Username       string
	Password       string
	Interval       int
	FailedRetry    int
	AlertPreHour   int
	StartAlertTime string
	EndAlertTime   string
}

// 执行http探测
func (h *HttpDetector) HttpLoginDetect() error {
	req, err := grequests.Get(h.Url, &grequests.RequestOptions{
		Header: map[string]string{"Authorization": "Basic d2ViYXBwOndlYmFwcA=="},
		Params: map[string]string{
			"grant_type": "password",
			"username":   h.Username,
			"password":   h.Password,
		},
		Timeout: time.Second * 2,
	})
	if err != nil {
		return err
	}
	if req.StatusCode() != 200 {
		return errors.New("登录测试失败,响应码:%s" + strconv.Itoa(req.StatusCode()))
	}
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

// 按照规则执行
func (h *HttpDetector) HttpDetect() {
	hour := time.Now().Format("15")
	alertCount := 0
	err := Timer.Add("http_check", time.Second*time.Duration(h.Interval), func() {
		timeRange := h.timeRange()
		var err error
		for i := 1; i < h.FailedRetry+1; i++ {
			err = h.HttpLoginDetect()
			if err == nil {
				logger.Infof(fmt.Sprintf("%s: 第%d探测成功", h.Url, i))
				break
			}
			logger.Warnf(fmt.Sprintf("%s: 第%d探测失败,错误信息:%s, 等待3秒重试", h.Url, i, err.Error()))
			time.Sleep(time.Second * 3)
		}
		if err != nil {
			// 每小时告警重置
			if hour != time.Now().Format("15") {
				alertCount = 0
			}
			// 控制每小时发送3次, 时间在指定的时间范围
			if _, ok := timeRange[hour]; alertCount < h.AlertPreHour && ok {
				// 发送告警
				alertCount += 1
				err := sms.SendSms(fmt.Sprintf(smsTemplate, h.Url, "登录测试失败"))
				if err != nil {
					logger.Errorf("发送告警短信失败,错误信息:%s", err.Error())
				} else {
					logger.Info("发送告警短信成功")
				}
			}
		}
	})
	if err != nil {
		panic(err)
	}
}

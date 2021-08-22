package http_detector

import (
	"monit/settings"
)

func Runner() {
	dectors := settings.Config.HttpDetectors
	for _, d := range dectors {
		var dector = &HttpDetector{
			Url:            d.Url,
			Username:       d.Username,
			Password:       d.Password,
			Interval:       d.Interval,
			FailedRetry:    d.FailedRetry,
			AlertPreHour:   d.AlertPreHour,
			StartAlertTime: d.StartAlertTime,
			EndAlertTime:   d.EndAlertTime,
		}
		dector.HttpDetect()
	}
}

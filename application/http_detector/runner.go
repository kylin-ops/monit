package http_detector

import (
	"encoding/json"
	"monit/settings"
)

func Runner() {
	d, _ := json.Marshal(settings.Config.HttpDetectors)
	detector := &HttpDetector{}
	_ = json.Unmarshal(d, &detector)
	detector.HttpDetect()
}

package settings

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type AliyunSms struct {
	RegionId        string `json:"region_id"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	SignName        string `json:"sign_name"`
	TemplateCode    string `json:"template_code"`
	Phones          string `json:"phones"`
}

func (s *AliyunSms) SendSms(smsContext string) error {
	client, err := dysmsapi.NewClientWithAccessKey(s.RegionId, s.AccessKeyId, s.AccessKeySecret)
	if err != nil {
		return err
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = s.Phones
	request.SignName = s.SignName
	request.TemplateCode = s.TemplateCode
	request.TemplateParam = fmt.Sprintf(`{"content": "%s"}`, smsContext)
	request.Content = []byte("my sms")
	_, err = client.SendSms(request)
	return err
}

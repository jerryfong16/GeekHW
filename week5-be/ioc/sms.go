package ioc

import "geek-hw-week5/internal/service/sms"

func InitLocalSMSService() sms.SMSService {
	return sms.NewLocalSMSService()
}

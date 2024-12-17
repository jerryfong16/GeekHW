package ioc

import "geek-hw-week4/internal/service/sms"

func InitLocalSMSService() sms.SMSService {
	return sms.NewLocalSMSService()
}

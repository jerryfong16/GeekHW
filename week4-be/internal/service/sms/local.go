package sms

import (
	"context"
	"log"
)

type LocalSMSService struct {
}

func (svc *LocalSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	log.Println("验证码是", args)
	return nil
}

func NewLocalSMSService() *LocalSMSService {
	return &LocalSMSService{}
}

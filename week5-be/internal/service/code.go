package service

import (
	"context"
	"errors"
	"fmt"
	"geek-hw-week4/internal/repository"
	"geek-hw-week4/internal/service/sms"
	"math/rand"
)

var ErrCodeSendTooMany = repository.ErrCodeVerifyTooMany

type CodeService interface {
	Send(ctx context.Context, businessType, phone string) error
	Verify(ctx context.Context, businessType, phone, inputCode string) (bool, error)
}

type SMSCodeService struct {
	repo repository.CodeRepository
	sms  sms.SMSService
}

func (svc *SMSCodeService) Send(ctx context.Context, businessType, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, businessType, phone, code)
	if err != nil {
		return err
	}
	const codeTplId = "1234567"
	return svc.sms.Send(ctx, codeTplId, []string{code}, phone)
}

func (svc *SMSCodeService) Verify(ctx context.Context, businessType, phone, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, businessType, phone, inputCode)
	if errors.Is(err, repository.ErrCodeVerifyTooMany) {
		return false, nil
	}
	return ok, err
}

func (svc *SMSCodeService) generate() string {
	// 0-999999
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.SMSService) CodeService {
	return &SMSCodeService{repo: repo, sms: smsSvc}
}

package worker

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/ismail118/simple-bank/mail"
	"github.com/ismail118/simple-bank/repository"
)

type RedisTaskProcessorMock struct {
}

func NewRedisTaskProcessorMock(redisOpt asynq.RedisClientOpt, store repository.Store, mailer mail.SenderEmail, gatewaySeverAddress string) TaskProcessor {
	return &RedisTaskProcessorMock{}
}

func (p *RedisTaskProcessorMock) Start() error {
	return nil
}

func (p *RedisTaskProcessorMock) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	return nil
}

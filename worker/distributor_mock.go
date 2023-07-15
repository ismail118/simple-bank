package worker

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
)

type RedisTaskDistributorMock struct {
}

func NewRedisTaskDistributorMock(redisOpt asynq.RedisClientOpt) TaskDistributor {
	return &RedisTaskDistributorMock{}
}

func (d *RedisTaskDistributorMock) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	if payload.Username == "user2" {
		return fmt.Errorf("some error")
	}
	return nil
}

package worker

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/ismail118/simple-bank/mail"
	"github.com/ismail118/simple-bank/repository"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server              *asynq.Server
	store               repository.Store
	mailer              mail.SenderEmail
	gatewaySeverAddress string
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store repository.Store, mailer mail.SenderEmail, gatewaySeverAddress string) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 6,
			QueueDefault:  3,
			QueueLow:      1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().
				Err(err).
				Str("type", task.Type()).
				Bytes("payload", task.Payload()).
				Msg("process task failed")
		}),
		Logger: NewLogger(),
	})

	return &RedisTaskProcessor{
		server:              server,
		store:               store,
		mailer:              mailer,
		gatewaySeverAddress: gatewaySeverAddress,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, p.ProcessTaskSendVerifyEmail)

	return p.server.Start(mux)
}

package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/ismail118/simple-bank/models"
	"github.com/ismail118/simple-bank/util"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (d *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload)
	info, err := d.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return err
	}

	log.Info().
		Str("type", task.Type()).
		Str("queue", info.Queue).
		Bytes("payload", jsonPayload).
		Int("max_retry", info.MaxRetry).
		Msg("enqueued task")

	return nil
}

func (p *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		// if task error on unmarshall there is no point to retrying it
		// and we tell asynq about that by return error asynq.SkipRetry
		return fmt.Errorf("error unmarshall payload: %w", asynq.SkipRetry)
	}

	user, err := p.store.GetUsersByUsername(ctx, payload.Username)
	if err != nil {
		return fmt.Errorf("error get user payload: %w", err)
	}
	if user.Username == "" {
		return fmt.Errorf("username doesn't exist: %w", err)
	}

	verifyEmail := models.VerifyEmail{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	}

	id, err := p.store.InsertVerifyEmail(ctx, verifyEmail)
	if err != nil {
		return fmt.Errorf("cannot insert verify_email: %s", err)
	}
	verifyEmail.ID = id

	subject := "Welcome to Simple Bank"
	verifyUrl := fmt.Sprintf("http://%s/v1/verify_email?id=%d&secret_code=%s", p.gatewaySeverAddress, verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here<a/> to verify your email address.<br/>
	`, user.FullName, verifyUrl)
	to := []string{user.Email}
	err = p.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %s", err)
	}

	log.Info().
		Str("type", task.Type()).
		Str("email", user.Email).
		Bytes("payload", task.Payload()).
		Msg("process task")

	return nil
}

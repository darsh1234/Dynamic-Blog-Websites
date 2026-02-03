package email

import (
	"context"
	"log/slog"
)

type StubSender struct {
	logger *slog.Logger
}

func NewStubSender(logger *slog.Logger) *StubSender {
	return &StubSender{logger: logger}
}

func (s *StubSender) Send(_ context.Context, message Message) error {
	s.logger.Info("stub email sender",
		"to", message.To,
		"subject", message.Subject,
		"body", message.Body,
	)
	return nil
}

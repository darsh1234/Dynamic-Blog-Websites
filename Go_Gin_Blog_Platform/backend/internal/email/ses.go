package email

import (
	"context"
	"fmt"
)

type SESConfig struct {
	Region  string
	From    string
	FromARN string
}

type SESSender struct {
	config SESConfig
}

func NewSESSender(config SESConfig) *SESSender {
	return &SESSender{config: config}
}

func (s *SESSender) Send(_ context.Context, _ Message) error {
	if s.config.Region == "" || s.config.From == "" {
		return fmt.Errorf("ses sender is not fully configured")
	}

	return fmt.Errorf("ses sender is not wired in this build")
}

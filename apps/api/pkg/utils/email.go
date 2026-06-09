package utils

import (
	"context"
	"net/smtp"

	"github.com/rentoutdoor/api/internal/infrastructure/config"
)

func Send(ctx context.Context, to string, subject string, body string, cfg config.SMTPConfig) error {
	auth := smtp.PlainAuth(
		"",
		cfg.User,
		cfg.Pass,
		cfg.Host,
	)

	msg := []byte(
		"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
			body,
	)

	return smtp.SendMail(cfg.Host+":"+cfg.Port, auth, cfg.From, []string{to}, msg)
}

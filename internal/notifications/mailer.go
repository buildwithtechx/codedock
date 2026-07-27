package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"time"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/services/system"
)

type MailerService struct {
	notifSettingsService *system.NotificationSettingsService
}

func NewMailerService(notifSettings *system.NotificationSettingsService) (*MailerService, error) {
	if err := LoadTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load email templates: %w", err)
	}
	return &MailerService{
		notifSettingsService: notifSettings,
	}, nil
}

func (s *MailerService) SendSystemEmail(ctx context.Context, templateName string, toAddress string, subject string, data any) error {
	settings, _ := s.notifSettingsService.GetNotificationSettings(ctx)
	cfg := config.Get()

	var buf bytes.Buffer
	if err := HTMLTemplates.ExecuteTemplate(&buf, templateName, data); err != nil {
		return fmt.Errorf("executing template %s: %w", templateName, err)
	}
	htmlContent := buf.String()

	resendKey := cfg.Resend.APIKey
	if settings != nil && settings.ResendEnabled && settings.ResendAPIKey != "" {
		resendKey = settings.ResendAPIKey
	}

	fromAddress := cfg.SMTP.From
	if settings != nil && settings.SMTPFromAddress != "" {
		fromAddress = settings.SMTPFromAddress
	}

	if resendKey != "" {
		return sendResendEmail(ctx, resendKey, fromAddress, toAddress, subject, htmlContent)
	}

	var host, port, user, pass string
	if settings != nil && settings.SMTPEnabled {
		host = settings.SMTPHost
		port = fmt.Sprintf("%d", settings.SMTPPort)
		user = settings.SMTPUser
		pass = settings.SMTPPassword
		if settings.SMTPFromAddress != "" {
			fromAddress = settings.SMTPFromAddress
		}
	} else {
		host = cfg.SMTP.Host
		port = fmt.Sprintf("%d", cfg.SMTP.Port)
		user = cfg.SMTP.User
		pass = cfg.SMTP.Pass
	}

	if host == "" || fromAddress == "" {
		return fmt.Errorf("no email provider configured (missing RESEND_API_KEY or SMTP_HOST/SMTP_FROM)")
	}

	msg := fmt.Appendf(nil, "To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
		"%s", toAddress, fromAddress, subject, htmlContent)

	var auth smtp.Auth
	if user != "" && pass != "" {
		auth = smtp.PlainAuth("", user, pass, host)
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	err := smtp.SendMail(addr, auth, fromAddress, []string{toAddress}, msg)
	if err != nil {
		return fmt.Errorf("smtp.SendMail: %w", err)
	}

	return nil
}

func sendResendEmail(ctx context.Context, apiKey, from, to, subject, htmlContent string) error {
	if from == "" {
		from = "Codedock <no-reply@codedock.run>"
	}
	payload := map[string]any{
		"from":    from,
		"to":      []string{to},
		"subject": subject,
		"html":    htmlContent,
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling resend payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(jsonBytes))
	if err != nil {
		return fmt.Errorf("creating resend request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending resend request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend API error (%d): %s", resp.StatusCode, string(body))
	}
	return nil
}

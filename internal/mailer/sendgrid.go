package mailer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	sendgridMail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGridMailer(fromEmail, apiKey string) *SendGridMailer {
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    sendgrid.NewSendClient(apiKey),
	}
}

func (m *SendGridMailer) Send(ctx context.Context, templateFile string, mail Mail, isSandbox bool) error {
	if len(mail.To) == 0 {
		return errors.New("no to addresses provided to send email")
	}

	from := sendgridMail.NewEmail(fromName, m.fromEmail)
	sandbox := &sendgridMail.MailSettings{
		SandboxMode: &sendgridMail.Setting{
			Enable: &isSandbox,
		},
	}

	for _, recipient := range mail.To {
		if err := ctx.Err(); err != nil {
			return err
		}

		to := sendgridMail.NewEmail(recipient.Name, recipient.Email)
		templatePath := fmt.Sprintf("templates/%s", templateFile)
		tmpl, err := template.ParseFS(FS, templatePath)
		if err != nil {
			return err
		}

		subject := new(bytes.Buffer)
		if err := tmpl.ExecuteTemplate(subject, "subject", mail.Args); err != nil {
			return err
		}
		mail.Subject = subject.String()

		htmlBody := new(bytes.Buffer)
		if err := tmpl.ExecuteTemplate(htmlBody, "body", mail.Args); err != nil {
			return err
		}

		message := sendgridMail.NewSingleEmail(from, mail.Subject, to, "", htmlBody.String())
		message.SetMailSettings(sandbox)

		if err := m.sendWithRetry(ctx, message); err != nil {
			return err
		}
	}

	return nil
}

func (m *SendGridMailer) sendWithRetry(ctx context.Context, message *sendgridMail.SGMailV3) error {
	backoff := initialBackoff

	for attempt := 1; attempt <= maxSendAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		// Successfully sent the email
		response, err := m.client.SendWithContext(ctx, message)
		if err == nil && response.StatusCode >= 200 && response.StatusCode < 300 {
			return nil
		}

		// Non-retryable HTTP outcome (4xx except 429)
		if err == nil && response != nil && !retryableStatus(response.StatusCode) {
			return fmt.Errorf("failed to send email: %d: %s", response.StatusCode, response.Body)
		}

		// Last attempt: return best error we have
		if attempt == maxSendAttempts {
			if err != nil {
				return err
			}
			if response != nil {
				return fmt.Errorf("sendgrid: status %d: %s", response.StatusCode, response.Body)
			}
			return errors.New("sendgrid: empty response")
		}
		// Decide wait: prefer SendGrid reset time on 429, else exponential backoff
		delay := backoff
		if response != nil && response.StatusCode == http.StatusTooManyRequests {
			if d := rateLimitDelay(response); d > 0 {
				delay = d
			}
		}

		if delay > maxBackoff {
			delay = maxBackoff
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		if backoff < maxBackoff {
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}

	return errors.New("failed to send email after max send attempts")
}

func retryableStatus(code int) bool {
	return code == http.StatusTooManyRequests || code >= 500
}

func rateLimitDelay(resp *rest.Response) time.Duration {
	if resp == nil || resp.Headers == nil {
		return 0
	}
	v := resp.Headers["X-RateLimit-Reset"]
	if len(v) == 0 || v[0] == "" {
		return 0
	}
	sec, err := strconv.Atoi(v[0])
	if err != nil {
		return 0
	}
	t := time.Unix(int64(sec), 0)
	d := time.Until(t)
	if d < 0 {
		return 0
	}
	return d
}

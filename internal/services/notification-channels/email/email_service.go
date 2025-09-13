//nolint:gosec,gocyclo // it's ok
package email

import (
	"bytes"
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"html/template"
	"log/slog"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-mail/mail"

	"github.com/rom8726/etoggle/internal/domain"
)

const (
	timeout = 30 * time.Second
)

//go:embed templates/reset_password_email.tmpl
var resetPasswordEmailTemplate string

//go:embed templates/2fa_code_email.tmpl
var twoFACodeEmailTemplate string

type Service struct {
	cfg          *Config
	usersRepo    UsersRepository
	projectsRepo ProjectsRepository

	sendEmailFunc func(ctx context.Context, toEmails []string, subject, body string) error
}

type Config struct {
	SMTPHost      string
	Username      string
	Password      string
	CertFile      string
	KeyFile       string
	AllowInsecure bool
	UseTLS        bool

	BaseURL string
	From    string
	LogoURL string
}

func New(
	cfg *Config,
	usersRepo UsersRepository,
	projectsRepo ProjectsRepository,
) *Service {
	service := &Service{
		cfg:          cfg,
		usersRepo:    usersRepo,
		projectsRepo: projectsRepo,
	}
	service.sendEmailFunc = service.SendEmail

	return service
}

func (s *Service) Type() domain.NotificationType {
	return domain.NotificationTypeEmail
}

// Send2FACodeEmail sends a 2FA confirmation code for a specific action (disable/reset).
func (s *Service) Send2FACodeEmail(ctx context.Context, email, code, action string) error {
	subject := "eToggle: 2FA confirmation code"
	var actionText string
	switch action {
	case "disable":
		actionText = "to disable two-factor authentication"
	case "reset":
		actionText = "to reset two-factor authentication"
	default:
		actionText = "for your action"
	}

	tmpl, err := template.New("2fa_code_email").Parse(twoFACodeEmailTemplate)
	if err != nil {
		return err
	}
	var body bytes.Buffer
	err = tmpl.Execute(&body, struct {
		Code       string
		ActionText string
	}{
		Code:       code,
		ActionText: actionText,
	})
	if err != nil {
		return err
	}

	return s.SendEmail(ctx, []string{email}, subject, body.String())
}

func (s *Service) SendResetPasswordEmail(ctx context.Context, email, token string) error {
	slog.Debug("sending reset password email", "base_url", s.cfg.BaseURL)

	tpl, err := template.New("reset_password").Parse(resetPasswordEmailTemplate)
	if err != nil {
		slog.Error("failed to parse reset password template", "error", err)

		return fmt.Errorf("parse template: %w", err)
	}

	resetLink := s.cfg.BaseURL + "/reset-password?token=" + token
	slog.Debug("generated reset link", "reset_link", resetLink)

	renderData := struct {
		ResetLink string
	}{
		ResetLink: resetLink,
	}

	var body bytes.Buffer
	if err := tpl.Execute(&body, renderData); err != nil {
		slog.Error("failed to execute reset password template", "error", err)

		return fmt.Errorf("execute template: %w", err)
	}

	err = s.SendEmail(ctx, []string{email}, "eToggle: Reset Your Password", body.String())
	if err != nil {
		slog.Error("failed to send reset password email", "error", err)

		return err
	}

	slog.Info("reset password email sent successfully")

	return nil
}

// SendEmail builds a MIME message and sends it via SMTP.
//
//nolint:gosec,nestif // it's ok here
func (s *Service) SendEmail(ctx context.Context, toEmails []string, subject, bodyHTML string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	slog.Debug("starting email send", "to", strings.Join(toEmails, ", "), "subject", subject, "smtp_host", s.cfg.SMTPHost)

	from := s.cfg.From
	if from == "" {
		from = s.cfg.Username
	}

	// --- Build message ------------------------------------------------------
	msg := mail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", toEmails...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html; charset=UTF-8", bodyHTML)

	// --- Dialer -------------------------------------------------------------
	host, portStr, err := net.SplitHostPort(s.cfg.SMTPHost)
	if err != nil {
		slog.Error("invalid smtp host configuration", "smtp_host", s.cfg.SMTPHost, "error", err)

		return fmt.Errorf("invalid smtp host: %w", err)
	}
	port, _ := strconv.Atoi(portStr)

	slog.Debug("creating SMTP dialer", "host", host, "port", port, "username", s.cfg.Username)
	dialer := mail.NewDialer(host, port, s.cfg.Username, s.cfg.Password)
	dialer.Timeout = timeout

	if s.cfg.UseTLS {
		var certs []tls.Certificate
		if s.cfg.CertFile != "" {
			cert, err := tls.LoadX509KeyPair(s.cfg.CertFile, s.cfg.KeyFile)
			if err != nil {
				slog.Error("failed to load TLS certificate", "cert_file", s.cfg.CertFile,
					"key_file", s.cfg.KeyFile, "error", err)

				return fmt.Errorf("load TLS key pair: %w", err)
			}
			certs = []tls.Certificate{cert}
		}

		dialer.TLSConfig = &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: s.cfg.AllowInsecure,
			Certificates:       certs,
		}
		slog.Debug("using TLS for SMTP connection", "host", host, "port", port)
	} else {
		// For MailHog, explicitly set TLS config to nil to avoid any TLS attempts
		if port == 1025 {
			dialer.TLSConfig = nil
			slog.Warn("explicitly set TLS config to nil for MailHog")
		}
		slog.Debug("using unencrypted SMTP connection", "host", host, "port", port)
	}

	// --- Send with context --------------------------------------------------
	slog.Debug("attempting to send email", "host", host, "port", port, "use_tls", s.cfg.UseTLS)

	// Special handling for MailHog to avoid TLS issues
	if port == 1025 {
		slog.Debug("using direct SMTP for MailHog")
		err = s.sendEmailDirectSMTP(ctx, host, port, s.cfg.Username, from, toEmails, subject, bodyHTML)
	} else {
		errCh := make(chan error, 1)
		go func() {
			errCh <- dialer.DialAndSend(msg)
		}()

		select {
		case <-ctx.Done():
			slog.Error("email send timeout", "to", strings.Join(toEmails, ", "), "subject", subject,
				"error", ctx.Err())

			return fmt.Errorf("send mail: %w", ctx.Err())
		case err = <-errCh:
		}
	}

	if err != nil {
		slog.Error("failed to send email", "to", strings.Join(toEmails, ", "), "subject", subject,
			"smtp_host", s.cfg.SMTPHost, "host", host, "port", port, "use_tls", s.cfg.UseTLS, "error", err)

		return fmt.Errorf("send mail: %w", err)
	}

	slog.Debug("email sent successfully", "to", strings.Join(toEmails, ", "), "subject", subject)

	return nil
}

// sendEmailsParallel sends emails in parallel with a limit on the number of workers.
func (s *Service) sendEmailsParallel(
	ctx context.Context,
	maxWorkers int,
	emails []emailData,
) error {
	if len(emails) == 0 {
		return nil
	}

	// Create a channel to limit the number of concurrent sends
	semaphore := make(chan struct{}, maxWorkers)

	// Create a WaitGroup to wait for all sends to complete
	var wg sync.WaitGroup

	// Channel to collect errors
	errorChan := make(chan error, len(emails))

	// Function to send a single email
	sendSingleEmail := func(email emailData) {
		defer wg.Done()

		// Get a slot for sending
		semaphore <- struct{}{}
		defer func() { <-semaphore }()

		err := s.sendEmailFunc(ctx, email.toEmails, email.subject, email.body)
		if err != nil {
			errorChan <- fmt.Errorf("send email to %s: %w", strings.Join(email.toEmails, ", "), err)

			return
		}
	}

	// Start sending emails in goroutines
	for _, email := range emails {
		wg.Add(1)
		go sendSingleEmail(email)
	}

	// Wait for all sends to complete
	wg.Wait()
	close(errorChan)

	// Check for errors
	for err := range errorChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// emailData contains data for sending a single email.
type emailData struct {
	toEmails []string
	subject  string
	body     string
}

// sendEmailDirectSMTP sends an email directly using net/smtp.
func (s *Service) sendEmailDirectSMTP(
	_ context.Context,
	host string,
	port int,
	username, from string,
	toEmails []string,
	subject, bodyHTML string,
) error {
	slog.Debug("starting direct SMTP send", "host", host, "port", port, "username", username)

	// Create email message
	msg := []byte("To: " + strings.Join(toEmails, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + bodyHTML + "\r\n")

	// Connect to SMTP server
	conn, err := smtp.Dial(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		slog.Error("failed to connect to SMTP server", "host", host, "port", port, "error", err)

		return err
	}
	defer func() { _ = conn.Close() }()

	// Say hello
	if err := conn.Hello("localhost"); err != nil {
		slog.Error("failed to say hello to SMTP server", "host", host, "port", port, "error", err)

		return err
	}

	// Set sender
	if err := conn.Mail(from); err != nil {
		slog.Error("failed to set from address", "host", host, "port", port, "error", err)

		return err
	}

	// Set recipients
	for _, to := range toEmails {
		if err := conn.Rcpt(to); err != nil {
			slog.Error("failed to set to address", "host", host, "port", port, "to", to, "error", err)

			return err
		}
	}

	// Send data
	writeCloser, err := conn.Data()
	if err != nil {
		slog.Error("failed to open data connection", "host", host, "port", port, "error", err)

		return err
	}
	defer func() { _ = writeCloser.Close() }()

	if _, err := writeCloser.Write(msg); err != nil {
		slog.Error("failed to write email data", "host", host, "port", port, "error", err)

		return err
	}

	slog.Debug("email sent successfully via direct SMTP", "host", host, "port", port)

	return nil
}

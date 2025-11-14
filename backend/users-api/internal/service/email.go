package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"users-api/internal/config"
)

// EmailService define las operaciones para envío de correos electrónicos
type EmailService interface {
	SendVerificationEmail(toEmail, token string) error
	SendPasswordResetEmail(toEmail, token string) error
	GenerateToken() (string, error)
}

type emailService struct {
	config *config.Config
}

// NewEmailService crea una nueva instancia del servicio de email
func NewEmailService(cfg *config.Config) EmailService {
	return &emailService{config: cfg}
}

func (s *emailService) GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *emailService) SendVerificationEmail(toEmail, token string) error {
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", s.config.AppURL, token)

	subject := "Verifica tu correo electrónico - CarPooling"
	body := fmt.Sprintf(`
		<h2>Bienvenido a CarPooling</h2>
		<p>Por favor, verifica tu correo electrónico haciendo clic en el siguiente enlace:</p>
		<a href="%s">Verificar Email</a>
		<p>Este enlace es válido por 24 horas.</p>
	`, verificationURL)

	return s.sendEmail(toEmail, subject, body)
}

func (s *emailService) SendPasswordResetEmail(toEmail, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.config.AppURL, token)

	subject := "Restablece tu contraseña - CarPooling"
	body := fmt.Sprintf(`
		<h2>Restablecer Contraseña</h2>
		<p>Has solicitado restablecer tu contraseña. Haz clic en el siguiente enlace:</p>
		<a href="%s">Restablecer Contraseña</a>
		<p>Este enlace es válido por 1 hora.</p>
		<p>Si no solicitaste este cambio, ignora este correo.</p>
	`, resetURL)

	return s.sendEmail(toEmail, subject, body)
}

func (s *emailService) sendEmail(to, subject, body string) error {
	// Configuración SMTP
	from := s.config.SMTPFrom
	password := s.config.SMTPPassword
	smtpHost := s.config.SMTPHost
	smtpPort := s.config.SMTPPort

	// Mensaje
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", from, to, subject, body))

	// Autenticación
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Enviar email
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	return err
}

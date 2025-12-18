package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
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
	// URL para web
	webURL := fmt.Sprintf("%s/verify-email?token=%s", s.config.AppURL, token)

	// Deep link para mobile
	mobileURL := fmt.Sprintf("carpooling://verify-email?token=%s", token)

	subject := "Verifica tu correo electrónico - CarPooling"
	body := fmt.Sprintf(`
		<h2>Bienvenido a CarPooling</h2>
		<p>Por favor, verifica tu correo electrónico usando una de las siguientes opciones:</p>

		<h3>Opción 1: Desde la Web</h3>
		<p><a href="%s" style="background-color: #0ea5e9; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block;">Verificar desde Web</a></p>

		<h3>Opción 2: Desde la App Móvil</h3>
		<p><a href="%s" style="background-color: #14b8a6; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block;">Abrir en la App</a></p>

		<p style="margin-top: 20px; color: #6b7280; font-size: 14px;">
			Si el botón "Abrir en la App" no funciona, copia y pega este enlace en tu navegador móvil:<br>
			<code style="background-color: #f3f4f6; padding: 4px 8px; border-radius: 4px;">%s</code>
		</p>

		<p style="margin-top: 20px; color: #dc2626; font-size: 14px;">
			⚠️ Este enlace es válido por 24 horas.
		</p>
	`, webURL, mobileURL, mobileURL)

	return s.sendEmail(toEmail, subject, body)
}

func (s *emailService) SendPasswordResetEmail(toEmail, token string) error {
	// URL para web
	webURL := fmt.Sprintf("%s/reset-password?token=%s", s.config.AppURL, token)

	// Deep link para mobile
	mobileURL := fmt.Sprintf("carpooling://reset-password?token=%s", token)

	subject := "Restablece tu contraseña - CarPooling"
	body := fmt.Sprintf(`
		<h2>Restablecer Contraseña</h2>
		<p>Has solicitado restablecer tu contraseña. Usa una de las siguientes opciones:</p>

		<h3>Opción 1: Desde la Web</h3>
		<p><a href="%s" style="background-color: #0ea5e9; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block;">Restablecer desde Web</a></p>

		<h3>Opción 2: Desde la App Móvil</h3>
		<p><a href="%s" style="background-color: #14b8a6; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block;">Abrir en la App</a></p>

		<p style="margin-top: 20px; color: #6b7280; font-size: 14px;">
			Si el botón "Abrir en la App" no funciona, copia y pega este enlace en tu navegador móvil:<br>
			<code style="background-color: #f3f4f6; padding: 4px 8px; border-radius: 4px;">%s</code>
		</p>

		<p style="margin-top: 20px; color: #dc2626; font-size: 14px;">
			⚠️ Este enlace es válido por 1 hora.
		</p>

		<p style="margin-top: 20px; color: #6b7280; font-size: 14px;">
			Si no solicitaste este cambio, ignora este correo.
		</p>
	`, webURL, mobileURL, mobileURL)

	return s.sendEmail(toEmail, subject, body)
}

func (s *emailService) sendEmail(to, subject, body string) error {
	// Configuración SMTP
	from := s.config.SMTPFrom
	password := s.config.SMTPPassword
	smtpHost := s.config.SMTPHost
	smtpPort := s.config.SMTPPort

	// Log de intento de envío
	log.Printf("[EMAIL] Intentando enviar email a %s (Subject: %s)", to, subject)

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

	if err != nil {
		log.Printf("[EMAIL ERROR] Fallo al enviar email a %s: %v", to, err)
		return err
	}

	log.Printf("[EMAIL SUCCESS] Email enviado exitosamente a %s", to)
	return nil
}

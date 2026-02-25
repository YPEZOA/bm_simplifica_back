package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct {
	smtpHost string
	smtpPort string
	email    string
	password string
}

func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost: getEnv("SMTP_HOST", "smtp.gmail.com"),
		smtpPort: getEnv("SMTP_PORT", "587"),
		email:    getEnv("SMTP_EMAIL", "bmsimplifica@gmail.com"),
		password: getEnv("SMTP_PASSWORD", ""),
	}
}

func (e *EmailService) SendNewUserNotification(name, email, phone, message string) error {
	subject := "📧 Nuevo Usuario Solicita Registro - BM Simplifica"
	body := fmt.Sprintf(`
<h2>Nueva Solicitud de Registro</h2>

<p>Se ha recibido una nueva solicitud de registro de usuario:</p>

<ul>
	<li><strong>Nombre:</strong> %s</li>
	<li><strong>Email:</strong> %s</li>
	<li><strong>Teléfono:</strong> %s</li>
	<li><strong>Mensaje:</strong> %s</li>
</ul>

<hr>

<h3>Próximos Pasos:</h3>
<ol>
	<li>Verificar la información del usuario</li>
	<li>Crear cuenta en el sistema admin</li>
	<li>Generar contraseña segura</li>
	<li>Crear empresa asociada</li>
	<li>Subir archivos iniciales si es necesario</li>
	<li>Comunicar credenciales al cliente</li>
</ol>

<p><em>Este es un mensaje automático del sistema BM Simplifica</em></p>
	`, name, email, phone, message)

	headers := "MIME-version: 1.0\r\n"
	headers += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	headers += "From: " + e.email + "\r\n"
	headers += "To: " + e.email + "\r\n"
	headers += "Subject: " + subject + "\r\n"

	msg := headers + "\r\n" + body

	return e.sendEmail(e.email, msg)
}

func (e *EmailService) SendWelcomeEmail(userEmail, userName, tempPassword string) error {
	subject := "🎉 ¡Bienvenido a BM Simplifica!"
	body := fmt.Sprintf(`
<h2>¡Bienvenido a BM Simplifica!</h2>

<p>Estimado <strong>%s</strong>,</p>

<p>Tu cuenta ha sido creada exitosamente. Ya puedes acceder a nuestra plataforma para gestionar tus documentos empresariales.</p>

<div style="background-color: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0;">
	<h3>Tus Credenciales de Acceso:</h3>
	<ul>
		<li><strong>Email:</strong> %s</li>
		<li><strong>Contraseña Temporal:</strong> <code style="background: #e9ecef; padding: 4px 8px;">%s</code></li>
	</ul>
</div>

<p><strong>⚠️ Importante:</strong> Por seguridad, te recomendamos cambiar tu contraseña en tu primer inicio de sesión.</p>

<h3>¿Qué puedes hacer en la plataforma?</h3>
<ul>
	<li>Ver todas tus empresas registradas</li>
	<li>Explorar los documentos que hemos preparado para ti</li>
	<li>Descargar archivos cuando los necesites</li>
</ul>

<p>Para cualquier duda o soporte, no dudes en contactarnos.</p>

<hr>

<p>Atentamente,<br>El equipo de BM Simplifica<br>📧 bmsimplifica@gmail.com</p>
	`, userName, userEmail, tempPassword)

	headers := "MIME-version: 1.0\r\n"
	headers += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	headers += "From: " + e.email + "\r\n"
	headers += "To: " + userEmail + "\r\n"
	headers += "Subject: " + subject + "\r\n"

	msg := headers + "\r\n" + body

	return e.sendEmail(userEmail, msg)
}

func (e *EmailService) sendEmail(to, msg string) error {
	auth := smtp.PlainAuth("", e.email, e.password, e.smtpHost)

	// Connect to server without TLS first
	client, err := smtp.Dial(e.smtpHost + ":" + e.smtpPort)
	if err != nil {
		return fmt.Errorf("error al conectar al servidor SMTP: %v", err)
	}
	defer client.Quit()

	// Check if server supports STARTTLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         e.smtpHost,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("error al iniciar STARTTLS: %v", err)
		}
	}

	// Authenticate
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("error al autenticar: %v", err)
	}

	// Set sender and recipient
	if err := client.Mail(e.email); err != nil {
		return fmt.Errorf("error al establecer remitente: %v", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("error al establecer destinatario: %v", err)
	}

	// Send message
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("error al enviar datos: %v", err)
	}
	defer w.Close()

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("error al escribir mensaje: %v", err)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Package email contains a minimal snippet for sending an HTML email with attachments.
//
// Notes:
// - This uses gomail (gopkg.in/gomail.v2).
// - Attachment filename is RFC 2047 encoded to avoid garbled Chinese filenames
//   in some clients (macOS Mail / iOS Mail).
package email

import (
	"encoding/base64"
	"fmt"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

// Config holds SMTP credentials.
type Config struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

// SendHTMLWithAttachments sends an HTML email with optional attachments.
func SendHTMLWithAttachments(cfg Config, subject string, to, cc []string, htmlBody string, attachments []string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", cfg.From)
	if len(to) > 0 {
		m.SetHeader("To", to...)
	}
	if len(cc) > 0 {
		m.SetHeader("Cc", cc...)
	}
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	for _, p := range attachments {
		// Fix Chinese filename garbling in some mail clients.
		name := filepath.Base(p)
		encoded := fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(name)))
		m.Attach(p, gomail.Rename(encoded))
	}

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Pass)
	return d.DialAndSend(m)
}

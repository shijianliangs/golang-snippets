# Send email (HTML + attachments)

This snippet sends an HTML email with optional attachments.

- Library: `gopkg.in/gomail.v2`
- Includes RFC 2047 encoding for attachment file names (helps with Chinese filenames on some clients).

## Install

```bash
go get gopkg.in/gomail.v2
```

## Usage

```go
cfg := email.Config{
  Host: "smtp.example.com",
  Port: 465,
  User: "your_user",
  Pass: "your_pass",
  From: "Your Name <your_user@example.com>",
}

err := email.SendHTMLWithAttachments(
  cfg,
  "Subject",
  []string{"to@example.com"},
  nil,
  "<h1>Hello</h1>",
  []string{"./report.html", "./截图.png"},
)
```

## Notes
- For SSL/TLS ports (like 465), gomail handles it via the dialer.
- Keep secrets out of code. Use env vars or secret managers.

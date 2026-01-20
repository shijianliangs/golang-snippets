package notification

import (
	"encoding/base64"
	"fmt"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

// SendEmail 发送带附件的HTML格式邮件
func SendEmail(subject string, to []string, cc []string, body string, attachment []string) error {

	
	m := gomail.NewMessage()

	// 设置发件人和收件人
	m.SetHeader("From", "<FromUser>")
	m.SetHeader("To", to...)
	m.SetHeader("Cc", cc...)

	// 设置邮件主题
	m.SetHeader("Subject", subject)

	// 设置邮件正文
	m.SetBody("text/html", body)

	// 添加附件，指定文件路径
	for _, v := range attachment {
    // 这里解决附件中文乱码问题(mac自带邮件客户端+ios自带邮件客户端)
		filename := fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(filepath.Base(v))))
		m.Attach(
			v,
			gomail.Rename(filename),
			// gomail.SetHeader(map[string][]string{
			// 	"Content-Disposition": {
			// 		fmt.Sprintf(`attachment; filename="%s"`, filename),
			// 	},
			// }),
		)
	}

	// 创建一个新的SMTP客户端并发送邮件
	dialer := gomail.NewDialer("<EmailHost>", "<EmailPort>", "<EmailUser>", "<EmailPassword>")

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

package config

import "gopkg.in/gomail.v2"

type Mail struct {
	MailHost string
	MailPort int
	MailUser string
	MailPwd  string
}

func SendGoMail(mailAddress []string, subject string, body string, nickname string, mail Mail) error {

	m := gomail.NewMessage()
	// 这种方式可以添加别名，即 nickname， 也可以直接用<code>m.SetHeader("From", MAIL_USER)</code>
	m.SetHeader("From", nickname)
	// 发送给多个用户
	m.SetHeader("To", mailAddress...)
	// 设置邮件主题
	m.SetHeader("Subject", subject)
	// 设置邮件正文
	m.SetBody("text/html", body)
	d := gomail.NewDialer(mail.MailHost, mail.MailPort, mail.MailUser, mail.MailPwd)
	// 发送邮件
	err := d.DialAndSend(m)
	return err
}

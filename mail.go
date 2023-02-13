package main

import (
	"crypto/tls"
	"log"

	"gopkg.in/gomail.v2"
)

type mailer struct {
	MailServer string
	Port       int
	Username   string
	Password   string
}

func (m *mailer) dialer() *gomail.Dialer {
	d := gomail.NewDialer(
		m.MailServer,
		m.Port,
		m.Username,
		m.Password,
	)
	d.TLSConfig = &tls.Config{
		ServerName:         m.MailServer,
		InsecureSkipVerify: false,
	}
	return d
}

func (m *mailer) sendMail(mailToName, mailToAddress, mailFromName, mailFromAddress, subject, template string, attachmentPath string) error {
	d := m.dialer()
	s, err := d.Dial()
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", mailFromAddress, mailFromName)
	msg.SetAddressHeader("To", mailToAddress, mailToName)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", template)
	if attachmentPath != "" {
		msg.Attach(attachmentPath)
	}

	if err := gomail.Send(s, msg); err != nil {
		log.Printf("Could not send email to %q: %v", mailToAddress, err)
		return err
	}

	return nil
}

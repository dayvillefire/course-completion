package main

import "testing"

func Test_Mailer(t *testing.T) {
	c, err := loadConfig("config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	m := mailer{
		MailServer: c.Mail.ServerName,
		Port:       c.Mail.ServerPort,
		Username:   c.Mail.Username,
		Password:   c.Mail.Password,
	}
	err = m.sendMail(
		"Me",
		c.Mail.Username,
		"Also me",
		c.Mail.Username,
		"Testing course completion emails",
		"<b>This</b> is a test.<br/><br/>Signed,</br/>Me.",
		"images/Star_of_life.png",
	)
	if err != nil {
		t.Fatal(err)
	}
}

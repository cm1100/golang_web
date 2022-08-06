package main

import (
	"fmt"
	"myapp2/pkg/models"

	//"html/template"
	"io/ioutil"
	"log"
	"strings"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {

	go func() {

		for {

			msg := <-app.MailChan
			sendMessage(msg)

		}
	}()

}

func sendMessage(m models.MailData) {

	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			log.Println(err)
			return
		}
		mailTemplate := string(data)
		msgToSpend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSpend)

	}
	err = email.Send(client)

	if err != nil {
		log.Println(err)
	} else {

		log.Println("email sent")
	}

}

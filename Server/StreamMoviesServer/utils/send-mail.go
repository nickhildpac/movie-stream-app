package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nickhildpac/movie-stream-app/Server/StreamMoviesServer/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func ListenForMail(mailChan chan models.MailData) {
	go func() {
		for {
			msg := <-mailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		log.Println(err)
	}
	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	log.Println("Sending email to: ", m.To)
	if m.Template == "" {
		log.Println("from plain text", m.To)
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		log.Println("from template", m.To)
		data, err := os.ReadFile(fmt.Sprintf("./templates/%s", m.Template))
		if err != nil {
			log.Println(err.Error())
		}
		mailTemplate := string(data)
		log.Println(mailTemplate)
		msgToSend := strings.Replace(mailTemplate, "{{reset_link}}", m.Content, 3)
		email.SetBody(mail.TextHTML, msgToSend)
	}
	err = email.Send(client)
	if err != nil {
		log.Println(err)
	}
}

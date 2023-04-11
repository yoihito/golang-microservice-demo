package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"time"

	"notification/config"
	"notification/internal/infrustructure"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	consumer := infrustructure.NewRabbitMqConsumer(config.RabbitMqUrl, []infrustructure.RabbitMqQueue{{config.AudioQueue}})

	for {
		if consumer.IsReady() {
			break
		}
		<-time.After(10 * time.Second)
	}

	msgs, err := consumer.Consume(config.AudioQueue, false)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles("templates/mail.html")
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			log.Printf("Received Message: %s\n", msg.Body())
			if err := processEvent(config, tmpl, msg); err != nil {
				log.Println(err)
				msg.Nack()
				continue
			}
			msg.Ack()
		}
	}()
	<-forever
}

type AudioExtractedEvent struct {
	ObjectId         string `json:"objectId"`
	OriginalFilename string `json:"originalFilename"`
	Email            string `json:"email"`
}

func processEvent(config config.Config, tmpl *template.Template, msg infrustructure.Delivery) error {
	var event AudioExtractedEvent
	if err := json.Unmarshal(msg.Body(), &event); err != nil {
		return err
	}

	data := NewTemplateData(config, event)
	message, err := GenerateMessage(tmpl, &data)
	if err != nil {
		return err
	}

	if err = SendMail(config.SmtpHost+":"+config.SmtpPort, config.FromEmail, event.Email, message); err != nil {
		return err
	}
	log.Printf("Received Event: %s\n", event)

	return nil
}

type TemplateData struct {
	Subject      string
	DownloadLink string
}

func NewTemplateData(config config.Config, event AudioExtractedEvent) TemplateData {
	link := fmt.Sprintf("%s/download/%s", config.DownloadHost, event.ObjectId)
	return TemplateData{
		Subject:      "Audio is extracted!",
		DownloadLink: link,
	}
}

func GenerateMessage(tmpl *template.Template, data *TemplateData) (*bytes.Buffer, error) {
	buffer := bytes.NewBufferString(fmt.Sprintf("Subject: %s\r\nContent-Type: text/html; charset=utf-8\r\n\r\n", data.Subject))
	if err := tmpl.Execute(buffer, data); err != nil {
		return nil, err
	}
	buffer.WriteString("\r\n")
	return buffer, nil
}

func SendMail(host, fromEmail, toEmail string, buf *bytes.Buffer) (err error) {
	c, err := smtp.Dial(host)
	if err != nil {
		return err
	}

	defer func() {
		tempErr := c.Close()
		if err == nil {
			err = tempErr
		}
	}()

	c.Mail(fromEmail)
	c.Rcpt(toEmail)

	wc, err := c.Data()
	if err != nil {
		return err
	}

	defer func() {
		tempErr := wc.Close()
		if err == nil {
			err = tempErr
		}
	}()

	if _, err = buf.WriteTo(wc); err != nil {
		return err
	}

	return nil
}

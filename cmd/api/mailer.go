package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/xhit/go-simple-mail/v2"
)

//go:embed templates
var emailTemplateFS embed.FS

func (app *application) sendEmail(from, to, subject, tmpl string, data interface{}) error {
	templateToRender := fmt.Sprintf("templates/%s.html.tmpl", tmpl)

	// "email-html" is the template name, but i could have choice another name
	t, err := template.New("email-html").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		app.errorLog.Println(err)
		return err
	}

	// cast the template to string
	// into tpl is the html version of the template
	formattedMessage := tpl.String()

	// now we need to do the plainText
	templateToRender = fmt.Sprintf("templates/%s.plain.tmpl", tmpl)
	t, err = template.New("email-plan").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}
	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		app.errorLog.Println(err)
		return err
	}

	plainMessage := tpl.String()

	app.infoLog.Println(formattedMessage, plainMessage)

	// send the mail
	server := mail.NewSMTPClient()
	server.Host = app.config.smtp.host
	server.Port = app.config.smtp.port
	server.Username = app.config.smtp.username
	server.Password = app.config.smtp.password
	server.Encryption = mail.EncryptionNone
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject)

	// mail.TextHTML and mail.TextPlain are const built in the mail package
	email.SetBody(mail.TextHTML, formattedMessage)
	email.AddAlternative(mail.TextPlain, plainMessage)

	err = email.Send(smtpClient)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	app.infoLog.Println("Mail sent")

	return nil
}

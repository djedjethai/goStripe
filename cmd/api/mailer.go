package main

import (
	"embed"
	"fmt"

	"github.com/xhit/go-simple-mail/v2"
)

//go:embed templates
var emailTemplateFS embed.FS

func (app *application) sendEmail(from, to, subject, tmpl string, data interface{}) error {
	templateToRender := fmt.Sprintf("templates/%s.html.tmpl", tmpl)

	// "email-html" is the template name, but i could have choice another name
	t, err := template.new("email-html").ParseFS(emailTemplateFS, templateToRender)
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
	templateToRender := fmt.Sprintf("template/%s.plain.tmpl", tmpl)

	// TOTOTOTO FINISHHHH .....
	t, err := template.new("email-html").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	return nil
}

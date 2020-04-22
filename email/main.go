package email

import (
	"os"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Email struct {
	to   *mail.Email
	from *mail.Email
}

type Hit struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func NewEmail() Email {
	return Email{
		to:   mail.NewEmail(os.Getenv("TO_EMAIL_NAME"), os.Getenv("TO_EMAIL")),
		from: mail.NewEmail(os.Getenv("FROM_EMAIL_NAME"), os.Getenv("FROM_EMAIL")),
	}
}

func (e Email) Send(hits []Hit) (*rest.Response, error) {
	m := mail.NewV3Mail()
	m.SetFrom(e.from)

	m.SetTemplateID(os.Getenv("EMAIL_TEMPLATE"))

	p := mail.NewPersonalization()
	p.AddTos(e.to)

	p.SetDynamicTemplateData("hits", hits)
	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"

	var Body = mail.GetRequestBody(m)
	request.Body = Body
	return sendgrid.API(request)
}

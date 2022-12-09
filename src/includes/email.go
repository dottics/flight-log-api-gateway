package includes

import (
	"fmt"
	"github.com/dottics/dutil"
	"github.com/dottics/emailserv"
	"github.com/google/uuid"
	"html/template"
	"net/mail"
	"net/url"
	"os"
	"path"
)

func loadTemplate() *template.Template {
	wd := os.Getenv("WORKDIR")
	templates := path.Join(wd, "templates/*.html")
	tpl := template.Must(template.New("").ParseGlob(templates))
	return tpl
}

type ForgotPasswordData struct {
	Token                   uuid.UUID
	ResetPasswordLink       string
	RevokeResetPasswordLink string
	HomeLink                string
	ContactUsLink           string
}

// NewForgotPasswordData gets all the basic forgot password email body
// information from the environment and returns a new instance of the
// forgot password email body to populate the email template.
func NewForgotPasswordData(t uuid.UUID) *ForgotPasswordData {
	scheme := os.Getenv("APP_SCHEME")
	host := os.Getenv("APP_HOST")
	u := url.URL{
		Scheme: scheme,
		Host:   host,
	}
	q := url.Values{
		"r": []string{t.String()},
	}
	m := &ForgotPasswordData{}
	m.Token = t
	// Reset password
	u.Path = "/reset-password"
	u.RawQuery = q.Encode()
	m.ResetPasswordLink = u.String()
	// Revoke password
	u.Path = "/revoke-password"
	u.RawQuery = q.Encode()
	m.RevokeResetPasswordLink = u.String()
	// Home
	u.Path = "/"
	u.RawQuery = ""
	m.HomeLink = u.String()
	// Contact us
	u.Path = "/contact-us"
	u.RawQuery = ""
	m.ContactUsLink = u.String()
	return m
}

type ForgotPasswordMsg struct {
	Message *emailserv.Message
	Data    *ForgotPasswordData
}

// NewForgotPasswordMsg does the basic scaffolding and data manipulation
// for the forgot password email.
func NewForgotPasswordMsg(to mail.Address, t uuid.UUID) *ForgotPasswordMsg {
	msg := &ForgotPasswordMsg{
		Message: &emailserv.Message{
			Headers: map[string][]string{
				"Mime-Version": {"1.0"},
				"Content-Type": {"text/html", "charset=UTF-8"},
			},
			From: mail.Address{
				Name:    "No-Reply Dottics",
				Address: "mail@dottics.com",
			},
			To: []mail.Address{to},
			ReplyTo: mail.Address{
				Name:    "Johannes Scribante",
				Address: "js@dottics.com",
			},
			Subject: "Dottics Forgot Password",
		},
		Data: NewForgotPasswordData(t),
	}
	return msg
}

// ExecuteTemplate creates an HTML file for the forgot password email. It
// executes the template to populate the template with data. Then it reads
// the HTML file and converts it to a string for the emailserv.Message
// body.
func (msg *ForgotPasswordMsg) ExecuteTemplate() dutil.Error {
	tpl := loadTemplate()
	wd := os.Getenv("WORKDIR")
	u, err := uuid.NewUUID()
	if err != nil {
		e := dutil.NewErr(500, "uuid", []string{"unable to generate uuid", err.Error()})
		return e
	}

	name := path.Join(wd, "documents", fmt.Sprintf("forgot-password-%s.html", u.String()))
	file, err := os.Create(name)
	if err != nil {
		e := dutil.NewErr(500, "document", []string{"unable to create document", err.Error()})
		return e
	}

	err = tpl.ExecuteTemplate(file, "forgot-password.html", msg.Data)
	if err != nil {
		e := dutil.NewErr(500, "template", []string{"unable to execute template", err.Error()})
		return e
	}

	xb, err := os.ReadFile(name)
	if err != nil {
		e := dutil.NewErr(500, "readFile", []string{"unable to read file", err.Error()})
		return e
	}

	msg.Message.Body = string(xb)
	//log.Printf("\n%s\n", msg.Message.Body)
	return nil
}

// SendForgotPassword handles the execution of the forgot password HTML
// template and then sends the email.
func SendForgotPassword(msg ForgotPasswordMsg) dutil.Error {
	// no token required for the email microservice at the moment (2022-04-17)
	s := emailserv.NewService("")
	return s.SendMail(msg.Message)
}

type ContactUsData struct {
	Name          string
	Email         string
	Message       string
	HomeLink      string
	ContactUsLink string
}
type ContactUsMsg struct {
	Message *emailserv.Message
	Data    *ContactUsData
}

func NewContactUsMsg(replyTo mail.Address, message string) *ContactUsMsg {
	scheme := os.Getenv("APP_SCHEME")
	host := os.Getenv("APP_HOST")
	u := url.URL{
		Scheme: scheme,
		Host:   host,
	}
	d := &ContactUsData{
		Name:    replyTo.Name,
		Email:   replyTo.Address,
		Message: message,
	}
	// Home
	u.Path = "/"
	u.RawQuery = ""
	d.HomeLink = u.String()
	// Contact us
	u.Path = "/contact-us"
	u.RawQuery = ""
	d.ContactUsLink = u.String()

	return &ContactUsMsg{
		Message: &emailserv.Message{
			Headers: map[string][]string{
				"Mime-Version": {"1.0"},
				"Content-Type": {"text/html", "charset=UTF-8"},
			},
			From: mail.Address{
				Name:    "Dottics Contact Us",
				Address: "mail@dottics.com",
			},
			To: []mail.Address{
				{Name: "Dottics Team", Address: "howzit@dottics.com"},
			},
			CC:      []mail.Address{replyTo},
			ReplyTo: replyTo,
			Subject: "Dottics Contact Us",
		},
		Data: d,
	}
}

func (msg *ContactUsMsg) ExecuteTemplate() dutil.Error {
	tpl := loadTemplate()
	wd := os.Getenv("WORKDIR")
	u, err := uuid.NewUUID()
	if err != nil {
		e := dutil.NewErr(500, "uuid", []string{"unable to generate uuid", err.Error()})
		return e
	}

	name := path.Join(wd, "documents", fmt.Sprintf("contact-us-%s.html", u.String()))
	file, err := os.Create(name)
	if err != nil {
		e := dutil.NewErr(500, "document", []string{"unable to create document", err.Error()})
		return e
	}

	err = tpl.ExecuteTemplate(file, "contact-us.html", msg.Data)
	if err != nil {
		e := dutil.NewErr(500, "template", []string{"unable to execute template", err.Error()})
		return e
	}

	xb, err := os.ReadFile(name)
	if err != nil {
		e := dutil.NewErr(500, "readFile", []string{"unable to read file", err.Error()})
		return e
	}

	msg.Message.Body = string(xb)
	return nil
}

// SendMail will become a generic function to send an email via the
// emailserv microservice package to the email microservice.
func (msg *ContactUsMsg) SendMail() dutil.Error {
	ms := emailserv.NewService("")
	return ms.SendMail(msg.Message)
}

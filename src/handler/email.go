package handler

import (
	"github.com/dottics/dutil"
	"github.com/dottics/flight-log-api-gateway/src/includes"
	security "github.com/dottics/securityserv"
	"net/http"
	"net/mail"
)

// ForgotPassword handles the generation of the forgot password email
// and exchanges with the email microservice to send the email.
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	// decode the request json body
	p := security.PasswordResetTokenPayload{}
	e := dutil.Decode(w, r, &p)
	if e != nil {
		Error(w, r, e)
		return
	}

	// generate the password reset token
	t, e := includes.PasswordResetToken(p)
	if e != nil {
		Error(w, r, e)
		return
	}

	to := mail.Address{Address: p.Email}
	msg := includes.NewForgotPasswordMsg(to, t)
	e = msg.ExecuteTemplate()
	if e != nil {
		Error(w, r, e)
		return
	}

	e = includes.SendForgotPassword(*msg)
	if e != nil {
		Error(w, r, e)
		return
	}

	resp := dutil.Resp{
		Status:  200,
		Message: "forgot password email sent successfully",
	}
	resp.Respond(w, r)
}

// ContactUs handles the generation of the contact-us email and exchanges
// with the email microservice to send the contact-us email.
func ContactUs(w http.ResponseWriter, r *http.Request) {
	msgData := includes.ContactUsMsgPayload{}
	e := dutil.Decode(w, r, &msgData)
	if e != nil {
		Error(w, r, e)
		return
	}

	replyTo := mail.Address{
		Name:    msgData.Name,
		Address: msgData.Email,
	}
	msg := includes.NewContactUsMsg(replyTo, msgData.Message)
	e = msg.ExecuteTemplate()
	if e != nil {
		Error(w, r, e)
		return
	}
	e = msg.SendMail()
	if e != nil {
		Error(w, r, e)
		return
	}

	resp := dutil.Resp{
		Status:  200,
		Message: "contact-us email sent successfully",
	}
	resp.Respond(w, r)
}

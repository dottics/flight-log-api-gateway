package includes

import (
	"fmt"
	"github.com/dottics/dutil"
	"github.com/dottics/emailserv"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/microtest"
	"net/mail"
	"os"
	"path"
	"strings"
	"testing"
)

func TestNewForgotPasswordData(t *testing.T) {
	wd, _ := os.Getwd()
	_ = os.Setenv("WORKDIR", path.Join(wd, "../.."))
	_ = os.Setenv("APP_SCHEME", "https")
	_ = os.Setenv("APP_HOST", "test.dottics.com")

	m := NewForgotPasswordData(uuid.MustParse("73730848-a9ed-4d25-9892-7948799cdc7a"))

	resetPasswordLink := "https://test.dottics.com/reset-password?r=73730848-a9ed-4d25-9892-7948799cdc7a"
	if m.ResetPasswordLink != resetPasswordLink {
		t.Errorf("expected %s got %s", resetPasswordLink, m.ResetPasswordLink)
	}
	revokePasswordLink := "https://test.dottics.com/revoke-password?r=73730848-a9ed-4d25-9892-7948799cdc7a"
	if m.RevokeResetPasswordLink != revokePasswordLink {
		t.Errorf("expected %s got %s", revokePasswordLink, m.RevokeResetPasswordLink)
	}
	homeLink := "https://test.dottics.com/"
	if m.HomeLink != homeLink {
		t.Errorf("expected %s got %s", homeLink, m.HomeLink)
	}
	contactUsLink := "https://test.dottics.com/contact-us"
	if m.ContactUsLink != contactUsLink {
		t.Errorf("expected %s got %s", contactUsLink, m.ContactUsLink)
	}
}

func TestNewForgotPasswordMsg(t *testing.T) {
	wd, _ := os.Getwd()
	_ = os.Setenv("WORKDIR", path.Join(wd, "../.."))
	_ = os.Setenv("APP_SCHEME", "https")
	_ = os.Setenv("APP_HOST", "test.dottics.com")
	to := mail.Address{
		Name:    "James Bond",
		Address: "james@bond.com",
	}
	token := uuid.MustParse("73730848-a9ed-4d25-9892-7948799cdc7a")

	msg := NewForgotPasswordMsg(to, token)
	if msg.Message.To[0] != to {
		t.Errorf("expected to address %v got %v", to, msg.Message.To[0])
	}
	if msg.Data.Token != token {
		t.Errorf("expected token %s got %s", token, msg.Data.Token)
	}
	resetPasswordLink := "https://test.dottics.com/reset-password?r=73730848-a9ed-4d25-9892-7948799cdc7a"
	if msg.Data.ResetPasswordLink != resetPasswordLink {
		t.Errorf("expected token %s got %s", resetPasswordLink, msg.Data.ResetPasswordLink)
	}
}

func TestForgotPasswordMsg_ExecuteTemplate(t *testing.T) {
	wd, _ := os.Getwd()
	_ = os.Setenv("WORKDIR", path.Join(wd, "../.."))
	_ = os.Setenv("APP_SCHEME", "https")
	_ = os.Setenv("APP_HOST", "test.dottics.com")
	to := mail.Address{
		Name:    "James Bond",
		Address: "james@bond.com",
	}
	token := uuid.MustParse("73730848-a9ed-4d25-9892-7948799cdc7a")
	msg := NewForgotPasswordMsg(to, token)
	e := msg.ExecuteTemplate()
	if e != nil {
		t.Errorf("expected error %v got %v", nil, e)
	}
	if len(msg.Message.Body) == 0 {
		t.Errorf("expected body length %d > 0", len(msg.Message.Body))
	}
	if !strings.Contains(msg.Message.Body, msg.Data.ResetPasswordLink) {
		t.Errorf("expected template to execute and contain %s", msg.Data.ResetPasswordLink)
	}
}

func TestSendForgotPassword(t *testing.T) {
	wd, _ := os.Getwd()
	_ = os.Setenv("WORKDIR", path.Join(wd, "../.."))
	tests := []struct {
		name     string
		msg      ForgotPasswordMsg
		exchange *microtest.Exchange
		e        dutil.Error
	}{
		{
			name: "Validation Error",
			msg: ForgotPasswordMsg{
				Message: &emailserv.Message{
					Subject: "Missing To",
				},
			},
			exchange: nil,
			e: &dutil.Err{
				Status: 400,
				Errors: map[string][]string{
					"body":    {"required"},
					"from":    {"address required"},
					"replyTo": {"address required"},
					"to":      {"minimum 1 address"},
				},
			},
		},
		{
			name: "Bad Request",
			msg: ForgotPasswordMsg{
				Message: &emailserv.Message{
					From:    mail.Address{Name: "Test User", Address: "test@dottics.com"},
					To:      []mail.Address{{Name: "another user", Address: "another@dottics.com"}},
					ReplyTo: mail.Address{Name: "Reply-To", Address: "reply-to@dottics.com"},
					Subject: "Test message",
					Body:    `<html><body>here is my email</body></html>`,
				},
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   `{"message":"BadRequest","data":null,"errors":{"permission":["please ensure you have permission"]}}`,
				},
			},
			e: &dutil.Err{
				Status: 403,
				Errors: map[string][]string{
					"permission": {"please ensure you have permission"},
				},
			},
		},
		{
			name: "Successful",
			msg: ForgotPasswordMsg{
				Message: &emailserv.Message{
					From:    mail.Address{Name: "Test User", Address: "test@dottics.com"},
					To:      []mail.Address{{Name: "another user", Address: "another@dottics.com"}},
					ReplyTo: mail.Address{Name: "Reply-To", Address: "reply-to@dottics.com"},
					Subject: "Test message",
					Body:    `<html><body>here is my email</body></html>`,
				},
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   `{"message":"email send successful","data":null,"errors":null}`,
				},
			},
			e: nil,
		},
	}

	ms := microtest.NewMockServer("EMAIL_SERVICE_SCHEME", "EMAIL_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tests {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			ms.Append(tc.exchange)

			e := SendForgotPassword(tc.msg)
			if !dutil.ErrorEqual(e, tc.e) {
				t.Errorf("expected error %v got %v", tc.e, e)
			}
		})
	}
}

func TestNewContactUsMsg(t *testing.T) {
	wd, _ := os.Getwd()
	_ = os.Setenv("WORKDIR", path.Join(wd, "../.."))
	_ = os.Setenv("APP_SCHEME", "https")
	_ = os.Setenv("APP_HOST", "test.dottics.com")
	to := mail.Address{
		Name:    "James Bond",
		Address: "james@bond.com",
	}
	message := "Hi there,\n\nCan I please be in the closed review for the budget app\n\n" +
		"I am still new to budgeting.\n\nRegards\nJames Bond"

	msg := NewContactUsMsg(to, message)

	if msg.Data.HomeLink != "https://test.dottics.com/" {
		t.Errorf("expected home link '%s' got '%s'", "https://test.dottics.com/", msg.Data.HomeLink)
	}
	if msg.Data.Name != to.Name {
		t.Errorf("expected name '%s' got '%s'", to.Name, msg.Data.Name)
	}
	if msg.Data.Email != to.Address {
		t.Errorf("expected name '%s' got '%s'", to.Address, msg.Data.Email)
	}
	if msg.Data.Message != message {
		t.Errorf("expected message '%s' got %s'", message, msg.Data.Message)
	}
}

func TestContactUsMsg_ExecuteTemplate(t *testing.T) {
	wd, _ := os.Getwd()
	_ = os.Setenv("WORKDIR", path.Join(wd, "../.."))
	_ = os.Setenv("APP_SCHEME", "https")
	_ = os.Setenv("APP_HOST", "test.dottics.com")
	to := mail.Address{
		Name:    "James Bond",
		Address: "james@bond.com",
	}
	message := "Hi there,\n\nCan I please be in the closed review for the budget app\n\n" +
		"I am still new to budgeting.\n\nRegards\nJames Bond"

	msg := NewContactUsMsg(to, message)
	e := msg.ExecuteTemplate()
	if e != nil {
		t.Errorf("expected error %v got %v", nil, e)
	}
	if len(msg.Message.Body) == 0 {
		t.Errorf("expected body length %d > 0", len(msg.Message.Body))
	}
	if !strings.Contains(msg.Message.Body, msg.Data.Message) {
		t.Errorf("expected template to execute and contain %s", msg.Data.Message)
	}
}

func TestContactUsMsg_SendMail(t *testing.T) {
	wd, _ := os.Getwd()
	_ = os.Setenv("WORKDIR", path.Join(wd, "../.."))
	tests := []struct {
		name     string
		msg      ContactUsMsg
		exchange *microtest.Exchange
		e        dutil.Error
	}{
		{
			name: "Bad Request",
			msg: ContactUsMsg{
				Message: &emailserv.Message{
					From:    mail.Address{Name: "Test User", Address: "test@dottics.com"},
					To:      []mail.Address{{Name: "another user", Address: "another@dottics.com"}},
					ReplyTo: mail.Address{Name: "Reply-To", Address: "reply-to@dottics.com"},
					Subject: "Test message",
					Body:    `<html><body>here is my email</body></html>`,
				},
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   `{"message":"BadRequest","data":null,"errors":{"permission":["please ensure you have permission"]}}`,
				},
			},
			e: &dutil.Err{
				Status: 403,
				Errors: map[string][]string{
					"permission": {"please ensure you have permission"},
				},
			},
		},
		{
			name: "Successful",
			msg: ContactUsMsg{
				Message: &emailserv.Message{
					From:    mail.Address{Name: "Test User", Address: "test@dottics.com"},
					To:      []mail.Address{{Name: "another user", Address: "another@dottics.com"}},
					ReplyTo: mail.Address{Name: "Reply-To", Address: "reply-to@dottics.com"},
					Subject: "Test message",
					Body:    `<html><body>here is my email</body></html>`,
				},
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   `{"message":"email send successful","data":null,"errors":null}`,
				},
			},
			e: nil,
		},
	}

	ms := microtest.NewMockServer("EMAIL_SERVICE_SCHEME", "EMAIL_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tests {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			ms.Append(tc.exchange)

			e := tc.msg.SendMail()
			if !dutil.ErrorEqual(e, tc.e) {
				t.Errorf("expected error %v got %v", tc.e, e)
			}
		})
	}
}

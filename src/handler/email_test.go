package handler

import (
	"bytes"
	"fmt"
	"github.com/johannesscr/micro/microtest"
	"io"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
)

func TestForgotPassword(t *testing.T) {
	wd, _ := os.Getwd()
	_ = os.Setenv("WORKDIR", path.Join(wd, "../.."))

	type E struct {
		status int
		data   string
	}
	tests := []struct {
		name          string
		payload       io.Reader
		secExchange   *microtest.Exchange
		emailExchange *microtest.Exchange
		E             E
	}{
		{
			name:          "Payload Decode Error",
			payload:       strings.NewReader(`{"email:"name@example.com"}`),
			secExchange:   nil,
			emailExchange: nil,
			E: E{
				status: 400,
				data:   `{"message":"Bad Request","data":null,"errors":{"decode":["unable to decode data","invalid character 'n' after object key"]}}`,
			},
		},
		{
			name:    "bad request password reset token",
			payload: strings.NewReader(`{"email":"name@example.com"}`),
			secExchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 400,
					Body:   `{"message":"BadRequest","data":{},"errors":{"user":["not found"]}}`,
				},
			},
			emailExchange: nil,
			E: E{
				status: 400,
				data:   `{"message":"Bad Request","data":null,"errors":{"user":["not found"]}}`,
			},
		},
		{
			name:    "bad request email could not send",
			payload: strings.NewReader(`{"email":"name@example.com"}`),
			secExchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   `{"message":"password reset token successful","data":{"password_reset_token":"93142963-531d-4ca2-8b78-b5dc61f48c04"},"errors":{}}`,
				},
			},
			emailExchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 400,
					Body:   `{"message":"Bad Request","data":null,"errors":{"email":["sender not found"]}}`,
				},
			},
			E: E{
				status: 400,
				data:   `{"message":"Bad Request","data":null,"errors":{"email":["sender not found"]}}`,
			},
		},
		{
			name:    "successful",
			payload: strings.NewReader(`{"email":"name@example.com"}`),
			secExchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   `{"message":"password reset token successful","data":{"password_reset_token":"93142963-531d-4ca2-8b78-b5dc61f48c04"},"errors":{}}`,
				},
			},
			emailExchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   `{"message":"email send successful","data":null,"errors":null}`,
				},
			},
			E: E{
				status: 200,
				data:   `{"message":"forgot password email sent successfully","data":null,"errors":null}`,
			},
		},
	}

	securityMS := microtest.NewMockServer("SECURITY_SERVICE_SCHEME", "SECURITY_SERVICE_HOST")
	defer securityMS.Server.Close()
	emailMS := microtest.NewMockServer("EMAIL_SERVICE_SCHEME", "EMAIL_SERVICE_HOST")
	defer emailMS.Server.Close()

	for i, tc := range tests {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			securityMS.Append(tc.secExchange)
			emailMS.Append(tc.emailExchange)

			req := microtest.NewRequest("post", "/forgot-password", nil, nil, tc.payload)
			rec := httptest.NewRecorder()
			ForgotPassword(rec, req)
			res, xb := microtest.ReadRecorder(rec)

			if res.StatusCode != tc.E.status {
				t.Errorf("expected status code %d got %d", tc.E.status, res.StatusCode)
			}
			d := string(bytes.TrimSpace(xb))
			if d != tc.E.data {
				t.Errorf("expected data '%s' got '%s'", tc.E.data, d)
			}
		})
	}
}

func TestContactUs(t *testing.T) {
	wd, _ := os.Getwd()
	_ = os.Setenv("WORKDIR", path.Join(wd, "../.."))

	type E struct {
		status int
		data   string
	}
	tests := []struct {
		name     string
		payload  io.Reader
		exchange *microtest.Exchange
		E        E
	}{
		{
			name:     "bad request",
			payload:  strings.NewReader(`{}`),
			exchange: nil,
			E: E{
				status: 400,
				data:   `{"message":"Bad Request","data":null,"errors":{"replyTo":["address required"]}}`,
			},
		},
		{
			name:    "successful",
			payload: strings.NewReader(`{"name":"James Bond", "email":"name@example.com","message":"here is\nmy message"}`),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   `{"message":"message sent successfully","data":{},"errors":{}}`,
				},
			},
			E: E{
				status: 200,
				data:   `{"message":"contact-us email sent successfully","data":null,"errors":null}`,
			},
		},
	}

	ms := microtest.NewMockServer("EMAIL_SERVICE_SCHEME", "EMAIL_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tests {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			ms.Append(tc.exchange)

			req := microtest.NewRequest("post", "/contact-us", nil, nil, tc.payload)
			rec := httptest.NewRecorder()
			ContactUs(rec, req)
			res, xb := microtest.ReadRecorder(rec)

			if res.StatusCode != tc.E.status {
				t.Errorf("expected status code %d got %d", tc.E.status, res.StatusCode)
			}
			d := string(bytes.TrimSpace(xb))
			if d != tc.E.data {
				t.Errorf("expected data '%s' got '%s'", tc.E.data, d)
			}
		})
	}
}

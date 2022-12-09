package handler

import (
	"bytes"
	"fmt"
	"github.com/dottics/dutil"
	"github.com/johannesscr/micro/microtest"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	type resp struct {
		status int
		data   string
	}
	type E struct {
		resp resp
	}
	tt := []struct {
		name string
		e    dutil.Err
		E    E
	}{
		{
			name: "400",
			e: dutil.Err{
				Status: 400,
				Errors: map[string][]string{
					"email": {"some email format error"},
				},
			},
			E: E{
				resp: resp{
					status: 400,
					data:   `{"message":"Bad Request","data":null,"errors":{"email":["some email format error"]}}`,
				},
			},
		},
		{
			name: "401",
			e: dutil.Err{
				Status: 401,
				Errors: map[string][]string{
					"auth": {"some auth error"},
				},
			},
			E: E{
				resp: resp{
					status: 401,
					data:   `{"message":"Unauthorized","data":null,"errors":{"auth":["some auth error"]}}`,
				},
			},
		},
		{
			name: "403",
			e: dutil.Err{
				Status: 403,
				Errors: map[string][]string{
					"auth": {"some other auth error"},
				},
			},
			E: E{
				resp: resp{
					status: 403,
					data:   `{"message":"Forbidden","data":null,"errors":{"auth":["some other auth error"]}}`,
				},
			},
		},
		{
			name: "500",
			e: dutil.Err{
				Status: 500,
				Errors: map[string][]string{
					"internal_server_error": {"some internal server error"},
				},
			},
			E: E{
				resp: resp{
					status: 500,
					data:   `{"message":"Internal Server Error","data":null,"errors":{"internal_server_error":["some internal server error"]}}`,
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()

			Error(rec, req, &tc.e)

			res, xb := microtest.ReadRecorder(rec)
			if res.StatusCode != tc.E.resp.status {
				t.Errorf("expected '%v' got '%v'", tc.E.resp.status, res.StatusCode)
			}
			s := string(bytes.TrimSpace(xb))
			if s != tc.E.resp.data {
				t.Errorf("expected '%v' got '%v'", tc.E.resp.data, s)
			}
		})
	}
}

func TestHome(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	Home(rec, req)

	res, xb := microtest.ReadRecorder(rec)
	if res.StatusCode != 200 {
		t.Errorf("expecrted %d got %d", 200, res.StatusCode)
	}
	d := `{"message":"Welcome to the Budget API Gateway","data":{"alive":true},"errors":null}`
	if string(bytes.TrimSpace(xb)) != d {
		t.Errorf("expected '%v' got  '%v'", d, string(bytes.TrimSpace(xb)))
	}
}

func TestLogin(t *testing.T) {
	type E struct {
		status int
		body   string
		token  string
	}
	tt := []struct {
		name     string
		body     string
		exchange *microtest.Exchange
		E        E
	}{
		{
			name: "incorrect data format",
			body: `{"email":"jamesbond.com","password":"007"}`,
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 400,
					Body:   `{"message":"BadRequest: Unable to process request","data":{},"errors":{"email":["incorrect format"]}}`,
				},
			},
			E: E{
				status: 400,
				body:   `{"message":"Bad Request","data":null,"errors":{"email":["incorrect format"]}}`,
				token:  "",
			},
		},
		{
			name: "successful login",
			body: `{"email":"james@bond.com","password":"007"}`,
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Header: map[string][]string{
						"Content-Type": {"application/json"},
						"X-User-Token": {"the-best-secret-token-ever"},
					},
					Body: `{"message":"login successful","data":{"user":{"uuid":"ba0cb87d-6644-4b92-9e1d-0503c5563a3e","first_name":"james","last_name":"bond","email":"james@bond.com","active":true},"permission":["aio2","91j8","s1ga","ai3h"]},"errors":{}}`,
				},
			},
			E: E{
				status: 200,
				body:   `{"message":"login successful","data":{"user":{"uuid":"ba0cb87d-6644-4b92-9e1d-0503c5563a3e","first_name":"james","last_name":"bond","email":"james@bond.com","contact_number":"","password_reset_token":"","active":true},"permission_codes":["aio2","91j8","s1ga","ai3h"]},"errors":null}`,
				token:  "the-best-secret-token-ever",
			},
		},
	}

	ms := microtest.NewMockServer("SECURITY_SERVICE_SCHEME", "SECURITY_SERVICE_HOST")
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ms.Append(tc.exchange)

			req := httptest.NewRequest("POST", "/login", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()

			Login(rec, req)
			res, xb := microtest.ReadRecorder(rec)
			if res.StatusCode != tc.E.status {
				t.Errorf("expected %d got %d", tc.E.status, res.StatusCode)
			}
			s := string(bytes.TrimSpace(xb))
			if s != tc.E.body {
				t.Errorf("expected '%v' got '%v'", tc.E.body, s)
			}
			token := res.Header.Get("X-Token")
			if tc.E.token != token {
				t.Errorf("expected '%v' got '%v'", tc.E.token, token)
			}

			// test the data is passed correctly
			//rc, _ := ioutil.ReadAll(req.Body)
			//log.Println("'", string(bytes.TrimSpace(rc)), "'")
			////log.Println(req.Body)
			//log.Println("EX BODY", tc.exchange.Request.Body)
			//rb, _ := ioutil.ReadAll(tc.exchange.Request.Body)
			//rs := string(bytes.TrimSpace(rb))
			//if rs != tc.body {
			//	t.Errorf("expected '%v' got '%v'", tc.body, rs)
			//}
		})
	}
}

func TestLogout(t *testing.T) {
	type E struct {
		status   int
		body     string
		reqToken string
	}

	tt := []struct {
		name     string
		token    string
		exchange *microtest.Exchange
		E        E
	}{
		{
			name:  "401 no token",
			token: "",
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 401,
					Body:   `{"message":"Unauthorized: unable to process request","data":{},"errors":{"authentication":["Auth token is missing","Please login"]}}`,
				},
			},
			E: E{
				status:   401,
				body:     `{"message":"Unauthorized","data":null,"errors":{"authentication":["Auth token is missing","Please login"]}}`,
				reqToken: "",
			},
		},
		{
			name:  "500 server error",
			token: "token for the 500",
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 500,
					Body:   `{"message":"InternalServerError: unable to process request","data":{},"errors":{"internal_server_error":["some unexpected server error"]}}`,
				},
			},
			E: E{
				status:   500,
				body:     `{"message":"Internal Server Error","data":null,"errors":{"internal_server_error":["some unexpected server error"]}}`,
				reqToken: "token for the 500",
			},
		},
		{
			name:  "logout successful",
			token: "token for the 200",
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   `{"message":"logout successful","data":null,"errors":null}`,
				},
			},
			E: E{
				status:   200,
				body:     `{"message":"logout successful","data":null,"errors":null}`,
				reqToken: "token for the 200",
			},
		},
	}

	ms := microtest.NewMockServer("SECURITY_SERVICE_SCHEME", "SECURITY_SERVICE_HOST")
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add the exchange between the api and the micro-service
			ms.Append(tc.exchange)
			req := httptest.NewRequest("DELETE", "/logout", nil)
			req.Header.Set("X-Token", tc.token)
			rec := httptest.NewRecorder()

			Logout(rec, req)
			res, xb := microtest.ReadRecorder(rec)
			// test the headers were passed correctly
			h := tc.exchange.Request.Header.Get("X-User-Token")
			if h != tc.E.reqToken {
				t.Errorf("expected '%v' got '%v'", tc.E.reqToken, h)
			}
			// test response
			if res.StatusCode != tc.E.status {
				t.Errorf("expected %d got %d", tc.E.status, res.StatusCode)
			}
			b := string(bytes.TrimSpace(xb))
			if b != tc.E.body {
				t.Errorf("expected '%v' got '%v'", tc.E.body, b)
			}
		})
	}
}

func TestResetPassword(t *testing.T) {
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
			name:    "Bad Request",
			payload: strings.NewReader(`{"email":"name@example.com","password_reset_token":"63a33ec6-1b11-4635-9def-391fc17bc6a0","password":"new password"}`),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 400,
					Body:   `{"message":"BadRequest","data":{},"errors":{"password":["must be more than 6 characters"]}}`,
				},
			},
			E: E{
				status: 400,
				data:   `{"message":"Bad Request","data":null,"errors":{"password":["must be more than 6 characters"]}}`,
			},
		},
		{
			name:    "Successful",
			payload: strings.NewReader(`{"email":"name@example.com","password_reset_token":"63a33ec6-1b11-4635-9def-391fc17bc6a0","password":"new password"}`),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   `{"message":"password reset successful","data":{},"errors":{}}`,
				},
			},
			E: E{
				status: 200,
				data:   `{"message":"password reset successful","data":null,"errors":null}`,
			},
		},
	}

	ms := microtest.NewMockServer("SECURITY_SERVICE_SCHEME", "SECURITY_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tests {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			ms.Append(tc.exchange)

			req := microtest.NewRequest("POST", "/reset-password", nil, nil, tc.payload)
			rec := httptest.NewRecorder()
			ResetPassword(rec, req)
			res, xb := microtest.ReadRecorder(rec)

			d := string(bytes.TrimSpace(xb))
			if res.StatusCode != tc.E.status {
				t.Errorf("expected status %d got %d", tc.E.status, res.StatusCode)
			}
			if d != tc.E.data {
				t.Errorf("expected data '%s' got '%s'", tc.E.data, d)
			}
		})
	}
}

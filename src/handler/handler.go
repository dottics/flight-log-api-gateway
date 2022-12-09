package handler

import (
	"bytes"
	"github.com/dottics/dutil"
	"github.com/dottics/securityserv"
	"io/ioutil"
	"net/http"
)

// Error handles all error responses
func Error(w http.ResponseWriter, r *http.Request, err dutil.Error) {
	e := dutil.Inst(err)
	res := dutil.Resp{
		Status:  e.Status,
		Message: http.StatusText(e.Status),
		Errors:  e.Errors,
	}
	res.Respond(w, r)
}

// Home is the health check for the server
func Home(w http.ResponseWriter, r *http.Request) {
	msg := struct {
		Alive bool `json:"alive"`
	}{
		Alive: true,
	}
	res := dutil.Resp{
		Status:  200,
		Message: "Welcome to the Budget API Gateway",
		Data:    msg,
	}
	res.Respond(w, r)
}

// Login handles the login of the budget api and  maps to the security service
func Login(w http.ResponseWriter, r *http.Request) {
	s := security.NewService("")
	xb, _ := ioutil.ReadAll(r.Body)
	_ = r.Body.Close()

	token, u, xs, e := s.Login(bytes.NewReader(xb))
	if e != nil {
		Error(w, r, e)
		return
	}

	data := struct {
		User            security.User            `json:"user"`
		PermissionCodes security.PermissionCodes `json:"permission_codes"`
	}{
		User:            u,
		PermissionCodes: xs,
	}

	resp := dutil.Resp{
		Status: 200,
		Header: map[string][]string{
			"X-Token": {token},
		},
		Message: "login successful",
		Data:    data,
	}
	resp.Respond(w, r)
}

// Logout handles the logout of the budget api and  maps to the security service
func Logout(w http.ResponseWriter, r *http.Request) {
	s := security.NewService(r.Header.Get("X-Token"))

	e := s.Logout()
	if e != nil {
		Error(w, r, e)
		return
	}

	resp := dutil.Resp{
		Status:  200,
		Message: "logout successful",
	}
	resp.Respond(w, r)
}

// ResetPassword handles the reset password of the budget API gateway and maps
// to the security service
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	s := security.NewService("")

	p := security.ResetPasswordPayload{}
	e := dutil.Decode(w, r, &p)
	if e != nil {
		Error(w, r, e)
		return
	}

	e = s.ResetPassword(p)
	if e != nil {
		Error(w, r, e)
		return
	}

	resp := dutil.Resp{
		Status:  200,
		Message: "password reset successful",
	}
	resp.Respond(w, r)
}

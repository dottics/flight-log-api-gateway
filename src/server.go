package src

import (
	"github.com/dottics/dutil"
	"github.com/dottics/flight-log-api-gateway/src/handler"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	Redis  bool
	Router *mux.Router
}

func NewServer() *Server {
	s := &Server{}
	s.Router = mux.NewRouter()

	// register routes
	s.routes()
	return s
}

// ServeHTTP is what makes the Server an HandlerFunc needed for the
// http.ListenAndServe function.
//func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	s.Router.ServeHTTP(w, r)
//}

/* ROUTING */

// routes sets all the possible endpoints available on the API.
func (s *Server) routes() {
	s.Router.HandleFunc("/", s.prop(handler.Home)).Methods("OPTIONS", "GET")
	// Auth
	s.Router.HandleFunc("/login", s.prop(handler.Login)).Methods("OPTIONS", "POST")
	s.Router.HandleFunc("/logout", s.prop(handler.Logout)).Methods("OPTIONS", "DELETE")
	s.Router.HandleFunc("/forgot-password", s.prop(handler.ForgotPassword)).Methods("OPTIONS", "POST")
	s.Router.HandleFunc("/reset-password", s.prop(handler.ResetPassword)).Methods("OPTIONS", "POST")
	s.Router.HandleFunc("/contact-us", s.prop(handler.ContactUs)).Methods("OPTIONS", "POST")
	// aircraft types
	//s.Router.HandleFunc("/aircraft-type", s.prop(handler.AircraftTypes)).Methods("OPTIONS", "GET")
	//// CRUD flight log
	//s.Router.HandleFunc("/flight-log", s.prop(handler.FlightLog)).Methods("OPTIONS", "GET")
	//s.Router.HandleFunc("/flight-log/-", s.prop(handler.FlightLogs)).Methods("OPTIONS", "GET")
	//s.Router.HandleFunc("/flight-log", s.prop(handler.CreateFlightLog)).Methods("OPTIONS", "POST")
	//s.Router.HandleFunc("/flight-log/-", s.prop(handler.UpdateFlightLog)).Methods("OPTIONS", "PUT")
	//s.Router.HandleFunc("/flight-log/-", s.prop(handler.DeleteFlightLog)).Methods("OPTIONS", "DELETE")
}

// prop propagates the http.ResponseWriter and http.Request to the handler
// function f. Primary use is as a form of middleware to allow for CORS.
func (s *Server) prop(f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			resp := dutil.Resp{
				Status:  200,
				Message: "Pre-Flight Allowed",
			}
			resp.Respond(w, r)
			return
		}
		f(w, r)
	}
}

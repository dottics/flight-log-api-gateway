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
	// Budget
	//s.Router.HandleFunc("/budget", s.prop(handler.Budgets)).Methods("OPTIONS", "GET")
	//s.Router.HandleFunc("/budget/-", s.prop(handler.Budget)).Methods("OPTIONS", "GET")
	//s.Router.HandleFunc("/budget/-/group", s.prop(handler.Groups)).Methods("OPTIONS", "GET")
	//s.Router.HandleFunc("/budget/group/-/item", s.prop(handler.GetItems)).Methods("OPTIONS", "GET")
	//s.Router.HandleFunc("/budget/group/item/-/event", s.prop(handler.GetEvents)).Methods("OPTIONS", "GET")
	//s.Router.HandleFunc("/calc/item/month-total", s.prop(handler.ItemMonthlyTotal)).Methods("OPTIONS", "POST")
	//// CRUD event
	//s.Router.HandleFunc("/event", s.prop(handler.CreateEvent)).Methods("OPTIONS", "POST")
	//s.Router.HandleFunc("/event/-", s.prop(handler.UpdateEvent)).Methods("OPTIONS", "PUT")
	//s.Router.HandleFunc("/event/-", s.prop(handler.DeleteEvent)).Methods("OPTIONS", "DELETE")
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

package handler

import (
	"github.com/dottics/dutil"
	"github.com/dottics/flight-log-api-gateway/src/includes"
	"github.com/dottics/flightserv"
	"net/http"
)

// AircraftTypes is the handler that manages fetching the aircraft types that
// are available to the app.
func AircraftTypes(w http.ResponseWriter, r *http.Request) {
	// get token
	xt, e := includes.GetAircraftTypes(r.Header.Get("X-Token"))
	if e != nil {
		Error(w, r, e)
		return
	}

	data := struct {
		flightserv.AircraftTypes `json:"aircraftTypes"`
	}{
		AircraftTypes: xt,
	}
	res := dutil.Resp{
		Status:  200,
		Message: "aircraft types found",
		Data:    data,
	}
	res.Respond(w, r)
}

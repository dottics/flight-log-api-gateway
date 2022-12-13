package handler

import (
	"github.com/dottics/dutil"
	"github.com/dottics/flight-log-api-gateway/src/includes"
	"github.com/dottics/flightserv"
	"github.com/google/uuid"
	"net/http"
)

// FlightLog is the handler to get a flight log from the flight log
// service.
func FlightLog(w http.ResponseWriter, r *http.Request) {
	// get token
	token := r.Header.Get("X-Token")
	userUUID, err := uuid.Parse(r.URL.Query().Get("userUuid"))
	if err != nil {
		e := &dutil.Err{
			Status: 400,
			Errors: map[string][]string{
				"userUuid": {err.Error()},
			},
		}
		Error(w, r, e)
		return
	}
	UUID, err := uuid.Parse(r.URL.Query().Get("uuid"))
	if err != nil {
		e := &dutil.Err{
			Status: 400,
			Errors: map[string][]string{
				"userUuid": {err.Error()},
			},
		}
		Error(w, r, e)
	}

	log, e := includes.GetFlightLog(token, userUUID, UUID)
	if e != nil {
		Error(w, r, e)
		return
	}

	type Data struct {
		FlightLog flightserv.FlightLog `json:"flightLog"`
	}
	res := dutil.Resp{
		Status:  200,
		Message: "flight log found",
		Data: Data{
			FlightLog: log,
		},
	}
	res.Respond(w, r)
}

// FlightLogs is the handler to get a flight logs from the flight log
// service.
func FlightLogs(w http.ResponseWriter, r *http.Request) {
	// get token
	token := r.Header.Get("X-Token")
	userUUID, err := uuid.Parse(r.URL.Query().Get("userUuid"))
	if err != nil {
		e := &dutil.Err{
			Status: 400,
			Errors: map[string][]string{
				"userUuid": {err.Error()},
			},
		}
		Error(w, r, e)
		return
	}

	xLog, e := includes.GetFlightLogs(token, userUUID)
	if e != nil {
		Error(w, r, e)
		return
	}

	type Data struct {
		FlightLogs flightserv.FlightLogs `json:"flightLogs"`
	}
	res := dutil.Resp{
		Status:  200,
		Message: "flight logs found",
		Data: Data{
			FlightLogs: xLog,
		},
	}
	res.Respond(w, r)
}

// CreateFlightLog is the handler to create a new flight log for the user in the
// flight log service.
func CreateFlightLog(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Token")
	flightLogData := flightserv.FlightLog{}
	e := dutil.Decode(w, r, &flightLogData)
	if e != nil {
		Error(w, r, e)
		return
	}

	log, e := includes.CreateFlightLog(token, flightLogData)
	if e != nil {
		Error(w, r, e)
		return
	}

	type Data struct {
		FlightLog flightserv.FlightLog `json:"flightLog"`
	}
	res := dutil.Resp{
		Status:  201,
		Message: "flight log created",
		Data: Data{
			FlightLog: log,
		},
	}
	res.Respond(w, r)
}

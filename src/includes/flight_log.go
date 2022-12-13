package includes

import (
	"github.com/dottics/dutil"
	"github.com/dottics/flightserv"
	"github.com/google/uuid"
)

// GetFlightLogs fetches all the flight logs of a user from the flight log
// service.
func GetFlightLogs(token string, UserUUID uuid.UUID) (flightserv.FlightLogs, dutil.Error) {
	ms := flightserv.NewService(token)
	xlogs, e := ms.GetFlightLogs(UserUUID)
	return xlogs, e
}

// GetFlightLog fetches a specific flight log of a user from the flight log
// service.
func GetFlightLog(token string, UserUUID, UUID uuid.UUID) (flightserv.FlightLog, dutil.Error) {
	ms := flightserv.NewService(token)
	log, e := ms.GetFlightLog(UserUUID, UUID)
	return log, e
}

// CreateFlightLog exchanges with the flight log service to create a new flight
// log.
func CreateFlightLog(token string, log flightserv.FlightLog) (flightserv.FlightLog, dutil.Error) {
	ms := flightserv.NewService(token)
	log, e := ms.CreateFlightLog(log)
	return log, e
}

// UpdateFlightLog exchanges with the flight log service to update a flight log.
func UpdateFlightLog(token string, log flightserv.FlightLog) (flightserv.FlightLog, dutil.Error) {
	ms := flightserv.NewService(token)
	log, e := ms.UpdateFlightLog(log)
	return log, e
}

// DeleteFlightLog exchanges with the flight log service to delete a user's
// flight log. And only returns an error if an error occurred.
func DeleteFlightLog(token string, UserUUID, UUID uuid.UUID) dutil.Error {
	ms := flightserv.NewService(token)
	e := ms.DeleteFlightLog(UserUUID, UUID)
	return e
}

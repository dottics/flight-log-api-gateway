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
	xl, e := ms.GetFlightLogs(UserUUID)
	return xl, e
}

// GetFlightLog fetches a specific flight log of a user from the flight log
// service.
func GetFlightLog(UserUUID, UUID uuid.UUID) (flightserv.FlightLog, dutil.Error) {
	return flightserv.FlightLog{}, nil
}

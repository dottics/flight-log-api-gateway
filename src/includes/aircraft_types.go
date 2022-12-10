package includes

import (
	"github.com/dottics/dutil"
	"github.com/dottics/flightserv"
)

// GetAircraftTypes fetches all the aircraft types from the microservice.
func GetAircraftTypes(token string) (flightserv.AircraftTypes, dutil.Error) {
	ms := flightserv.NewService(token)
	xf, e := ms.GetAircraftTypes()
	return xf, e
}

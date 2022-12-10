package includes

import (
	"fmt"
	"github.com/dottics/dutil"
	"github.com/dottics/flightserv"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/microtest"
	"testing"
)

func TestGetFlightLogs(t *testing.T) {
	tt := []struct {
		name       string
		userUUID   uuid.UUID
		ex         *microtest.Exchange
		flightLogs flightserv.FlightLogs
		e          dutil.Error
	}{
		{
			name:     "403 permission required",
			userUUID: uuid.MustParse("d12cf499-7fb5-4495-9e7c-e9ceb1e0ad77"),
			ex: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body: `{
						"message":"Forbidden",
						"errors":{
							"permission": ["Please ensure you have permission"]
						}
					}`,
				},
			},
			flightLogs: flightserv.FlightLogs{},
			e: &dutil.Err{
				Status: 403,
				Errors: map[string][]string{
					"permission": {"Please ensure you have permission"},
				},
			},
		},
		{
			name:     "200 fetches the user's flight logs",
			userUUID: uuid.MustParse("85fc976b-866b-4c4a-b28c-0996c9671a90"),
			ex: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message": "flight logs found",
						"data": {
							"flightLogs": [
								{
									"uuid": "2d916236-62a6-40b5-9b22-83202bbe6a48",
									"userUuid": "85fc976b-866b-4c4a-b28c-0996c9671a90"
								},
								{
									"uuid": "62c491b1-4aad-4889-b4f0-281063d7f5c8",
									"userUuid": "85fc976b-866b-4c4a-b28c-0996c9671a90"
								}
							]
						}
					}`,
				},
			},
			flightLogs: flightserv.FlightLogs{
				flightserv.FlightLog{
					UUID:     uuid.MustParse("2d916236-62a6-40b5-9b22-83202bbe6a48"),
					UserUUID: uuid.MustParse("85fc976b-866b-4c4a-b28c-0996c9671a90"),
				},
				flightserv.FlightLog{
					UUID:     uuid.MustParse("62c491b1-4aad-4889-b4f0-281063d7f5c8"),
					UserUUID: uuid.MustParse("85fc976b-866b-4c4a-b28c-0996c9671a90"),
				},
			},
			e: nil,
		},
	}

	ms := microtest.NewMockServer("FLIGHT_SERVICE_SCHEME", "FLIGHT_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tt {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			ms.Append(tc.ex)

			xlogs, e := GetFlightLogs("test-token", tc.userUUID)
			// test request formatted correctly
			query := fmt.Sprintf("userUUID=%s", tc.userUUID.String())
			if tc.ex.Request.URL.RawQuery != query {
				t.Errorf("expected query string %s got %s", query, tc.ex.Request.URL.RawQuery)
			}

			if !dutil.ErrorEqual(tc.e, e) {
				t.Errorf("expected error %v got %v", tc.e, e)
			}
			if !tc.flightLogs.EqualTo(xlogs) {
				t.Errorf("expected flight logs %v got %v", tc.flightLogs, xlogs)
			}
		})
	}
}

func TestGetFlightLog(t *testing.T) {
	tt := []struct {
		name      string
		userUUID  uuid.UUID
		UUID      uuid.UUID
		ex        *microtest.Exchange
		flightLog flightserv.FlightLog
		e         dutil.Error
	}{
		{
			name:     "403 permission required",
			userUUID: uuid.MustParse("d12cf499-7fb5-4495-9e7c-e9ceb1e0ad77"),
			UUID:     uuid.MustParse("481f8179-a1ff-4099-a69a-814679fa5a83"),
			ex: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body: `{
						"message":"Forbidden",
						"errors":{
							"permission": ["Please ensure you have permission"]
						}
					}`,
				},
			},
			flightLog: flightserv.FlightLog{},
			e: &dutil.Err{
				Status: 403,
				Errors: map[string][]string{
					"permission": {"Please ensure you have permission"},
				},
			},
		},
		{
			name:     "200 fetches the user's flight logs",
			userUUID: uuid.MustParse("85fc976b-866b-4c4a-b28c-0996c9671a90"),
			UUID:     uuid.MustParse("2d916236-62a6-40b5-9b22-83202bbe6a48"),
			ex: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message": "flight log found",
						"data": {
							"flightLog": {
								"uuid": "2d916236-62a6-40b5-9b22-83202bbe6a48",
								"userUuid": "85fc976b-866b-4c4a-b28c-0996c9671a90"
							}
						}
					}`,
				},
			},
			flightLog: flightserv.FlightLog{
				UUID:     uuid.MustParse("2d916236-62a6-40b5-9b22-83202bbe6a48"),
				UserUUID: uuid.MustParse("85fc976b-866b-4c4a-b28c-0996c9671a90"),
			},
			e: nil,
		},
	}

	ms := microtest.NewMockServer("FLIGHT_SERVICE_SCHEME", "FLIGHT_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tt {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			ms.Append(tc.ex)

			log, e := GetFlightLog("test-token", tc.userUUID, tc.UUID)
			// test request formatted correctly
			query := fmt.Sprintf("UUID=%s&userUUID=%s", tc.UUID.String(), tc.userUUID.String())
			if tc.ex.Request.URL.RawQuery != query {
				t.Errorf("expected query string %s got %s", query, tc.ex.Request.URL.RawQuery)
			}

			if !dutil.ErrorEqual(tc.e, e) {
				t.Errorf("expected error %v got %v", tc.e, e)
			}
			if tc.flightLog != log {
				t.Errorf("expected flight logs %v got %v", tc.flightLog, log)
			}
		})
	}
}

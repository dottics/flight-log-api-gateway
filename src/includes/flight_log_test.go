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

func TestCreateFlightLog(t *testing.T) {
	type E struct {
		log flightserv.FlightLog
		e   dutil.Error
	}
	tt := []struct {
		name string
		log  flightserv.FlightLog
		ex   *microtest.Exchange
		E
	}{
		{
			name: "403 Forbidden",
			log: flightserv.FlightLog{
				UUID:     uuid.MustParse("76b1841b-d6c5-4696-8fbe-b062056ccd76"),
				UserUUID: uuid.MustParse("b8530d7d-6678-432f-a9be-6b5847c3ef2f"),
			},
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
			E: E{
				log: flightserv.FlightLog{},
				e: &dutil.Err{
					Status: 403,
					Errors: map[string][]string{
						"permission": {"Please ensure you have permission"},
					},
				},
			},
		},
		{
			name: "201 Create flight log successful",
			log: flightserv.FlightLog{
				UserUUID: uuid.MustParse("b8530d7d-6678-432f-a9be-6b5847c3ef2f"),
			},
			ex: &microtest.Exchange{
				Response: microtest.Response{
					Status: 201,
					Body: `{
						"message": "flight log created",
						"data": {
							"flightLog": {
								"uuid": "76b1841b-d6c5-4696-8fbe-b062056ccd76",
								"userUuid": "b8530d7d-6678-432f-a9be-6b5847c3ef2f"
							}
						}
					}`,
				},
			},
			E: E{
				log: flightserv.FlightLog{
					UUID:     uuid.MustParse("76b1841b-d6c5-4696-8fbe-b062056ccd76"),
					UserUUID: uuid.MustParse("b8530d7d-6678-432f-a9be-6b5847c3ef2f"),
				},
				e: nil,
			},
		},
	}

	ms := microtest.NewMockServer("FLIGHT_SERVICE_SCHEME", "FLIGHT_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tt {
		name := fmt.Sprintf("%d %s", i, tc.name)

		t.Run(name, func(t *testing.T) {
			ms.Append(tc.ex)

			log, e := CreateFlightLog("test-token", tc.log)
			if tc.ex.Request.Method != "POST" {
				t.Errorf("expected create to be a POST request got %s", tc.ex.Request.Method)
			}
			if !dutil.ErrorEqual(tc.E.e, e) {
				t.Errorf("expected error %v got %v", tc.E.e, e)
			}
			if log != tc.E.log {
				t.Errorf("expected flight log %v got %v", tc.E.log, log)
			}
		})
	}
}

func TestUpdateFlightLog(t *testing.T) {
	type E struct {
		log flightserv.FlightLog
		e   dutil.Error
	}
	tt := []struct {
		name string
		log  flightserv.FlightLog
		ex   *microtest.Exchange
		E
	}{
		{
			name: "403 Forbidden",
			log: flightserv.FlightLog{
				UUID:     uuid.MustParse("76b1841b-d6c5-4696-8fbe-b062056ccd76"),
				UserUUID: uuid.MustParse("b8530d7d-6678-432f-a9be-6b5847c3ef2f"),
			},
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
			E: E{
				log: flightserv.FlightLog{},
				e: &dutil.Err{
					Status: 403,
					Errors: map[string][]string{
						"permission": {"Please ensure you have permission"},
					},
				},
			},
		},
		{
			name: "200 Update flight log successful",
			log: flightserv.FlightLog{
				UserUUID: uuid.MustParse("b8530d7d-6678-432f-a9be-6b5847c3ef2f"),
			},
			ex: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message": "flight log updated",
						"data": {
							"flightLog": {
								"uuid": "76b1841b-d6c5-4696-8fbe-b062056ccd76",
								"userUuid": "b8530d7d-6678-432f-a9be-6b5847c3ef2f"
							}
						}
					}`,
				},
			},
			E: E{
				log: flightserv.FlightLog{
					UUID:     uuid.MustParse("76b1841b-d6c5-4696-8fbe-b062056ccd76"),
					UserUUID: uuid.MustParse("b8530d7d-6678-432f-a9be-6b5847c3ef2f"),
				},
				e: nil,
			},
		},
	}

	ms := microtest.NewMockServer("FLIGHT_SERVICE_SCHEME", "FLIGHT_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tt {
		name := fmt.Sprintf("%d %s", i, tc.name)

		t.Run(name, func(t *testing.T) {
			ms.Append(tc.ex)

			log, e := UpdateFlightLog("test-token", tc.log)
			if tc.ex.Request.Method != "PUT" {
				t.Errorf("expected update to be a PUT request got %s", tc.ex.Request.Method)
			}
			if !dutil.ErrorEqual(tc.E.e, e) {
				t.Errorf("expected error %v got %v", tc.E.e, e)
			}
			if log != tc.E.log {
				t.Errorf("expected flight log %v got %v", tc.E.log, log)
			}
		})
	}
}

func TestDeleteFlightLog(t *testing.T) {
	type E struct {
		e dutil.Error
	}
	tt := []struct {
		name     string
		userUUID uuid.UUID
		UUID     uuid.UUID
		ex       *microtest.Exchange
		E
	}{
		{
			name:     "403 Forbidden",
			userUUID: uuid.MustParse("b8530d7d-6678-432f-a9be-6b5847c3ef2f"),
			UUID:     uuid.MustParse("a6ca1dba-9f15-44de-8b1d-1f6ee7bf92d8"),
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
			E: E{
				e: &dutil.Err{
					Status: 403,
					Errors: map[string][]string{
						"permission": {"Please ensure you have permission"},
					},
				},
			},
		},
		{
			name:     "200 delete flight log successful",
			userUUID: uuid.MustParse("78386229-9562-46cb-9441-afa9ce42fd70"),
			UUID:     uuid.MustParse("9403f65d-8fe1-4c87-bd49-49c051f6c4bb"),
			ex: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message": "flight log deleted"
					}`,
				},
			},
			E: E{
				e: nil,
			},
		},
	}

	ms := microtest.NewMockServer("FLIGHT_SERVICE_SCHEME", "FLIGHT_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tt {
		name := fmt.Sprintf("%d %s", i, tc.name)

		t.Run(name, func(t *testing.T) {
			ms.Append(tc.ex)

			e := DeleteFlightLog("test-token", tc.userUUID, tc.UUID)
			query := fmt.Sprintf("UUID=%s&userUUID=%s", tc.UUID.String(), tc.userUUID.String())
			if tc.ex.Request.URL.RawQuery != query {
				t.Errorf("expected query parameters %s got %s", query, tc.ex.Request.URL.RawQuery)
			}
			if tc.ex.Request.Method != "DELETE" {
				t.Errorf("expected update to be a DELETE request got %s", tc.ex.Request.Method)
			}
			if !dutil.ErrorEqual(tc.E.e, e) {
				t.Errorf("expected error %v got %v", tc.E.e, e)
			}
		})
	}
}

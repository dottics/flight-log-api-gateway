package handler

import (
	"fmt"
	"github.com/johannesscr/micro/microtest"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFlightLog(t *testing.T) {
	type E struct {
		status int
		body   string
	}
	tt := []struct {
		name string
		qs   map[string]string
		ex   *microtest.Exchange
		E
	}{
		{
			name: "403 Forbidden",
			qs: map[string]string{
				"userUuid": "fdb9cb92-ad97-43d6-b65a-17af732c4dec",
				"uuid":     "3a035c92-b601-4979-9429-7b78adcf7fe0",
			},
			ex: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body: `{
						"message":"Forbidden",
						"errors":{
							"permission":["Please ensure you have permission"]
						}
					}`,
				},
			},
			E: E{
				status: 403,
				body:   `{"message":"Forbidden","data":null,"errors":{"permission":["Please ensure you have permission"]}}`,
			},
		},
		{
			name: "200 get a specific flight log",
			qs: map[string]string{
				"userUuid": "09f43744-6737-446a-8d4b-eb4a136e08e6",
				"uuid":     "b594d3db-16c2-4a5f-b6ec-e08663a4c7b3",
			},
			ex: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"flight log found",
						"data": {
							"flightLog": {
								"uuid":"b594d3db-16c2-4a5f-b6ec-e08663a4c7b3",
								"userUuid":"09f43744-6737-446a-8d4b-eb4a136e08e6"
							}
						}
					}`,
				},
			},
			E: E{
				status: 200,
				body:   `{"message":"flight log found","data":{"flightLog":{"uuid":"b594d3db-16c2-4a5f-b6ec-e08663a4c7b3","userUuid":"09f43744-6737-446a-8d4b-eb4a136e08e6","aircraftTypeUuid":"00000000-0000-0000-0000-000000000000","aircraftType":"","date":"0001-01-01T00:00:00Z","registration":"","pilotInCommand":"","details":"","instrumentNavAids":"","instrumentPlace":"","instrumentActual":"","instrumentFSTD":0,"instructorSE":0,"instructorME":0,"instructorFSTD":0,"fstd":0,"engineType":"","dayType":"","dual":0,"pic":0,"picus":0,"copilot":0,"dayLandings":0,"nightLandings":0,"remarks":""}},"errors":null}`,
			},
		},
	}

	ms := microtest.NewMockServer("FLIGHT_SERVICE_SCHEME", "FLIGHT_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tt {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			ms.Append(tc.ex)
			req := microtest.NewRequest("GET", "/log/-", tc.qs, nil, nil)
			rec := httptest.NewRecorder()

			FlightLog(rec, req)
			query := fmt.Sprintf("userUuid=%s&uuid=%s", tc.qs["userUuid"], tc.qs["uuid"])
			if req.URL.RawQuery != query {
				t.Errorf("expected query parameters '%s' got '%s'", query, req.URL.RawQuery)
			}
			res, xb := microtest.ReadRecorder(rec)
			if res.StatusCode != tc.E.status {
				t.Errorf("expected status code %d got %d", tc.E.status, res.StatusCode)
			}
			body := strings.TrimSpace(string(xb))
			if body != tc.E.body {
				t.Errorf("expected body '%s' got '%s'", tc.E.body, body)
			}
		})
	}
}

func TestFlightLogs(t *testing.T) {
	type E struct {
		status int
		body   string
	}
	tt := []struct {
		name string
		qs   map[string]string
		ex   *microtest.Exchange
		E
	}{
		{
			name: "403 Forbidden",
			qs: map[string]string{
				"userUuid": "fdb9cb92-ad97-43d6-b65a-17af732c4dec",
			},
			ex: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body: `{
						"message":"Forbidden",
						"errors":{
							"permission":["Please ensure you have permission"]
						}
					}`,
				},
			},
			E: E{
				status: 403,
				body:   `{"message":"Forbidden","data":null,"errors":{"permission":["Please ensure you have permission"]}}`,
			},
		},
		{
			name: "200 get flight logs",
			qs: map[string]string{
				"userUuid": "09f43744-6737-446a-8d4b-eb4a136e08e6",
			},
			ex: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"flight log found",
						"data": {
							"flightLogs": [
								{
									"uuid":"b594d3db-16c2-4a5f-b6ec-e08663a4c7b3",
									"userUuid":"09f43744-6737-446a-8d4b-eb4a136e08e6"
								},
								{
									"uuid":"c377b3d4-8dd4-43c9-95bb-5957dddef582",
									"userUuid":"09f43744-6737-446a-8d4b-eb4a136e08e6"
								}
							]
						}
					}`,
				},
			},
			E: E{
				status: 200,
				body:   `{"message":"flight logs found","data":{"flightLogs":[{"uuid":"b594d3db-16c2-4a5f-b6ec-e08663a4c7b3","userUuid":"09f43744-6737-446a-8d4b-eb4a136e08e6","aircraftTypeUuid":"00000000-0000-0000-0000-000000000000","aircraftType":"","date":"0001-01-01T00:00:00Z","registration":"","pilotInCommand":"","details":"","instrumentNavAids":"","instrumentPlace":"","instrumentActual":"","instrumentFSTD":0,"instructorSE":0,"instructorME":0,"instructorFSTD":0,"fstd":0,"engineType":"","dayType":"","dual":0,"pic":0,"picus":0,"copilot":0,"dayLandings":0,"nightLandings":0,"remarks":""},{"uuid":"c377b3d4-8dd4-43c9-95bb-5957dddef582","userUuid":"09f43744-6737-446a-8d4b-eb4a136e08e6","aircraftTypeUuid":"00000000-0000-0000-0000-000000000000","aircraftType":"","date":"0001-01-01T00:00:00Z","registration":"","pilotInCommand":"","details":"","instrumentNavAids":"","instrumentPlace":"","instrumentActual":"","instrumentFSTD":0,"instructorSE":0,"instructorME":0,"instructorFSTD":0,"fstd":0,"engineType":"","dayType":"","dual":0,"pic":0,"picus":0,"copilot":0,"dayLandings":0,"nightLandings":0,"remarks":""}]},"errors":null}`,
			},
		},
	}

	ms := microtest.NewMockServer("FLIGHT_SERVICE_SCHEME", "FLIGHT_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tt {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			ms.Append(tc.ex)
			req := microtest.NewRequest("GET", "/log", tc.qs, nil, nil)
			rec := httptest.NewRecorder()

			FlightLogs(rec, req)
			query := fmt.Sprintf("userUuid=%s", tc.qs["userUuid"])
			if req.URL.RawQuery != query {
				t.Errorf("expected query parameters '%s' got '%s'", query, req.URL.RawQuery)
			}
			res, xb := microtest.ReadRecorder(rec)
			if res.StatusCode != tc.E.status {
				t.Errorf("expected status code %d got %d", tc.E.status, res.StatusCode)
			}
			body := strings.TrimSpace(string(xb))
			if body != tc.E.body {
				t.Errorf("expected body '%s' got '%s'", tc.E.body, body)
			}
		})
	}
}

func TestCreateFlightLog(t *testing.T) {
	type E struct {
		status int
		body   string
	}
	tt := []struct {
		name string
		body io.Reader
		ex   *microtest.Exchange
		E
	}{
		{
			name: "403 Forbidden",
			body: strings.NewReader(`{}`),
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
				status: 403,
				body:   `{"message":"Forbidden","data":null,"errors":{"permission":["Please ensure you have permission"]}}`,
			},
		},
		{
			name: "201 Create flight log successful",
			body: strings.NewReader(`{}`),
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
				status: 201,
				body:   `{"message":"flight log created","data":{"flightLog":{"uuid":"76b1841b-d6c5-4696-8fbe-b062056ccd76","userUuid":"b8530d7d-6678-432f-a9be-6b5847c3ef2f","aircraftTypeUuid":"00000000-0000-0000-0000-000000000000","aircraftType":"","date":"0001-01-01T00:00:00Z","registration":"","pilotInCommand":"","details":"","instrumentNavAids":"","instrumentPlace":"","instrumentActual":"","instrumentFSTD":0,"instructorSE":0,"instructorME":0,"instructorFSTD":0,"fstd":0,"engineType":"","dayType":"","dual":0,"pic":0,"picus":0,"copilot":0,"dayLandings":0,"nightLandings":0,"remarks":""}},"errors":null}`,
			},
		},
	}

	ms := microtest.NewMockServer("FLIGHT_SERVICE_SCHEME", "FLIGHT_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tt {
		name := fmt.Sprintf("%d %s", i, tc.name)

		t.Run(name, func(t *testing.T) {
			ms.Append(tc.ex)
			req := microtest.NewRequest("POST", "/flight-log", nil, nil, tc.body)
			rec := httptest.NewRecorder()
			CreateFlightLog(rec, req)

			res, xb := microtest.ReadRecorder(rec)
			if res.StatusCode != tc.E.status {
				t.Errorf("expected status code %d got %d", tc.E.status, res.StatusCode)
			}

			b := strings.TrimSpace(string(xb))
			if b != tc.E.body {
				t.Errorf("expected body '%s' got '%s'", tc.E.body, b)
			}
		})
	}
}

func TestUpdateFlightLog(t *testing.T) {
	type E struct {
		status int
		body   string
	}
	tt := []struct {
		name string
		body io.Reader
		ex   *microtest.Exchange
		E
	}{
		{
			name: "403 Forbidden",
			body: strings.NewReader(`{}`),
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
				status: 403,
				body:   `{"message":"Forbidden","data":null,"errors":{"permission":["Please ensure you have permission"]}}`,
			},
		},
		{
			name: "201 Create flight log successful",
			body: strings.NewReader(`{}`),
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
				status: 200,
				body:   `{"message":"flight log updated","data":{"flightLog":{"uuid":"76b1841b-d6c5-4696-8fbe-b062056ccd76","userUuid":"b8530d7d-6678-432f-a9be-6b5847c3ef2f","aircraftTypeUuid":"00000000-0000-0000-0000-000000000000","aircraftType":"","date":"0001-01-01T00:00:00Z","registration":"","pilotInCommand":"","details":"","instrumentNavAids":"","instrumentPlace":"","instrumentActual":"","instrumentFSTD":0,"instructorSE":0,"instructorME":0,"instructorFSTD":0,"fstd":0,"engineType":"","dayType":"","dual":0,"pic":0,"picus":0,"copilot":0,"dayLandings":0,"nightLandings":0,"remarks":""}},"errors":null}`,
			},
		},
	}

	ms := microtest.NewMockServer("FLIGHT_SERVICE_SCHEME", "FLIGHT_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tt {
		name := fmt.Sprintf("%d %s", i, tc.name)

		t.Run(name, func(t *testing.T) {
			ms.Append(tc.ex)
			req := microtest.NewRequest("PUT", "/flight-log/-", nil, nil, tc.body)
			rec := httptest.NewRecorder()
			UpdateFlightLog(rec, req)

			res, xb := microtest.ReadRecorder(rec)
			if res.StatusCode != tc.E.status {
				t.Errorf("expected status code %d got %d", tc.E.status, res.StatusCode)
			}

			b := strings.TrimSpace(string(xb))
			if b != tc.E.body {
				t.Errorf("expected body '%s' got '%s'", tc.E.body, b)
			}
		})
	}
}

func TestDeleteFlightLog(t *testing.T) {
	type E struct {
		status int
		body   string
	}
	tt := []struct {
		name string
		qs   map[string]string
		ex   *microtest.Exchange
		E
	}{
		{
			name: "403 Forbidden",
			qs: map[string]string{
				"uuid":     "b8530d7d-6678-432f-a9be-6b5847c3ef2f",
				"userUuid": "a6ca1dba-9f15-44de-8b1d-1f6ee7bf92d8",
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
				status: 403,
				body:   `{"message":"Forbidden","data":null,"errors":{"permission":["Please ensure you have permission"]}}`,
			},
		},
		{
			name: "200 Delete flight log successful",
			qs: map[string]string{
				"uuid":     "78386229-9562-46cb-9441-afa9ce42fd70",
				"userUuid": "9403f65d-8fe1-4c87-bd49-49c051f6c4bb",
			},
			ex: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message": "flight log deleted"
					}`,
				},
			},
			E: E{
				status: 200,
				body:   `{"message":"flight log deleted","data":null,"errors":null}`,
			},
		},
	}

	ms := microtest.NewMockServer("FLIGHT_SERVICE_SCHEME", "FLIGHT_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tt {
		name := fmt.Sprintf("%d %s", i, tc.name)

		t.Run(name, func(t *testing.T) {
			ms.Append(tc.ex)
			req := microtest.NewRequest("DELETE", "/flight-log/-", tc.qs, nil, nil)
			rec := httptest.NewRecorder()
			DeleteFlightLog(rec, req)

			query := fmt.Sprintf("UUID=%s&userUUID=%s", tc.qs["uuid"], tc.qs["userUuid"])
			if tc.ex.Request.URL.RawQuery != query {
				t.Errorf("expected query parameters %s got %s", query, tc.ex.Request.URL.RawQuery)
			}
			res, xb := microtest.ReadRecorder(rec)
			if res.StatusCode != tc.E.status {
				t.Errorf("expected status code %d got %d", tc.E.status, res.StatusCode)
			}

			b := strings.TrimSpace(string(xb))
			if b != tc.E.body {
				t.Errorf("expected body '%s' got '%s'", tc.E.body, b)
			}
		})
	}
}

package handler

import (
	"bytes"
	"fmt"
	"github.com/johannesscr/micro/microtest"
	"net/http/httptest"
	"testing"
)

func TestGetAircraftTypes(t *testing.T) {
	type E struct {
		status int
		body   string
	}
	tt := []struct {
		name string
		ex   *microtest.Exchange
		E
	}{
		{
			name: "403 permission required",
			ex: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   `{"message":"Forbidden","errors":{"permission":["Please ensure you have permission"]}}`,
				},
			},
			E: E{
				status: 403,
				body:   `{"message":"Forbidden","data":null,"errors":{"permission":["Please ensure you have permission"]}}`,
			},
		},
		{
			name: "200 get aircraft types",
			ex: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"aircraft types found",
						"data":{
							"aircraftTypes":[
								{
									"uuid":"08e670a5-3562-4538-b48b-dd04c5387617",
									"name":"A380",
									"description":""
								},
								{
									"uuid":"fa26565d-3bf6-49b0-a1d6-138a431bf7cd",
									"name":"F210",
									"description":""
								}
							]
						}
					}`,
				},
			},
			E: E{
				status: 200,
				body:   `{"message":"aircraft types found","data":{"aircraftTypes":[{"uuid":"08e670a5-3562-4538-b48b-dd04c5387617","name":"A380","description":""},{"uuid":"fa26565d-3bf6-49b0-a1d6-138a431bf7cd","name":"F210","description":""}]},"errors":null}`,
			},
		},
	}

	ms := microtest.NewMockServer("FLIGHT_SERVICE_SCHEME", "FLIGHT_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tt {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			ms.Append(tc.ex)
			req := httptest.NewRequest("GET", "/aircraft-type", nil)
			rec := httptest.NewRecorder()

			GetAircraftTypes(rec, req)
			res, xb := microtest.ReadRecorder(rec)
			if res.StatusCode != tc.E.status {
				t.Errorf("expected %d got %d", tc.E.status, res.StatusCode)
			}
			s := string(bytes.TrimSpace(xb))
			if s != tc.E.body {
				t.Errorf("expected '%v' got '%v'", tc.E.body, s)
			}
		})
	}
}

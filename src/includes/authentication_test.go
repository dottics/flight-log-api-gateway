package includes

import (
	"fmt"
	"github.com/dottics/dutil"
	security "github.com/dottics/securityserv"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/microtest"
	"testing"
)

func TestPasswordResetToken(t *testing.T) {
	tests := []struct {
		name     string
		payload  security.PasswordResetTokenPayload
		exchange *microtest.Exchange
		uuid     uuid.UUID
		e        dutil.Error
	}{
		{
			name:    "Bad Request",
			payload: security.PasswordResetTokenPayload{Email: "t@test.dottics.com"},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 400,
					Body:   `{"message":"BadRequest","data":null,"errors":{"user":["not found"]}}`,
				},
			},
			uuid: uuid.UUID{},
			e: &dutil.Err{
				Status: 400,
				Errors: map[string][]string{
					"user": {"not found"},
				},
			},
		},
		{
			name:    "UUID Parse Error",
			payload: security.PasswordResetTokenPayload{Email: "t@test.dottics.com"},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   `{"message":"password reset token successful","data":{"password_reset_token":"not-a-valid-uuid"},"errors":null}`,
				},
			},
			uuid: uuid.UUID{},
			e: &dutil.Err{
				Status: 400,
				Errors: map[string][]string{
					"uuid": {"invalid UUID length: 16"},
				},
			},
		},
		{
			name:    "Successful",
			payload: security.PasswordResetTokenPayload{Email: "t@test.dottics.com"},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   `{"message":"password reset token successful","data":{"password_reset_token":"3d9579de-162f-450a-8a7a-0cde0305d530"},"errors":null}`,
				},
			},
			uuid: uuid.MustParse("3d9579de-162f-450a-8a7a-0cde0305d530"),
			e:    nil,
		},
	}

	ms := microtest.NewMockServer("SECURITY_SERVICE_SCHEME", "SECURITY_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tests {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			ms.Append(tc.exchange)

			token, e := PasswordResetToken(tc.payload)
			if token != tc.uuid {
				t.Errorf("expected password reset token uuid %s got %s", tc.uuid.String(), token.String())
			}
			if !dutil.ErrorEqual(e, tc.e) {
				t.Errorf("expected error %v got %v", tc.e, e)
			}
		})
	}
}

func TestResetPassword(t *testing.T) {
	type E struct {
		e dutil.Error
	}
	tests := []struct {
		name     string
		payload  security.ResetPasswordPayload
		exchange *microtest.Exchange
		E        E
	}{
		{
			name: "Error",
			payload: security.ResetPasswordPayload{
				Email:              "name@example.com",
				PasswordResetToken: "this is a bad token",
				Password:           "new password",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 400,
					Body:   `{"message":"BadRequest","data":null,"errors":{"password_reset_token":["invalid"]}}`,
				},
			},
			E: E{
				e: &dutil.Err{
					Errors: map[string][]string{
						"password_reset_token": {"invalid"},
					},
				},
			},
		},
		{
			name: "Success",
			payload: security.ResetPasswordPayload{
				Email:              "name@example.com",
				PasswordResetToken: "3d9579de-162f-450a-8a7a-0cde0305d530",
				Password:           "new password",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   `{"message":"password reset successful","data":null,"errors":null}`,
				},
			},
			E: E{
				e: nil,
			},
		},
	}

	ms := microtest.NewMockServer("SECURITY_SERVICE_SCHEME", "SECURITY_SERVICE_HOST")
	defer ms.Server.Close()

	for i, tc := range tests {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			ms.Append(tc.exchange)

			e := ResetPassword(tc.payload)
			if !dutil.ErrorEqual(e, tc.E.e) {
				t.Errorf("expected error %v got %v", tc.E.e, e)
			}
		})
	}
}

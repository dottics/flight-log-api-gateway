package includes

import (
	"github.com/dottics/dutil"
	security "github.com/dottics/securityserv"
	"github.com/google/uuid"
)

// PasswordResetToken handles the exchange with the security microservice
// to generate the password reset token for the user.
func PasswordResetToken(p security.PasswordResetTokenPayload) (uuid.UUID, dutil.Error) {
	s := security.NewService("")
	t, e := s.PasswordResetToken(p)
	if e != nil {
		return uuid.UUID{}, e
	}
	tkn, err := uuid.Parse(t)
	if err != nil {
		e := dutil.NewErr(500, "uuid", []string{err.Error()})
		return uuid.UUID{}, e
	}
	return tkn, nil
	//return "", nil
}

// ResetPassword handles the exchange with the security microservice to
// reset a user's password.
func ResetPassword(p security.ResetPasswordPayload) dutil.Error {
	s := security.NewService("")
	return s.ResetPassword(p)
}

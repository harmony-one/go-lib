package errors

import (
	"github.com/pkg/errors"
)

var (
	// ErrMissingAccount is returned if the request keystore account can't be found (via its name)
	ErrMissingAccount = errors.New("keystore account can't be nil - please make sure the account you want to use exists in the keystore")
)

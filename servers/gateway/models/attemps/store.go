package attempts

import (
	"errors"
	"time"
)

// ErrAttemptNotFound is returned if no Attempt is found
// for a given email.
var ErrAttemptNotFound = errors.New("attempt not found")

// Store stores user failed sign-in attempts.
type Store interface {
	Save(email string, attempt *Attempt, expiry time.Duration) error

	Get(email string, attempt *Attempt) error

	Delete(email string) error
}

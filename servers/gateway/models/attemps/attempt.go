package attempts

import (
	"time"
)

// DefaultExpireTime represents a short amount of time
// within which the sign-in block will be triggered.
const DefaultExpireTime = time.Minute * 2

// BlockTime represents how long the user
// has to wait if failed sign-in attempts exceeds the max attempt
// for a given email address.
const BlockTime = time.Minute * 10

// MaxAttempt represents the max number of
// failed sign-in attempts for a given email.
const MaxAttempt = 5

// Attempt represents how many times the user has failed sign-in.
type Attempt struct {
	Count     int
	IsBlocked bool
}

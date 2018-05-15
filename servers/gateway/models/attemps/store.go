package attempts

import (
	"errors"
	"time"

	"github.com/info344-s18/challenges-yi904835116/servers/gateway/indexes"
)

// ErrAttemptNotFound is returned if no Attempt is found
// for a given email.
var AttemptNotFound = errors.New("attempt not found")

// Store stores user failed sign-in attempts.
type Store interface {
	Save(email string, attempt *Attempt, expiry time.Duration) error

	Get(email string, attempt *Attempt) error

	Delete(email string) error

	Trie() *indexes.Trie
}

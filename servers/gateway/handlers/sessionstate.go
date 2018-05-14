package handlers

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!
import (
	"time"

	"github.com/info344-s18/challenges-yi904835116/servers/gateway/models/users"
)

// SessionState store data for an authenticated user.
type SessionState struct {
	BeginTime time.Time
	User      *users.User
}

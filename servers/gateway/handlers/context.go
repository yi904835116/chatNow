package handlers

import (
	"github.com/info344-s18/challenges-yi904835116/servers/gateway/indexes"
	"github.com/info344-s18/challenges-yi904835116/servers/gateway/models/users"
	"github.com/info344-s18/challenges-yi904835116/servers/gateway/sessions"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store

// HandlerContext will be a receiver on any of your HTTP
// handler functions that need access to
// globals, such as the key used for signing
// and verifying SessionIDs, the session store
// and the user store.
type HandlerContext struct {
	SigningKey string
	// The type is an Store interface
	// rather than an actual Store implementation.
	SessionStore sessions.Store
	UserStore    users.Store

	Trie *indexes.Trie
	// AttemptStore attempts.Store
}

// NewHandlerContext constructs a new HanderContext,
// ensuring that the dependencies are valid values.
func NewHandlerContext(signingKey string, sessionStore sessions.Store, userStore users.Store, trie *indexes.Trie) *HandlerContext {

	if len(signingKey) == 0 {
		panic("signing key has length of zero")
	}

	if sessionStore == nil {
		panic("nil session store")
	}

	if userStore == nil {
		panic("nil user store")
	}

	if trie == nil {
		panic("nil trie tree")
	}

	return &HandlerContext{signingKey, sessionStore, userStore, trie}

}

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"

	"github.com/info344-s18/challenges-yi904835116/servers/gateway/models/users"
	"github.com/info344-s18/challenges-yi904835116/servers/gateway/sessions"
)

const headerUser = "X-User"

//NewServiceProxy returns a new ReverseProxy
//for a microservice given a comma-delimited
//list of network addresses
func NewServiceProxy(addrs string, ctx *HandlerContext) *httputil.ReverseProxy {
	splitAddrs := strings.Split(addrs, ",")
	nextAddr := 0
	mx := sync.Mutex{}

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			mx.Lock()
			r.URL.Host = splitAddrs[nextAddr]
			nextAddr = (nextAddr + 1) % len(splitAddrs)
			mx.Unlock()

			// get current authenticated user
			user := getCurrentUser(r, ctx)
			if user != nil {
				userJSON, _ := json.Marshal(user)
				r.Header.Add("X-User", string(userJSON))
			} else {
				// If there is no user authenticated user,
				// remote it from request.
				r.Header.Del("X-User")
			}
		},
	}
}

func getCurrentUser(r *http.Request, ctx *HandlerContext) *users.User {
	sessionState := &SessionState{}
	_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, sessionState)
	if err != nil {
		return nil
	}
	return sessionState.User
}

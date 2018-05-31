package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nbutton23/zxcvbn-go"

	"github.com/info344-s18/challenges-yi904835116/servers/gateway/models/users"
	"github.com/info344-s18/challenges-yi904835116/servers/gateway/sessions"
)

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.

const defualtSearchUserNumber = 20

// UsersHandler handles requests for the "users" resource.
func (context *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		sessionState := &SessionState{}
		_, err := sessions.GetState(r, context.SigningKey, context.SessionStore, sessionState)
		if err != nil {
			http.Error(w, fmt.Sprintf("error unauthenticated users: %v", err), http.StatusUnauthorized)
			return
		}

		results := []*users.User{}

		w.Header().Add(headerContentType, contentTypeJSON)

		query := r.URL.Query().Get("q")
		if len(query) == 0 {
			err = json.NewEncoder(w).Encode(results)
			if err != nil {
				http.Error(w, "error searching query should not be empty ", http.StatusBadRequest)
				return
			}
		}

		userSet := make(map[int64]bool)

		userSet = context.Trie.Search(20, query)

		results, err = context.UserStore.ConvertIDToUsers(userSet)

		sort.Slice(results, func(i, j int) bool {
			return results[i].ID < results[j].ID
		})

		w.Header().Add(headerContentType, contentTypeJSON)
		err = json.NewEncoder(w).Encode(results)
		if err != nil {
			http.Error(w, fmt.Sprintf("error encoding search results to JSON: %v", err), http.StatusInternalServerError)
			return
		}

	case "POST":
		// containJSON(r.Header.Get(headerContentType), w)

		if !strings.HasPrefix(r.Header.Get(headerContentType), contentTypeJSON) {
			http.Error(w, "request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}

		// Create an empty User to hold decoded request body.
		newUser := &users.NewUser{}

		err := json.NewDecoder(r.Body).Decode(newUser)

		// msValuePtr := reflect.ValueOf(&ms)
		// msValue := msValuePtr.Elem()

		// for

		zxcvbn.PasswordStrength(newUser.Password, userInput)
		if err != nil {
			http.Error(w, "error in JSON decoding. invalid JSON in request body", http.StatusBadRequest)
			return
		}

		// Validate the NewUser.
		err = newUser.Validate()

		if err != nil {
			http.Error(w, fmt.Sprintf("error validating new user: %s", err), http.StatusBadRequest)
			return
		}

		user, err := newUser.ToUser()

		if err != nil {
			http.Error(w, fmt.Sprintf("error converting new user to user: %s", err), http.StatusBadRequest)
			return
		}
		// Ensure there isn't already a user in the user store with the same email address.
		_, err = context.UserStore.GetByEmail(newUser.Email)
		if err == nil {
			http.Error(w, "user with the same email already exists", http.StatusBadRequest)
			return
		}

		// Ensure there isn't already a user in the user store with the same user name.
		_, err = context.UserStore.GetByUserName(newUser.UserName)
		if err == nil {
			http.Error(w, "user with the same username already exists", http.StatusBadRequest)
			return
		}

		// Insert the new user into the user store.
		user, err = context.UserStore.Insert(user)
		if err != nil {
			http.Error(w, fmt.Sprintf("error inserting new user: %s", err), http.StatusInternalServerError)
			return
		}

		context.Trie.Insert(user.UserName, user.ID)
		context.Trie.Insert(user.FirstName, user.ID)
		context.Trie.Insert(user.LastName, user.ID)

		beginNewSession(context, user, w)

	default:
		http.Error(w, "expect POST/GET method only", http.StatusMethodNotAllowed)

	}

}

// SpecificUserHandler handles requests for a specific user.
func (context *HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get session state from session store.

	sessionState := &SessionState{}
	sessionID, err := sessions.GetState(r, context.SigningKey, context.SessionStore, sessionState)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting session state: %v", err), http.StatusUnauthorized)
		return
	}

	path := path.Base(r.URL.Path)

	var givenID int64

	if path != "me" {
		givenID, err = strconv.ParseInt(path, 10, 64)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing ID: %v", err), http.StatusInternalServerError)
		return
	}

	switch r.Method {

	// Get the current user from the session state and respond with that user encoded as JSON object.
	case "GET":

		var user *users.User

		if path == "me" {
			user, err = sessionState.User, nil
		} else {
			user, err = context.UserStore.GetByID(givenID)
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("no user is found with given ID: %v", err), http.StatusNotFound)
			return
		}

		w.Header().Add(headerContentType, contentTypeJSON)
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, "error encoding SessionState Struct to JSON", http.StatusInternalServerError)
			return
		}

	// Update the current user with the JSON in the request body,
	// and respond with the newly updated user, encoded as a JSON object.
	case "PATCH":
		// Get Updates struct from request body.
		if path != "me" || givenID != sessionState.User.ID {
			http.Error(w, "User ID is not valid or does not match current-authenticaled user", http.StatusForbidden)
			return
		}
		// containJSON(r.Header.Get(headerContentType), w)

		if !strings.HasPrefix(r.Header.Get(headerContentType), contentTypeJSON) {
			http.Error(w, "request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}

		// Remove the user old fields from the trie.
		context.Trie.Remove(sessionState.User.FirstName, sessionState.User.ID)
		context.Trie.Remove(sessionState.User.LastName, sessionState.User.ID)

		updates := &users.Updates{}
		err := json.NewDecoder(r.Body).Decode(updates)
		if err != nil {
			http.Error(w, "error decoding request body: invalid JSON in request body", http.StatusBadRequest)
			return
		}
		// Update session state.
		sessionState.User.FirstName = updates.FirstName
		sessionState.User.LastName = updates.LastName

		// Update session store.
		err = context.SessionStore.Save(sessionID, sessionState)
		if err != nil {
			http.Error(w, fmt.Sprintf("error saving updated session state to session store: %s", err), http.StatusInternalServerError)
			return
		}

		// Insert the updated user fields into the trie.
		context.Trie.Insert(sessionState.User.FirstName, sessionState.User.ID)
		context.Trie.Insert(sessionState.User.LastName, sessionState.User.ID)

		// Update user store.
		user, err := context.UserStore.Update(sessionState.User.ID, updates)

		if err != nil {
			http.Error(w, fmt.Sprintf("error updating user store: %s", err), http.StatusInternalServerError)
			return
		}

		w.Header().Add(headerContentType, contentTypeJSON)
		w.WriteHeader(http.StatusOK)

		// err = json.NewEncoder(w).Encode(sessionState.User)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, "error encoding SessionState Struct to JSON", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "expect GET or PATCH method only", http.StatusMethodNotAllowed)
		return
	}
}

// SessionsHandler handles requests for the "sessions" resource, and allows clients to begin a new session using an existing user's credentials.
func (context *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	// Method must be POST.
	if r.Method != "POST" {
		http.Error(w, "expect POST method only", http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.Header.Get(headerContentType), contentTypeJSON) {
		http.Error(w, "request body must be in JSON", http.StatusUnsupportedMediaType)
		return
	}

	// Decode the request body into a users.Credentials struct.
	credentials := &users.Credentials{}
	err := json.NewDecoder(r.Body).Decode(credentials)
	if err != nil {
		http.Error(w, "error decoding request body: invalid JSON in request body", http.StatusBadRequest)
		return
	}

	// Get the user with the provided email from the UserStore.
	// If not found, respond with an http.StatusUnauthorized error
	// and the message "invalid credentials".
	user, err := context.UserStore.GetByEmail(credentials.Email)

	// Authenticate the user using the provided password.
	// If that fails, respond with an http.StatusUnauthorized error
	// and the message "invalid credentials".
	if err != nil {
		err = user.Authenticate(credentials.Password)
	}

	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	beginNewSession(context, user, w)
}

// SpecificSessionHandler handles requests related to a specific authenticated session.
func (context *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "expect DELETE method only", http.StatusMethodNotAllowed)
		return
	}

	path := path.Base(r.URL.Path)

	if path != "mine" {
		http.Error(w, "given path is not valid", http.StatusForbidden)
		return
	}

	_, err := sessions.EndSession(r, context.SigningKey, context.SessionStore)
	if err != nil {
		http.Error(w, fmt.Sprintf("error ending session: %s", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("signed out"))
}

// begineNewSession begins a new session process
func beginNewSession(context *HandlerContext, user *users.User, w http.ResponseWriter) {
	sessionState := SessionState{
		BeginTime: time.Now(),
		User:      user,
	}

	// begin new session and save seesion state to the store
	_, err := sessions.BeginSession(context.SigningKey, context.SessionStore, sessionState, w)
	if err != nil {
		http.Error(w, fmt.Sprintf("error beginning session: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Add(headerContentType, contentTypeJSON)
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "error encoding User struct to JSON", http.StatusInternalServerError)
		return
	}
}

func containJSON(contentType string, w http.ResponseWriter) {
	if !strings.HasPrefix(contentType, contentTypeJSON) {
		http.Error(w, "request body must be in JSON", http.StatusUnsupportedMediaType)
		return
	}
}

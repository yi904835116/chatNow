package handlers

// import (
// 	"bytes"
// 	"database/sql"
// 	"log"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/go-redis/redis"
// 	"github.com/go-sql-driver/mysql"
// 	"github.com/info344-s18/challenges-yi904835116/servers/gateway/models/users"
// 	"github.com/info344-s18/challenges-yi904835116/servers/gateway/sessions"
// )

// func config() *HandlerContext {

// 	redisClient := redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379",
// 	})

// 	sessionStore := sessions.NewRedisStore(redisClient, time.Hour)

// 	// mysqlAddr := os.Getenv("MYSQLADDR")

// 	config := mysql.Config{
// 		Addr:   "localhost:3306",
// 		User:   "root",
// 		Passwd: "abcd",
// 		Net:    "tcp",
// 		DBName: "mysqlDB",
// 	}
// 	db, err := sql.Open("mysql", config.FormatDSN())

// 	if err != nil {
// 		log.Fatalf("error opening database: %v", err)
// 	}

// 	defer db.Close()

// 	userStore := users.NewMySQLStore(db)
// 	context := NewHandlerContext("12345", sessionStore, userStore)
// 	return context
// }

// func TestUsersHandler(t *testing.T) {
// 	// setting up configeration
// 	context := config()
// 	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
// 	// pass 'nil' as the third parameter.
// 	var userInfo = `{"Email":"test@test.com"}{"Password":"123456789"}
// 	{"PasswordConf":"123456789"}{"UserName":"testuser"}
// 	{"FirstName":"test"}{"LastName":"user"}`
// 	var jsonStr = []byte(userInfo)

// 	req, err := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonStr))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(context.UsersHandler)

// 	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
// 	// directly and pass in our Request and ResponseRecorder.
// 	handler.ServeHTTP(rr, req)

// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != http.StatusCreated {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusCreated)
// 	}

// 	// Check the response body is what we expect.
// 	expected := userInfo
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}

// }

// func TestSpecificUserHandler(t *testing.T) {
// 	// setting up configeration
// 	context := config()
// 	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
// 	// pass 'nil' as the third parameter
// 	newUser := &users.NewUser{
// 		Email:        "test@test.com",
// 		Password:     "123456789",
// 		PasswordConf: "123456789",
// 		UserName:     "testuser",
// 		FirstName:    "test",
// 		LastName:     "user",
// 	}

// 	user, _ := newUser.ToUser()

// 	context.UserStore.Insert(user)

// 	// test GET http request
// 	var userInfo = `{"Email":"test@test.com"}{"Password":"123456789"}
// 	{"PasswordConf":"123456789"}{"UserName":"testuser"}
// 	{"FirstName":"test"}{"LastName":"user"}`
// 	var jsonStr = []byte(userInfo)

// 	req, err := http.NewRequest("GET", "/v1/users/"+user.ID, bytes.NewBuffer(jsonStr))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(context.UsersHandler)

// 	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
// 	// directly and pass in our Request and ResponseRecorder.
// 	handler.ServeHTTP(rr, req)

// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != http.StatusCreated {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusCreated)
// 	}

// 	// Check the response body is what we expect.
// 	expected := userInfo
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}

// 	// test PATCH http request
// 	var update = `{"Email":"test@test.com"}{"Password":"123456789"}
// 	{"PasswordConf":"123456789"}{"UserName":"testuser"}
// 	{"FirstName":"test"}{"LastName":"user"}`
// 	var jsonStr = []byte(userInfo)

// 	req, err := http.NewRequest("GET", "/v1/users/"+user.ID, bytes.NewBuffer(jsonStr))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(context.UsersHandler)

// 	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
// 	// directly and pass in our Request and ResponseRecorder.
// 	handler.ServeHTTP(rr, req)

// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != http.StatusCreated {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusCreated)
// 	}

// 	// Check the response body is what we expect.
// 	expected := userInfo
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}

// }

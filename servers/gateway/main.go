package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-sql-driver/mysql"

	"github.com/info344-s18/challenges-yi904835116/servers/gateway/handlers"
	"github.com/info344-s18/challenges-yi904835116/servers/gateway/models/users"
	"github.com/info344-s18/challenges-yi904835116/servers/gateway/sessions"
)

func reqEnv(name string) string {
	val := os.Getenv(name)
	if len(val) == 0 {
		log.Fatalf("please set the %s environment variable", name)
	}
	return val
}

//main is the main entry point for the server
func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	//get the TLS key and cert paths from environment variables
	//this allows us to use a self-signed cert/key during development
	//and the Let's Encrypt cert/key in production
	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")

	if len(tlsKeyPath) == 0 || len(tlsCertPath) == 0 {
		log.Fatal("Please set TLSCERT and TLSKEY")
	}
	// sessionKey is the signing key for SessionID.
	sessionKey := os.Getenv("SESSIONKEY")
	if len(sessionKey) == 0 {
		log.Fatal("Please set SESSIONKEY environment variable")
	}

	// Set up Redis connection.
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = "localhost:6379"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	sessionStore := sessions.NewRedisStore(redisClient, time.Hour)

	// dsn := os.Getenv("DSN")

	dbAddr := reqEnv("MYSQL_ADDR")
	dbPwd := reqEnv("MYSQL_ROOT_PASSWORD")
	dbName := reqEnv("MYSQL_DATABASE")

	config := mysql.Config{
		Addr:   dbAddr,
		User:   "root",
		Passwd: dbPwd,
		DBName: dbName,
		Net:    "tcp",
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatalf("error opening mysql: %v", err)

	}

	defer db.Close()

	userStore := users.NewMySQLStore(db)

	context := handlers.NewHandlerContext(sessionKey, sessionStore, userStore)
	// , attemptStore, resetCodeStore

	mux := http.NewServeMux()

	// mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)

	// Gateway
	mux.HandleFunc("/v1/users", context.UsersHandler)
	mux.HandleFunc("/v1/users/", context.SpecificUserHandler)

	mux.HandleFunc("/v1/sessions", context.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", context.SpecificSessionHandler)

	// mux.HandleFunc("/v1/resetcodes", context.ResetCodesHandler)
	// mux.HandleFunc("/v1/passwords", context.ResetPasswordHandler)

	corsMux := handlers.NewCORShandler(mux)

	log.Printf("Server is listening at https://%s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, corsMux))
}

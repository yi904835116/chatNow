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

	userStore := users.NewMySQLStore(db)

	// _, err = db.Query("select * from user")
	// if err != nil {
	// 	log.Fatalf("error select all: %v", err)
	// }

	trieTree, err := userStore.Trie()

	defer db.Close()

	if err != nil {
		log.Fatalf("error constructing user trie tree: %v", err)
	}

	context := handlers.NewHandlerContext(sessionKey, sessionStore, userStore, trieTree)

	// Messaging microservice addresses.
	msgAddrs := os.Getenv("MESSAGES_ADDR")
	if len(msgAddrs) == 0 {
		log.Fatal("Please set MESSAGES_ADDR environment variables")
	}

	// Summary microservice addresses.
	sumAddrs := os.Getenv("SUMMARYS_ADDR")
	if len(sumAddrs) == 0 {
		log.Fatal("Please set SUMMARYS_ADDR environment variables")
	}

	mux := http.NewServeMux()

	// Gateway of user authentication
	mux.HandleFunc("/v1/users", context.UsersHandler)
	mux.HandleFunc("/v1/users/", context.SpecificUserHandler)

	mux.HandleFunc("/v1/sessions", context.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", context.SpecificSessionHandler)

	// Summary microservice.
	mux.Handle("/v1/summary", handlers.NewServiceProxy(sumAddrs, context))
	// Messaging microservice.
	mux.Handle("/v1/channels", handlers.NewServiceProxy(msgAddrs, context))
	mux.Handle("/v1/channels/", handlers.NewServiceProxy(msgAddrs, context))
	mux.Handle("/v1/messages/", handlers.NewServiceProxy(msgAddrs, context))

	corsMux := handlers.NewCORShandler(mux)

	log.Printf("Server is listening at https://%s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, corsMux))
}

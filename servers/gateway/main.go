package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"

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

	notifier := handlers.NewNotifier()
	mux.Handle("/v1/ws", context.NewWebSocketsHandler(notifier))
	mqAddr := os.Getenv("MQADDR")
	if len(mqAddr) == 0 {
		log.Fatal("Please set the MQADDR variable to the address of your MQ server")
	}
	go listenToMQ(mqAddr, notifier)

	corsMux := handlers.NewCORShandler(mux)

	log.Printf("Server is listening at https://%s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, corsMux))
}

const qName = "testQ"
const maxConnRetries = 5

func listenToMQ(addr string, notifier *handlers.Notifier) {
	conn, err := connectToMQ(addr)
	if err != nil {
		log.Fatalf("error connecting to MQ server: %s", err)
	}
	log.Printf("connected to MQ server")
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("error opening channel: %v", err)
	}
	log.Println("created MQ channel")

	ch.Confirm(false)

	defer ch.Close()

	q, err := ch.QueueDeclare(qName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("error declaring queue: %v", err)
	}
	log.Println("declared MQ queue")
	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("error listening to queue: %v", err)
	}
	log.Println("listening for new MQ messages...")
	for msg := range messages {
		// log.Printf("new message id %s received from MQ", string(msg.Body))
		// Load messages received from RabbitMQ's eventQ channel to
		// notifier's eventQ channel, so that messages will be
		// broadcasted to all clients throught websocket.
		notifier.Notify(msg.Body)
	}
}

func connectToMQ(addr string) (*amqp.Connection, error) {
	mqURL := "amqp://" + addr
	var conn *amqp.Connection
	var err error
	for i := 1; i <= maxConnRetries; i++ {
		conn, err = amqp.Dial(mqURL)
		if err == nil {
			return conn, nil
		}
		log.Printf("error connecting to MQ server at %s: %s", mqURL, err)
		log.Printf("will attempt another connection in %d seconds", i*2)
		time.Sleep(time.Duration(i*2) * time.Second)
	}
	defer conn.Close()
	return nil, err
}

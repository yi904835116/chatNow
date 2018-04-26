package attempts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// RedisStore represents a attempts.Store backed by Redis.
type RedisStore struct {
	// Redis client used to talk to redis server.
	Client *redis.Client
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client) *RedisStore {
	//initialize and return a new RedisStore struct
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}

	return &RedisStore{
		Client: client,
	}
}

// Save saves the  email and corresponding attempts to the store.
func (rs *RedisStore) Save(email string, attempt *Attempt, expiry time.Duration) error {
	j, err := json.Marshal(attempt)
	if nil != err {
		return fmt.Errorf("error marshalling struct to JSON: %v", err)
	}

	err = rs.Client.Set(email, j, expiry).Err()
	if err != nil {
		return fmt.Errorf("error saving session state to Redis: %v", err)
	}

	return nil
}

// Get gets attempt with the data previously saved for the input email.
func (rs *RedisStore) Get(email string, attempt *Attempt) error {
	val, err := rs.Client.Get(email).Bytes()
	if err != nil {
		return AttemptNotFound
	}

	err = json.Unmarshal(val, attempt)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON to struct: %v", err)
	}

	return nil
}

// Delete deletes all Attempt data associated with the email from the store.
func (rs *RedisStore) Delete(email string) error {
	err := rs.Client.Del(email).Err()
	if err != nil {
		return fmt.Errorf("error deleting data: %v", err)
	}
	return nil
}

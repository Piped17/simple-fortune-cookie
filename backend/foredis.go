package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var dbLink redis.Conn
var usingRedis = false

// fortune struct and datastore from your previous code
type fortune struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type datastore struct {
	m map[string]fortune
	*sync.RWMutex
}

var datastoreDefault = datastore{
	m:       map[string]fortune{},
	RWMutex: &sync.RWMutex{},
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func init() {
	redisHost := getEnv("REDIS_DNS", "redis") // default to "redis" for Docker Compose
	redisPort := "6379"

	// Try connecting to Redis with retries
	var err error
	for i := 0; i < 5; i++ {
		dbLink, err = redis.Dial("tcp", fmt.Sprintf("%s:%s", redisHost, redisPort))
		if err == nil {
			usingRedis = true
			fmt.Println("Connected to Redis at", redisHost)
			break
		}
		log.Printf("Attempt %d: Redis connection failed: %s", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if !usingRedis {
		log.Println("Failed to connect to Redis after 5 attempts, using in-memory only")
		return
	}

	// Load fortunes from Redis into memory
	resKeys, err := redis.Strings(dbLink.Do("HKEYS", "fortunes"))
	if err != nil {
		log.Println("Redis HKEYS failed:", err)
		return
	}

	fmt.Println("*** Loading Redis fortunes:")
	datastoreDefault = datastore{m: map[string]fortune{}, RWMutex: &sync.RWMutex{}}
	for _, key := range resKeys {
		val, err := redis.String(dbLink.Do("HGET", "fortunes", key))
		if err != nil {
			log.Println("Redis HGET failed for key", key, ":", err)
			continue
		}
		datastoreDefault.m[key] = fortune{ID: key, Message: val}
		fmt.Printf("%s => %s\n", key, val)
	}
}

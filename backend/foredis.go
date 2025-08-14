package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

var dbLink redis.Conn
var usingRedis = false

// getEnv is already defined in main.go, so no need to redeclare it here

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
	for _, key := range resKeys {
		val, err := redis.String(dbLink.Do("HGET", "fortunes", key))
		if err != nil {
			log.Println("Redis HGET failed for key", key, ":", err)
			continue
		}
		// Store in the existing datastoreDefault
		datastoreDefault.Lock()
		datastoreDefault.m[key] = fortune{ID: key, Message: val}
		datastoreDefault.Unlock()
		fmt.Printf("%s => %s\n", key, val)
	}
}

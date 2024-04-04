package controllers

import (
	"database/sql"
	"log"

	"github.com/go-redis/redis/v8"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/e_store")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func connectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		count, err := GetCount()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "yay hits: %d", count)
	})

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	http.ListenAndServe(":5000", loggedRouter)
}

func GetCount() (int64, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	result, err := rdb.Incr("hits").Result()
	if err != nil {
		return 0, err
	}

	//fmt.Printf("hits, %#v", result)
	return result, nil
}

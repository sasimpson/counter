package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Counter struct {
	Hits int64 `json:"hits"`
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		count, err := GetCount()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if (r.Header.Get("Accepts") == "application/json") || (r.Header.Get("Content-Type") == "application/json") {
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(count)
			return
		}
		fmt.Fprintf(w, "yay hits: %d", count.Hits)
	})

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	http.ListenAndServe(":5000", loggedRouter)
}

func GetCount() (Counter, error) {
	var err error
	var result Counter

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	result.Hits, err = rdb.Incr("hits").Result()
	if err != nil {
		return result, err
	}

	return result, nil
}

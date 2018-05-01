package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	store, err := newConsulStore()
	if err != nil {
		log.Fatalf("unable to access store: %s", err)
	}

	listenPort := "8080"

	if val, ok := os.LookupEnv("KV_API_PORT"); ok {
		listenPort = val
	}

	log.Printf("Going to listen on port %s", listenPort)

	a := app{store: store}

	r := mux.NewRouter()
	r.HandleFunc("/stores/{store-name}", a.createStore).Methods("PUT")
	r.HandleFunc("/stores/{store-name}", a.deleteStore).Methods("DELETE")
	r.HandleFunc("/stores/{store-name}/keys/{key-name}", a.getKey).Methods("GET")
	r.HandleFunc("/stores/{store-name}/keys/{key-name}", a.deleteKey).Methods("DELETE")
	r.HandleFunc("/stores/{store-name}/keys/{key-name}", a.putKey).Methods("PUT")
	r.HandleFunc("/status", a.status).Methods("GET")
	r.HandleFunc("/", a.root).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+listenPort, r))
}

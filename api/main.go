package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	store, err := newConsulStore()
	if err != nil {
		log.Fatalf("unable to access store: %s", err)
	}

	a := app{store: store}

	r := mux.NewRouter()
	r.HandleFunc("/stores/{store-name}", a.createStore).Methods("PUT")
	r.HandleFunc("/stores/{store-name}", a.deleteStore).Methods("DELETE")
	r.HandleFunc("/stores/{store-name}/keys/{key-name}", a.getKey).Methods("GET")
	r.HandleFunc("/stores/{store-name}/keys/{key-name}", a.deleteKey).Methods("DELETE")
	r.HandleFunc("/stores/{store-name}/keys/{key-name}", a.putKey).Methods("PUT")
	r.HandleFunc("/status", a.status).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", r))
}

package main

import(
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/stores/{store-name}", createStore).Methods("PUT")
	r.HandleFunc("/stores/{store-name}", deleteStore).Methods("DELETE")
	r.HandleFunc("/stores/{store-name}/keys/{key-name}", getKey).Methods("GET")
	r.HandleFunc("/stores/{store-name}/keys/{key-name}", deleteKey).Methods("DELETE")
	r.HandleFunc("/stores/{store-name}/keys/{key-name}", putKey).Methods("PUT")
	http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":8080", r))
}

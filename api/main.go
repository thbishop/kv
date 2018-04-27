package main

import(
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)


func createStore(w http.ResponseWriter, r *http.Request) {
    name := mux.Vars(r)["store-name"]
    fmt.Printf("Got path var of: %s\n", name)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(`{ "name": "` + name + `" }`))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/stores/{store-name}", createStore)
	http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":8080", r))
}

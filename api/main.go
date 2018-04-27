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
    store := newEtcdStore()
    err := store.Create(name)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError) // TODO this should also handle client errors
        w.Write([]byte(`{ "error": ` + err.Error() + `}`))
        return
    }

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

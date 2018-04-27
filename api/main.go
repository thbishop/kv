package main

import(
    "fmt"
    "io/ioutil"
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

func putKey(w http.ResponseWriter, r *http.Request) {
    storeName := mux.Vars(r)["store-name"]
    key := mux.Vars(r)["key-name"]
    fmt.Printf("Got store '%s' and key '%s'\n", storeName, key)
    store := newEtcdStore()

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`{ "error": "error reading request body: ` + err.Error() + `" }`))
        return
    }

    err = store.SetKey(storeName, key, body)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError) // TODO this should also handle client errors
        w.Write([]byte(`{ "error": "` + err.Error() + `" }`))
        return
    }

    w.WriteHeader(http.StatusNoContent)
    return
}

func getKey(w http.ResponseWriter, r *http.Request) {
    storeName := mux.Vars(r)["store-name"]
    key := mux.Vars(r)["key-name"]
    fmt.Printf("Got store '%s' and key '%s'\n", storeName, key)
    store := newEtcdStore()

    data, err := store.GetKey(storeName, key)
    // TODO what if the key is not found?
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError) // TODO this should also handle client errors
        w.Write([]byte(`{ "error": "` + err.Error() + `" }`))
        return
    }

    w.Header().Set("Content-Type", "application/octet-stream")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
    return
}

func deleteKey(w http.ResponseWriter, r *http.Request) {
    storeName := mux.Vars(r)["store-name"]
    key := mux.Vars(r)["key-name"]
    fmt.Printf("Got store '%s' and key '%s'\n", storeName, key)
    store := newEtcdStore()

    err := store.DeleteKey(storeName, key)
    // TODO what if the key is not found?
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError) // TODO this should also handle client errors
        w.Write([]byte(`{ "error": "` + err.Error() + `" }`))
        return
    }

    w.WriteHeader(http.StatusNoContent)
    return
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/stores/{store-name}", createStore).Methods("PUT")
	r.HandleFunc("/stores/{store-name}/keys/{key-name}", getKey).Methods("GET")
	r.HandleFunc("/stores/{store-name}/keys/{key-name}", deleteKey).Methods("DELETE")
	r.HandleFunc("/stores/{store-name}/keys/{key-name}", putKey).Methods("PUT")
	http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":8080", r))
}

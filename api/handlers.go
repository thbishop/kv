package main

import (
    "io/ioutil"
    "log"
    "github.com/gorilla/mux"
    "net/http"
)

func createStore(w http.ResponseWriter, r *http.Request) {
    name := mux.Vars(r)["store-name"]
    log.Printf("Got path var of: %s\n", name)
    store, err := newEtcdStore(name)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError) // TODO this should also handle client errors
        w.Write([]byte(`{ "error": ` + err.Error() + `}`))
        return
    }

    err = store.Create()
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

func deleteStore(w http.ResponseWriter, r *http.Request) {
    name := mux.Vars(r)["store-name"]
    log.Printf("Got path var of: %s\n", name)
    store, err := newEtcdStore(name)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError) // TODO this should also handle client errors
        w.Write([]byte(`{ "error": ` + err.Error() + `}`))
        return
    }

    err = store.Delete()
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError) // TODO this should also handle client errors
        w.Write([]byte(`{ "error": ` + err.Error() + `}`))
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func putKey(w http.ResponseWriter, r *http.Request) {
    name := mux.Vars(r)["store-name"]
    key := mux.Vars(r)["key-name"]
    log.Printf("Got store '%s' and key '%s'\n", name, key)

    store, err := newEtcdStore(name)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError) // TODO this should also handle client errors
        w.Write([]byte(`{ "error": ` + err.Error() + `}`))
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`{ "error": "error reading request body: ` + err.Error() + `" }`))
        return
    }

    err = store.SetKey(key, body)
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
    name := mux.Vars(r)["store-name"]
    key := mux.Vars(r)["key-name"]
    log.Printf("Got store '%s' and key '%s'\n", name, key)

    store, err := newEtcdStore(name)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError) // TODO this should also handle client errors
        w.Write([]byte(`{ "error": ` + err.Error() + `}`))
        return
    }

    data, err := store.GetKey(key)
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
    name := mux.Vars(r)["store-name"]
    key := mux.Vars(r)["key-name"]
    log.Printf("Got store '%s' and key '%s'\n", name, key)

    store, err := newEtcdStore(name)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError) // TODO this should also handle client errors
        w.Write([]byte(`{ "error": ` + err.Error() + `}`))
        return
    }

    err = store.DeleteKey(key)
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

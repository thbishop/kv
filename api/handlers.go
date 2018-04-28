package main

import (
    "io/ioutil"
    "log"
    "github.com/gorilla/mux"
    "net/http"
)

type requestInfo struct {
    storeName string
    keyName string
}

func newRequestInfo(r *http.Request) requestInfo {
    rinfo := requestInfo{
        keyName: mux.Vars(r)["key-name"],
        storeName: mux.Vars(r)["store-name"],
    }
    log.Printf("Request info: %+v\n", rinfo)
    return rinfo
}

func respondWithServerError(w http.ResponseWriter, err error) {
    log.Printf("Responding with 500: %s", err)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusInternalServerError)
}

type app struct {
    store Store
}

func (a *app) createStore(w http.ResponseWriter, r *http.Request) {
    rinfo := newRequestInfo(r)

    err := a.store.CreateStore(rinfo.storeName)
    if err != nil {
        respondWithServerError(w, err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(`{ "name": "` + rinfo.storeName + `" }`))
}

func (a *app) deleteStore(w http.ResponseWriter, r *http.Request) {
    rinfo := newRequestInfo(r)

    err := a.store.DeleteStore(rinfo.storeName)
    if err != nil {
        respondWithServerError(w, err)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (a *app) putKey(w http.ResponseWriter, r *http.Request) {
    rinfo := newRequestInfo(r)

    // TODO enforce data size limit
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        respondWithServerError(w, err)
        return
    }

    err = a.store.SetKey(rinfo.storeName, rinfo.keyName, body)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError) // TODO this should also handle client errors
        w.Write([]byte(`{ "error": "` + err.Error() + `" }`))
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (a *app) getKey(w http.ResponseWriter, r *http.Request) {
    rinfo := newRequestInfo(r)

    data, err := a.store.GetKey(rinfo.storeName, rinfo.keyName)
    // TODO handle if key is not found?
    if err != nil {
        respondWithServerError(w, err)
        return
    }

    w.Header().Set("Content-Type", "application/octet-stream")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
}

func (a *app) deleteKey(w http.ResponseWriter, r *http.Request) {
    rinfo := newRequestInfo(r)

    err := a.store.DeleteKey(rinfo.storeName, rinfo.keyName)
    // TODO what if the key is not found?
    if err != nil {
        respondWithServerError(w, err)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

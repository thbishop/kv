package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

type requestInfo struct {
	storeName string
	keyName   string
}

func newRequestInfo(r *http.Request) requestInfo {
	rinfo := requestInfo{
		keyName:   mux.Vars(r)["key-name"],
		storeName: mux.Vars(r)["store-name"],
	}
	log.Printf("Request info: %+v\n", rinfo)
	return rinfo
}

func respondWithBadRequest(w http.ResponseWriter, val validation) {
	log.Printf("Responding with %d due to validation error: %+v", http.StatusBadRequest, val)

	errResp := clientErrorResponse{Error: val.message}
	b, err := json.Marshal(errResp)
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write(b)
}

func respondWithNotFound(w http.ResponseWriter, reason string) {
	log.Printf("Responding with %d: %s", http.StatusNotFound, reason)
	w.WriteHeader(http.StatusNotFound)
}

func respondWithPayloadToLarge(w http.ResponseWriter, size int64) {
	log.Printf("Responding with %d size is %d bytes", http.StatusRequestEntityTooLarge, size)
	w.WriteHeader(http.StatusRequestEntityTooLarge)
}

func respondWithServerError(w http.ResponseWriter, err error) {
	log.Printf("Responding with %d: %s", http.StatusInternalServerError, err)
	w.WriteHeader(http.StatusInternalServerError)
}

type validation struct {
	valid   bool
	message string
}

type clientErrorResponse struct {
	Error string `json:"error"`
}

func validateAlphanumericString(s string) (validation, error) {
	matched, err := regexp.MatchString("^[A-z0-9-_]+$", s)
	if err != nil {
		log.Printf("Unable to validate alphanumeric string; err: %s\n", err)
		return validation{}, err
	}

	v := validation{valid: matched}
	if !v.valid {
		v.message = fmt.Sprintf("invalid value '%s'; expected alphanumeric or '-'", s)
	}

	return v, nil
}

func validateLength(s string, min int, max int) validation {
	if len(s) >= min && len(s) <= max {
		return validation{valid: true}
	}
	return validation{
		valid:   false,
		message: fmt.Sprintf("invalid value '%s'; expected length of %d to %d", s, min, max),
	}
}

type app struct {
	store Store
}

func (a *app) createStore(w http.ResponseWriter, r *http.Request) {
	rinfo := newRequestInfo(r)

	v, err := validateAlphanumericString(rinfo.storeName)
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	if !v.valid {
		respondWithBadRequest(w, v)
		return
	}

	if v = validateLength(rinfo.storeName, 3, 16); !v.valid {
		respondWithBadRequest(w, v)
		return
	}

	exists, err := a.store.StoreExists(rinfo.storeName)
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	if exists {
		v = validation{message: "store already exists"}
		respondWithBadRequest(w, v)
		return
	}

	err = a.store.CreateStore(rinfo.storeName)
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

	exists, err := a.store.StoreExists(rinfo.storeName)
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	if !exists {
		respondWithNotFound(w, "store does not exist")
		return
	}

	err = a.store.DeleteStore(rinfo.storeName)
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *app) putKey(w http.ResponseWriter, r *http.Request) {
	rinfo := newRequestInfo(r)

	if v := validateLength(rinfo.keyName, 3, 16); !v.valid {
		respondWithBadRequest(w, v)
		return
	}

	exists, err := a.store.StoreExists(rinfo.storeName)
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	if !exists {
		respondWithNotFound(w, "store does not exist")
		return
	}

	v, err := validateAlphanumericString(rinfo.keyName)
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	if !v.valid {
		respondWithBadRequest(w, v)
		return
	}

	if r.ContentLength > 2048 {
		respondWithPayloadToLarge(w, r.ContentLength)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	err = a.store.SetKey(rinfo.storeName, rinfo.keyName, body)
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *app) getKey(w http.ResponseWriter, r *http.Request) {
	rinfo := newRequestInfo(r)

	exists, err := a.store.KeyExists(rinfo.storeName, rinfo.keyName)
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	if !exists {
		respondWithNotFound(w, "key does not exist")
		return
	}

	data, err := a.store.GetKey(rinfo.storeName, rinfo.keyName)
	if err != nil {
		if a.store.IsKeyMissing(err) {
			respondWithNotFound(w, "key does not exist")
			return
		}

		respondWithServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (a *app) deleteKey(w http.ResponseWriter, r *http.Request) {
	rinfo := newRequestInfo(r)

	exists, err := a.store.KeyExists(rinfo.storeName, rinfo.keyName)
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	if !exists {
		respondWithNotFound(w, "key does not exist")
		return
	}

	err = a.store.DeleteKey(rinfo.storeName, rinfo.keyName)
	if err != nil {
		respondWithServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *app) status(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *app) root(w http.ResponseWriter, r *http.Request) {
	doc := `
<!doctype html>
<!DOCTYPE html>
<html>
<body>
<h2>Key/Value Service</h2>
<br/>
<p>
Visit the <a href="https://github.com/thbishop/kv">github repo</a> for details
</p>
<body>
</html>
`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(doc))
}

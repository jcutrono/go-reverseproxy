package main

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
)

type KeyService struct{}

var (
	siteKeys map[string]*httputil.ReverseProxy
	urlKeys  map[string]string
)

// Configure - sets up all TenantService routes.
func (svc *KeyService) Configure(router *mux.Router) {
	siteKeys = make(map[string]*httputil.ReverseProxy)
	urlKeys = make(map[string]string)

	router.HandleFunc("/debug", svc.debug).Methods("GET")

	router.HandleFunc("/key/{id}", svc.get).Methods("GET")
	router.HandleFunc("/key", basicAuth(svc.setup)).Methods("POST")
	router.HandleFunc("/key", basicAuth(svc.regen)).Methods("PUT")
	router.HandleFunc("/key", basicAuth(svc.remove)).Methods("DELETE")
}

func (svc *KeyService) debug(resp http.ResponseWriter, req *http.Request) {

	keys, _ := json.Marshal(urlKeys)
	resp.Write(keys)
}

func (svc *KeyService) setup(resp http.ResponseWriter, req *http.Request) {

	url := req.Referer()

	siteKey, _ := uuid.NewV4()

	siteKeys[siteKey.String()] = CreateReverseProxy(url)
	urlKeys[url] = siteKey.String()

	resp.Write([]byte("https://{change on use}/?key=" + siteKey.String()))

	resp.WriteHeader(http.StatusOK)
}

func (svc *KeyService) get(resp http.ResponseWriter, req *http.Request) {

	url := req.Referer()
	key := mux.Vars(req)["id"]

	if key != "" && urlKeys[url] == key {
		resp.WriteHeader(http.StatusOK)
	}

	resp.WriteHeader(http.StatusUnauthorized)
}

func (svc *KeyService) regen(resp http.ResponseWriter, req *http.Request) {

	svc.remove(resp, req)
	svc.setup(resp, req)
}

func (svc *KeyService) remove(resp http.ResponseWriter, req *http.Request) {

	url := req.Referer()
	key := urlKeys[url]

	delete(siteKeys, key)
	delete(urlKeys, url)

	resp.WriteHeader(http.StatusOK)
}

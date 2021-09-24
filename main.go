package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	initializeApi(router)

	http.Handle("/api/", router)

	http.Handle("/", UrlInspector(http.FileServer(http.Dir("./web/"))))

	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	fmt.Println("Begin HTTP Listen on " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println(err)
	}
}

func initializeApi(router *mux.Router) {

	// setup api grouping
	apiRoutes := router.PathPrefix("/api").Subrouter()

	// setup api controllers
	keysvc := KeyService{}
	keysvc.Configure(apiRoutes)

	apiRoutes.Headers("Access-Control-Allow-Origin", "*")
	apiRoutes.Headers("Content-Type", "application/json")

	fmt.Println("API Router initialized")
}

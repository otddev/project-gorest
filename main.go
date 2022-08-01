package main

import (
	"github.com/angarcia/gorest/handlers"
	"github.com/angarcia/gorest/mw"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	// Port Configuration & HTTP Logger Initiate
	var port string = ":8181"
	router := mux.NewRouter().StrictSlash(true)
	router.Use(mw.HTTPLogger)

	// Create handler which will also initialize empty map for storage.
	rh := handlers.CreateHandler(make(map[string]map[string]interface{}))

	// API Route Definitions
	api := router.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/resources/{id}", rh.GetResourceHandler).Methods(http.MethodGet)
	api.HandleFunc("/resources", rh.GetResourcesHandler).Methods(http.MethodGet)
	api.HandleFunc("/resources", rh.CreateResourceHandler).Methods(http.MethodPost)
	api.HandleFunc("/resources/{id}", rh.UpdateResourceHandler).Methods(http.MethodPut)
	api.HandleFunc("/resources/{id}", rh.DeleteResourceHandler).Methods(http.MethodDelete)

	// Page Not Found Route Definition
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Error: The page request was not found.")
		http.Error(w, "PAGE NOT FOUND", http.StatusNotFound)
	})

	// Server Start
	log.Printf("Starting Server: '%s'", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}
}

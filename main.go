package main

import (
	"github.com/angarcia/gorest/controllers"
	"github.com/angarcia/gorest/docs"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func init() {
	// Initialize Swagger Documentation
	////////////////////////////////////////////////////////
	docs.SwaggerInfo.Title = "GOLang RestFul API Example"
	docs.SwaggerInfo.Description = "This is the API documentation for GOLang RestFul API Example Project"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}
	////////////////////////////////////////////////////////
}

func main() {

	var port string = ":8080"

	db := make(map[string]map[string]interface{})
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api/").Subrouter()

	api.HandleFunc("/resources/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.FindResource(w, r, db)
	}).Methods("GET")

	api.HandleFunc("/resources", func(w http.ResponseWriter, r *http.Request) {
		controllers.FindResources(w, r, db)
	}).Methods("GET")

	api.HandleFunc("/resources", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateResource(w, r, db)
	}).Methods("POST")

	api.HandleFunc("/resources/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateResource(w, r, db)
	}).Methods("PUT")

	api.HandleFunc("/resources/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteResource(w, r, db)
	}).Methods("DELETE")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Error: The page request was not found.")
		http.Error(w, "PAGE NOT FOUND", http.StatusNotFound)
	})

	log.Printf("Starting Server: '%s'", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}

}

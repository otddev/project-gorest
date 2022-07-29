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

	db := make(map[string]map[string]interface{})
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api/").Subrouter()

	// FindResource godoc
	// @Summary Find a specific resource based on ID.
	// @Description Find and get a response with specific resources based on the ID provided on request.
	// @Tags Resources
	// @Produce  json
	// @Success 200 {object} models.HTTPSuccess
	// @Failure 400 {object} models.HTTPBadRequest
	// @Router /resource [get]

	api.HandleFunc("/resources/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.FindResource(w, r, db)
	}).Methods("GET")

	// FindResources godoc
	// @Summary Find all the resources created within service.
	// @Description Find and get a response with full list of resources created withing the service by ID.
	// @Tags Resources
	// @Produce  json
	// @Success 200 {object} models.HTTPSuccess
	// @Failure 400 {object} models.HTTPBadRequest
	// @Router /resources/{id} [get]

	api.HandleFunc("/resources", func(w http.ResponseWriter, r *http.Request) {
		controllers.FindResources(w, r, db)
	}).Methods("GET")

	// CreateResource godoc
	// @Summary Create a new resource/hashmap within service.
	// @Description Create a new resource/hashmap within service for later consumption.
	// @Tags Resources
	// @Accept json
	// @Produce  json
	// @Success 201 {object} models.HTTPSuccess
	// @Failure 400 {object} models.HTTPBadRequest
	// @Failure 500 {object} models.HTTPServerErr
	// @Router /resources [post]

	api.HandleFunc("/resources", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateResource(w, r, db)
	}).Methods("POST")

	// UpdateResource godoc
	// @Summary Update an existing hashmap resource based on Id.
	// @Description Update an existing resource/hashmap within service based on id..
	// @Tags Resources
	// @Accept json
	// @Produce json
	// @Success 201 {object} models.HTTPSuccess
	// @Failure 400 {object} models.HTTPBadRequest
	// @Failure 500 {object} models.HTTPServerErr
	// @Router /resources/{id} [put]

	api.HandleFunc("/resources/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateResource(w, r, db)
	}).Methods("PUT")

	// DeleteResource godoc
	// @Summary Delete an existing hashmap resource based on id.
	// @Description Delete an existing resource/hashmap within service based on id.
	// @Tags Resources
	// @Accept json
	// @Produce json
	// @Success 201 {object} models.HTTPSuccess
	// @Failure 400 {object} models.HTTPBadRequest
	// @Failure 500 {object} models.HTTPServerErr
	// @Router /resources/{id} [delete]

	api.HandleFunc("/resources/{id}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteResource(w, r, db)
	}).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":80", router))

}

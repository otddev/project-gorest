package controllers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// CheckID | The functions allows the check if a correct key string was provided is correct and if so check if exists.
func CheckID(i string, db map[string]map[string]interface{}) error {

	_, err := uuid.Parse(i)
	if err != nil {
		err.Error()
	}

	if _, ok := db[i]; ok {
		errors.New("the id provided does not exist in database")
	}

	return nil
}

// FindResource godoc
// @Description Find a specific resources based on the ID provided on request.
// @Tags Resources
// @Produce  json
// @Success 200 {object} models.HTTPSuccess
// @Failure 400 {object} models.HTTPBadRequest
// @Router /resource [get]

func FindResource(w http.ResponseWriter, r *http.Request, db map[string]map[string]interface{}) {

	i := mux.Vars(r)["id"]
	err := CheckID(i, db)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(db[i])

	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Resource Returned: %v\n", len(db[i]))
}

// FindResources godoc
// @Summary Find all the resources created within service.
// @Description Find and get a response with full list of resources created withing the service by ID.
// @Tags Resources
// @Produce  json
// @Success 200 {object} models.HTTPSuccess
// @Failure 400 {object} models.HTTPBadRequest
// @Router /resources/{id} [get]

func FindResources(w http.ResponseWriter, r *http.Request, db map[string]map[string]interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(db)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Number of Maps Returned: %v\n", len(db))
}

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

func CreateResource(w http.ResponseWriter, r *http.Request, db map[string]map[string]interface{}) {

	obj := make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, "missing json body in request or invalid", http.StatusBadRequest)
		return
	}

	i, err := uuid.NewRandom()
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db[i.String()] = obj
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(obj)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("New Map Created: %v\n", obj)
}

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

func UpdateResource(w http.ResponseWriter, r *http.Request, db map[string]map[string]interface{}) {

	i := mux.Vars(r)["id"]
	err := CheckID(i, db)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	obj := make(map[string]interface{})

	err = json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db[i] = obj
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	err = json.NewEncoder(w).Encode(obj)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Map Resource Updated: %v\n", obj)
}

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

func DeleteResource(w http.ResponseWriter, r *http.Request, db map[string]map[string]interface{}) {

	i := mux.Vars(r)["id"]
	err := CheckID(i, db)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	delete(db, i)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	log.Printf("Map Resource Deleted: %v\n", i)
}

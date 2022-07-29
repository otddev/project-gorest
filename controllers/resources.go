package controllers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func FindResource(w http.ResponseWriter, r *http.Request, db map[string]map[string]interface{}) {

	i := mux.Vars(r)["id"]

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(db[i])

	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Resource Returned: %v\n", len(db[i]))
}

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

func CreateResource(w http.ResponseWriter, r *http.Request, db map[string]map[string]interface{}) {

	obj := make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func UpdateResource(w http.ResponseWriter, r *http.Request, db map[string]map[string]interface{}) {

	i := mux.Vars(r)["id"]
	obj := make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&obj)
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

func DeleteResource(w http.ResponseWriter, r *http.Request, db map[string]map[string]interface{}) {

	i := mux.Vars(r)["id"]

	if _, ok := db[i]; ok {
		delete(db, i)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		log.Printf("Map Resource Deleted: %v\n", i)
		return
	}

	log.Printf("Error: resource does not exist")
	http.Error(w, "resource does not exist", http.StatusInternalServerError)
}

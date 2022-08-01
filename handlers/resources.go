package handlers

import (
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

// ResourceHandler contains resource handler data
type ResourceHandler struct {
	ch *CommonHandler
	dh *DBHelper
}

// CreateHandler creation/initialization of resource handler.
func CreateHandler(db map[string]map[string]interface{}) *ResourceHandler {
	return &ResourceHandler{
		ch: &CommonHandler{Marshaler: nil, Unmarshaler: nil},
		dh: &DBHelper{
			db: db,
			mu: sync.Mutex{},
		},
	}
}

// CheckID | The functions allows the check if a correct key string was provided is correct and if so check is exists.
func CheckID(i string, db map[string]map[string]interface{}) error {

	_, err := uuid.Parse(i)
	if err != nil {
		return err
	}

	if _, ok := db[i]; ok {
		return nil
	}
	return errors.New("the id provided does not exist in database")
}

// GetResourcesHandler GET /api/resources/
func (rh *ResourceHandler) GetResourcesHandler(w http.ResponseWriter, r *http.Request) {
	rh.dh.mu.Lock()
	defer rh.dh.mu.Unlock()

	if len(rh.dh.db) > 0 {
		data, err := rh.ch.Marshal(rh.dh.db)
		if err != nil {
			log.Printf("error: %v", err)
			rh.ch.HttpError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(data)
		if err != nil {
			log.Printf("error: %v", err)
			rh.ch.HttpError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Resources Returned: %v\n", len(rh.dh.db))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Resources Returned: %v\n", len(rh.dh.db))
	return
}

// GetResourceHandler GET /api/resources/{id}
func (rh *ResourceHandler) GetResourceHandler(w http.ResponseWriter, r *http.Request) {
	rh.dh.mu.Lock()
	defer rh.dh.mu.Unlock()

	i := mux.Vars(r)["id"]
	err := CheckID(i, rh.dh.db)
	if err != nil {
		log.Printf("error: %v", err)
		rh.ch.HttpError(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := rh.ch.Marshal(rh.dh.db[i])
	if err != nil {
		log.Printf("error: %v", err)
		rh.ch.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(data)
	if err != nil {
		log.Printf("error: %v", err)
		rh.ch.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Resource Returned: %v\n", len(rh.dh.db[i]))
	return

}

// CreateResourceHandler POST /api/resources/
func (rh *ResourceHandler) CreateResourceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rh.dh.mu.Lock()
	defer rh.dh.mu.Unlock()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: %v", err)
		rh.ch.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	obj := make(map[string]interface{})

	err = rh.ch.Unmarshal(b, &obj)
	if err != nil {
		log.Printf("error: %v", err)
		rh.ch.HttpError(w, err.Error(), http.StatusBadRequest)
		return
	}

	i := uuid.New().String()
	rh.dh.db[i] = obj

	data, err := rh.ch.Marshal(rh.dh.db[i])
	if err != nil {
		log.Printf("error: %v", err)
		rh.ch.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(data)
	if err != nil {
		log.Printf("error: %v", err)
		rh.ch.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Resource Created: Id: %v\n", i)
	return
}

// UpdateResourceHandler PUT /api/resources/{id}
func (rh *ResourceHandler) UpdateResourceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rh.dh.mu.Lock()
	defer rh.dh.mu.Unlock()

	i := mux.Vars(r)["id"]
	err := CheckID(i, rh.dh.db)
	if err != nil {
		log.Printf("error: %v", err)
		rh.ch.HttpError(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: %v", err)
		rh.ch.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	obj := make(map[string]interface{})

	err = rh.ch.Unmarshal(b, &obj)
	if err != nil {
		log.Printf("error: %v", err)
		rh.ch.HttpError(w, err.Error(), http.StatusBadRequest)
		return
	}

	rh.dh.db[i] = obj

	data, err := rh.ch.Marshal(rh.dh.db[i])
	if err != nil {
		log.Printf("error: %v", err)
		rh.ch.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	_, err = w.Write(data)
	if err != nil {
		log.Printf("error: %v", err)
		rh.ch.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Map Resource Updated: %v\n", i)
	return
}

// DeleteResourceHandler DELETE /api/resources/{id}
func (rh *ResourceHandler) DeleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rh.dh.mu.Lock()
	defer rh.dh.mu.Unlock()

	i := mux.Vars(r)["id"]
	err := CheckID(i, rh.dh.db)
	if err != nil {
		log.Printf("error: %v", err)
		rh.ch.HttpError(w, err.Error(), http.StatusBadRequest)
		return
	}

	delete(rh.dh.db, i)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	log.Printf("Map Resource Deleted: %v\n", i)
	return
}

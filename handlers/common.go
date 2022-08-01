package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorHttp struct {
	Status int    `json:"status"`
	Msg    string `json:"message"`
}

// CommonHandler Main struct for basic dependencies for API handlers
type CommonHandler struct {
	Marshaler   func(v interface{}) ([]byte, error)
	Unmarshaler func(data []byte, v interface{}) error
}

// Marshal Will marshal provided data with Marshaler defined in ch.
// If un-set, json.Marshal will be used.
func (ch *CommonHandler) Marshal(v interface{}) ([]byte, error) {
	if ch.Marshaler == nil {
		return json.Marshal(v)
	}
	return ch.Marshaler(v)
}

// Unmarshal Will unmarshal provided data with the Unmarshaler defined on ch.
// If un-set, json.Unmarshal will be used.
func (ch *CommonHandler) Unmarshal(data []byte, v interface{}) error {
	if ch.Unmarshaler == nil {
		return json.Unmarshal(data, &v)
	}
	return ch.Unmarshaler(data, &v)
}

// HttpError Custom error function to report HTTP request errors in application/json format instead of test/plain
func (ch *CommonHandler) HttpError(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorHttp{Status: code, Msg: err})
}

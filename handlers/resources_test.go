package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("failed reader")
}

// TestResourceHandler_CreateResourceHandler POST /api/resources/
func TestResourceHandler_CreateResourceHandler(t *testing.T) {
	type args struct {
		w    *httptest.ResponseRecorder
		body string
		m    func(v interface{}) ([]byte, error)
		u    func(data []byte, v interface{}) error
		db   map[string]map[string]interface{}
		want int
	}
	tests := []struct {
		name   string
		args   args
		reader io.Reader
	}{
		{
			name:   "CreateResource - Success",
			reader: strings.NewReader(`{"name":"Bruce","lastname":"Wayne"}`),
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{},
				want: 201,
			},
		},
		{
			name:   "CreateResource - JSON Marshal Failure",
			reader: strings.NewReader(`{"name":"Bruce","lastname":"Wayne"}`),
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{},
				m:    func(v interface{}) ([]byte, error) { return []byte{}, errors.New("fake error") },
				u:    func(data []byte, v interface{}) error { return json.Unmarshal(data, &v) },
				want: 500,
			},
		},
		{
			name:   "CreateResource - Empty Request Failure",
			reader: strings.NewReader(``),
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{},
				want: 400,
			},
		},
		{
			name:   "CreateResource - Invalid Request Failure",
			reader: strings.NewReader(`text`),
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{},
				want: 400,
			},
		},
		{
			name:   "CreateResource - Failed ioutil.ReadAll()",
			reader: errReader(0),
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{},
				want: 500,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup Request
			reader := io.Reader(nil)
			if tt.reader != nil {
				reader = tt.reader
			}

			rh := ResourceHandler{
				ch: &CommonHandler{
					Marshaler:   tt.args.m,
					Unmarshaler: tt.args.u,
				},
				dh: &DBHelper{
					db: tt.args.db,
					mu: sync.Mutex{},
				},
			}
			r, err := http.NewRequest("POST", "", reader)
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Content-Length", "1")
			if err != nil {
				t.Fail()
			}

			rh.CreateResourceHandler(tt.args.w, r)
			if tt.args.w.Code != tt.args.want {
				t.Errorf("got %d want %d", tt.args.w.Code, tt.args.want)
			}
		})
	}
}

// TestResourceHandler_GetResourcesHandler GET /api/resources/
func TestResourceHandler_GetResourcesHandler(t *testing.T) {
	type args struct {
		w    *httptest.ResponseRecorder
		body string
		m    func(v interface{}) ([]byte, error)
		u    func(data []byte, v interface{}) error
		db   map[string]map[string]interface{}
		want int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "GetResources - Success",
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				want: 200,
			},
		},
		{
			name: "GetResources - Success / No Content",
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{},
				want: 204,
			},
		},
		{
			name: "GetResources - JSON Marshal Failure",
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				m:    func(v interface{}) ([]byte, error) { return []byte{}, errors.New("fake error") },
				want: 500,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup Request
			rh := ResourceHandler{
				ch: &CommonHandler{
					Marshaler:   tt.args.m,
					Unmarshaler: tt.args.u,
				},
				dh: &DBHelper{
					db: tt.args.db,
					mu: sync.Mutex{},
				},
			}
			r, err := http.NewRequest("GET", "/", strings.NewReader(``))
			r.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fail()
			}

			rh.GetResourcesHandler(tt.args.w, r)
			if tt.args.w.Code != tt.args.want {
				t.Errorf("got %d want %d", tt.args.w.Code, tt.args.want)
			}
		})
	}
}

// TestResourceHandler_GetResourceHandler GET /api/resources/{id}
func TestResourceHandler_GetResourceHandler(t *testing.T) {
	type args struct {
		w    *httptest.ResponseRecorder
		body string
		m    func(v interface{}) ([]byte, error)
		u    func(data []byte, v interface{}) error
		db   map[string]map[string]interface{}
		want int
		vars map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "GetResource - Success",
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				vars: map[string]string{"id": "0bf8651a-0923-47b8-aed3-e9fc1505e497"},
				want: 200,
			},
		},
		{
			name: "GetResource - Invalid UUID Failure",
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				vars: map[string]string{"id": "dummy"},
				want: 400,
			},
		},
		{
			name: "GetResource - UUID Not Exist Failure",
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				vars: map[string]string{"id": "0bf8651a-0923-47b8-aed3-e9fc1505e496"},
				want: 400,
			},
		},
		{
			name: "GetResource - JSON Marshal Failure",
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				vars: map[string]string{"id": "0bf8651a-0923-47b8-aed3-e9fc1505e497"},
				m:    func(v interface{}) ([]byte, error) { return []byte{}, errors.New("fake error") },
				want: 500,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup Request
			rh := ResourceHandler{
				ch: &CommonHandler{
					Marshaler:   tt.args.m,
					Unmarshaler: tt.args.u,
				},
				dh: &DBHelper{
					db: tt.args.db,
					mu: sync.Mutex{},
				},
			}

			r, err := http.NewRequest("GET", "", strings.NewReader(``))
			r = mux.SetURLVars(r, tt.args.vars)
			r.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fail()
			}

			rh.GetResourceHandler(tt.args.w, r)
			if tt.args.w.Code != tt.args.want {
				t.Errorf("got %d want %d", tt.args.w.Code, tt.args.want)
			}
		})
	}
}

// TestResourceHandler_DeleteResourceHandler DELETE /api/resources/{id}
func TestResourceHandler_DeleteResourceHandler(t *testing.T) {
	type args struct {
		w    *httptest.ResponseRecorder
		body string
		m    func(v interface{}) ([]byte, error)
		u    func(data []byte, v interface{}) error
		db   map[string]map[string]interface{}
		want int
		vars map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "DeleteResource - Success",
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				vars: map[string]string{"id": "0bf8651a-0923-47b8-aed3-e9fc1505e497"},
				want: 204,
			},
		},
		{
			name: "DeleteResource - Invalid UUID Failure",
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				vars: map[string]string{"id": "dummy"},
				want: 400,
			},
		},
		{
			name: "DeleteResource - UUID Not Exist Failure",
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				vars: map[string]string{"id": "0bf8651a-0923-47b8-aed3-e9fc1505e496"},
				want: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup Request
			rh := ResourceHandler{
				ch: &CommonHandler{
					Marshaler:   tt.args.m,
					Unmarshaler: tt.args.u,
				},
				dh: &DBHelper{
					db: tt.args.db,
					mu: sync.Mutex{},
				},
			}

			r, err := http.NewRequest("DELETE", "", strings.NewReader(``))
			r = mux.SetURLVars(r, tt.args.vars)
			r.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fail()
			}

			rh.DeleteResourceHandler(tt.args.w, r)
			if tt.args.w.Code != tt.args.want {
				t.Errorf("got %d want %d", tt.args.w.Code, tt.args.want)
			}
		})
	}
}

// TestResourceHandler_UpdateResourceHandler PUT /api/resources/{id}
func TestResourceHandler_UpdateResourceHandler(t *testing.T) {
	type args struct {
		w    *httptest.ResponseRecorder
		body string
		m    func(v interface{}) ([]byte, error)
		u    func(data []byte, v interface{}) error
		db   map[string]map[string]interface{}
		want int
		vars map[string]string
	}
	tests := []struct {
		reader io.Reader
		name   string
		args   args
	}{
		{
			name:   "UpdateResource - Success",
			reader: strings.NewReader(`{"name":"Bruce","lastname":"Wayne"}`),
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				vars: map[string]string{"id": "0bf8651a-0923-47b8-aed3-e9fc1505e497"},
				want: 202,
			},
		},
		{
			name:   "UpdateResource - Invalid UUID Failure",
			reader: strings.NewReader(`{"name":"Bruce","lastname":"Wayne"}`),
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				vars: map[string]string{"id": "dummy"},
				want: 400,
			},
		},
		{
			name:   "UpdateResource - UUID Not Exist Failure",
			reader: strings.NewReader(`{"name":"Bruce","lastname":"Wayne"}`),
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				vars: map[string]string{"id": "0bf8651a-0923-47b8-aed3-e9fc1505e496"},
				want: 400,
			},
		},
		{
			name:   "UpdateResource - JSON Marshal Failure",
			reader: strings.NewReader(`{"name":"Bruce","lastname":"Wayne"}`),
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				vars: map[string]string{"id": "0bf8651a-0923-47b8-aed3-e9fc1505e497"},
				m:    func(v interface{}) ([]byte, error) { return []byte{}, errors.New("fake error") },
				want: 500,
			},
		},
		{
			name:   "UpdateResource - Invalid Request Failure",
			reader: strings.NewReader(`text`),
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				vars: map[string]string{"id": "0bf8651a-0923-47b8-aed3-e9fc1505e497"},
				want: 400,
			},
		},
		{
			name:   "UpdateResource - Failed ioutil.ReadAll()",
			reader: errReader(0),
			args: args{
				w:    httptest.NewRecorder(),
				db:   map[string]map[string]interface{}{"0bf8651a-0923-47b8-aed3-e9fc1505e497": {"name": "Clark", "lastname": "Kent"}},
				vars: map[string]string{"id": "0bf8651a-0923-47b8-aed3-e9fc1505e497"},
				want: 500,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup Request
			reader := io.Reader(nil)
			if tt.reader != nil {
				reader = tt.reader
			}

			rh := ResourceHandler{
				ch: &CommonHandler{
					Marshaler:   tt.args.m,
					Unmarshaler: tt.args.u,
				},
				dh: &DBHelper{
					db: tt.args.db,
					mu: sync.Mutex{},
				},
			}

			r, err := http.NewRequest("GET", "", reader)
			r = mux.SetURLVars(r, tt.args.vars)
			r.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fail()
			}

			rh.UpdateResourceHandler(tt.args.w, r)
			if tt.args.w.Code != tt.args.want {
				t.Errorf("got %d want %d", tt.args.w.Code, tt.args.want)
			}
		})
	}
}

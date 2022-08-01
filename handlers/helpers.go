package handlers

import "sync"

// DBHelper Definition of main map as well mutex configuration lock/unlock data for multiple concurrent requests.
type DBHelper struct {
	db map[string]map[string]interface{}
	mu sync.Mutex
}

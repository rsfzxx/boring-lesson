// Package dto berisi bentuk data request dan response API.
package dto

type HealthResponse struct {
	Status   string            `json:"status"`   // "ok" atau "degraded"
	Services map[string]string `json:"services"` // status tiap dependency
}

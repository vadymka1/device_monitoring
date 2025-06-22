package http

import (
	"github.com/go-chi/chi/v5"
)

// SetupRoutes configures all HTTP routes for the device API
func SetupRoutes(r chi.Router, handler *DeviceHandler) {
	r.Get("/health/{id}", handler.HealthCheck)
	r.Get("/devices/statuses", handler.GetAllDeviceStatuses)
	r.Route("/devices", func(r chi.Router) {
		r.Get("/{id}", handler.GetDeviceByID)
		r.Get("/{id}/status", handler.GetDeviceStatus)
		r.Get("/{id}/check", handler.CheckDevice)
	})
}

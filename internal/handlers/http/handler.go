package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"device-monitoring/internal/logger"

	"github.com/go-chi/chi/v5"

	"device-monitoring/internal/services"
)

type DeviceHandler struct {
	service *services.DeviceService
	monitor *services.MonitorService
	logger  *logger.Logger
}

func NewDeviceHandler(service *services.DeviceService, monitor *services.MonitorService, logger *logger.Logger) *DeviceHandler {
	return &DeviceHandler{service: service, monitor: monitor, logger: logger}
}

func (h *DeviceHandler) GetDeviceByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.logger.Info("HTTP GET /devices/%s", id)
	device, err := h.service.GetDeviceByID(r.Context(), id)
	if err != nil {
		h.logger.Warn("Device not found %s", id)
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(device)
}

func (h *DeviceHandler) GetDeviceStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.logger.Info("HTTP GET /devices/%s", id)
	status, err := h.service.GetStatusByDeviceID(r.Context(), id)
	if err != nil {
		h.logger.Warn("Status not found for device %s", id)
		http.Error(w, "Status not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(status)
}

func (h *DeviceHandler) CheckDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.logger.Info("HTTP GET /devices/%s", id)
	device, err := h.service.GetDeviceByID(r.Context(), id)
	if err != nil {
		h.logger.Warn("Device not found %s", id)
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	fmt.Println(device.ID)
	status, err := h.monitor.CheckDevice(r.Context(), device)
	if err != nil {
		http.Error(w, "Failed to check device", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(status)
}

func (h *DeviceHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.logger.Info("HTTP GET /health/%s", id)

	if id == "" {
		http.Error(w, "Missing device ID", http.StatusBadRequest)
		return
	}

	status, err := h.service.GetStatusByDeviceID(r.Context(), id)
	if err != nil {
		h.logger.Warn("No status found for device %s", id)
		http.Error(w, "Device status not found", http.StatusNotFound)
		return
	}

	resp := map[string]string{
		"id":     id,
		"status": status.Status,
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *DeviceHandler) GetAllDeviceStatuses(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HTTP GET /devices/statuses")
	devices, err := h.service.GetAllDevices(r.Context())
	if err != nil {
		h.logger.Error("Failed to fetch devices: %v", err)
		http.Error(w, "Failed to fetch devices", http.StatusInternalServerError)
		return
	}
	var results []map[string]string
	for _, d := range devices {
		status, err := h.service.GetStatusByDeviceID(r.Context(), d.ID)
		if err != nil {
			results = append(results, map[string]string{
				"id":     d.ID,
				"status": "unknown",
			})
			continue
		}
		results = append(results, map[string]string{
			"id":     d.ID,
			"status": status.Status,
		})
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"devices": results,
	})
}

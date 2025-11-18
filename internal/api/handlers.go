package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lucast-ruiz/devices-api/internal/service"
)

type Handler struct {
	svc *service.DeviceService
}

func NewHandler(s *service.DeviceService) *Handler {
	return &Handler{svc: s}
}

// CreateDevice godoc
// @Summary Create a new device
// @Description Create a new device with name, brand and state
// @Tags devices
// @Accept json
// @Produce json
// @Param device body model.Device true "Device to create"
// @Success 201 {object} model.Device
// @Failure 400 {string} string "invalid body"
// @Failure 500 {string} string "failed to create device"
// @Router /devices [post]
func (h *Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Brand string `json:"brand"`
		State string `json:"state"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	device, err := h.svc.Create(r.Context(), req.Name, req.Brand, req.State)
	if err != nil {
		http.Error(w, "failed to create device", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(device)
}

// GetDeviceByID godoc
// @Summary Get a device by ID
// @Description Get a single device by its ID
// @Tags devices
// @Produce json
// @Param id path string true "Device ID"
// @Success 200 {object} model.Device
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal error"
// @Router /devices/{id} [get]
func (h *Handler) GetDeviceByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	device, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if device == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(device)
}

// ListDevices godoc
// @Summary List devices
// @Description List all devices or filter by brand/state
// @Tags devices
// @Produce json
// @Param brand query string false "Filter by brand"
// @Param state query string false "Filter by state"
// @Success 200 {array} model.Device
// @Failure 500 {string} string "internal error"
// @Router /devices [get]
func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	brand := r.URL.Query().Get("brand")
	state := r.URL.Query().Get("state")

	if brand != "" {
		devices, err := h.svc.GetByBrand(r.Context(), brand)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(devices)
		return
	}

	if state != "" {
		devices, err := h.svc.GetByState(r.Context(), state)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(devices)
		return
	}

	devices, err := h.svc.ListAll(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(devices)
}

// UpdateDevice godoc
// @Summary Update a device
// @Description Partially update a device
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Param device body model.Device true "Partial device fields"
// @Success 200 {object} model.Device
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /devices/{id} [patch]
func (h *Handler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Name  *string `json:"name"`
		Brand *string `json:"brand"`
		State *string `json:"state"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	device, err := h.svc.Update(r.Context(), id, req.Name, req.Brand, req.State)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if device == nil {
		http.NotFound(w, r)
		return
	}

	_ = json.NewEncoder(w).Encode(device)
}

// DeleteDevice godoc
// @Summary Delete a device
// @Description Delete a device by ID
// @Tags devices
// @Param id path string true "Device ID"
// @Success 204 "no content"
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Router /devices/{id} [delete]
func (h *Handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.svc.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

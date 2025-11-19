package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lucast-ruiz/devices-api/internal/model"
	"github.com/lucast-ruiz/devices-api/internal/service"
)

type Handler struct {
	svc *service.DeviceService
}

func NewHandler(s *service.DeviceService) *Handler {
	return &Handler{svc: s}
}

// CreateDeviceDTO represents the payload to create a device.
type CreateDeviceDTO struct {
	Name  string `json:"name"`
	Brand string `json:"brand"`
	State string `json:"state"`
}

// UpdateDeviceDTO represents the payload to partially update a device.
type UpdateDeviceDTO struct {
	Name  *string `json:"name"`
	Brand *string `json:"brand"`
	State *string `json:"state"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func parseIntQuery(value string, def int) int {
	if value == "" {
		return def
	}
	i, err := strconv.Atoi(value)
	if err != nil || i < 0 {
		return def
	}
	return i
}

// CreateDevice godoc
// @Summary Create a new device
// @Description Create a new device with name, brand and state
// @Tags devices
// @Accept json
// @Produce json
// @Param device body api.CreateDeviceDTO true "Device to create"
// @Success 201 {object} model.Device
// @Failure 400 {object} map[string]string "invalid body or validation error"
// @Failure 500 {object} map[string]string "failed to create device"
// @Router /devices [post]
func (h *Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var req CreateDeviceDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	// Validações básicas
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Brand) == "" {
		writeError(w, http.StatusBadRequest, "name and brand are required")
		return
	}
	if !model.IsValidState(req.State) {
		writeError(w, http.StatusBadRequest, "invalid state value")
		return
	}

	device, err := h.svc.Create(r.Context(), req.Name, req.Brand, req.State)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create device")
		return
	}

	writeJSON(w, http.StatusCreated, device)
}

// GetDeviceByID godoc
// @Summary Get a device by ID
// @Description Get a single device by its ID
// @Tags devices
// @Produce json
// @Param id path string true "Device ID"
// @Success 200 {object} model.Device
// @Failure 404 {object} map[string]string "not found"
// @Failure 500 {object} map[string]string "internal error"
// @Router /devices/{id} [get]
func (h *Handler) GetDeviceByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	device, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if device == nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	writeJSON(w, http.StatusOK, device)
}

// ListDevices godoc
// @Summary List devices
// @Description List all devices or filter by brand/state. If no filter is provided, results are paginated.
// @Tags devices
// @Produce json
// @Param brand query string false "Filter by brand"
// @Param state query string false "Filter by state"
// @Param limit query int false "Max items to return (default 100)"
// @Param offset query int false "Items to skip for pagination (default 0)"
// @Success 200 {array} model.Device
// @Failure 500 {object} map[string]string "internal error"
// @Router /devices [get]
func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	brand := query.Get("brand")
	state := query.Get("state")

	if brand != "" {
		devices, err := h.svc.GetByBrand(r.Context(), brand)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		writeJSON(w, http.StatusOK, devices)
		return
	}

	if state != "" {
		devices, err := h.svc.GetByState(r.Context(), state)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		writeJSON(w, http.StatusOK, devices)
		return
	}

	limit := parseIntQuery(query.Get("limit"), 100)
	offset := parseIntQuery(query.Get("offset"), 0)

	devices, err := h.svc.ListAll(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, devices)
}

// UpdateDevice godoc
// @Summary Update a device
// @Description Partially update a device
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Param device body api.UpdateDeviceDTO true "Partial device fields"
// @Success 200 {object} model.Device
// @Failure 400 {object} map[string]string "validation or business rule error"
// @Failure 404 {object} map[string]string "not found"
// @Failure 500 {object} map[string]string "internal error"
// @Router /devices/{id} [patch]
func (h *Handler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateDeviceDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	// Valida state se vier
	if req.State != nil && !model.IsValidState(*req.State) {
		writeError(w, http.StatusBadRequest, "invalid state value")
		return
	}

	device, err := h.svc.Update(r.Context(), id, req.Name, req.Brand, req.State)
	if err != nil {
		// Erros de regra de negócio
		msg := err.Error()
		if strings.Contains(msg, "cannot ") || strings.Contains(msg, "invalid") {
			writeError(w, http.StatusBadRequest, msg)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if device == nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	writeJSON(w, http.StatusOK, device)
}

// DeleteDevice godoc
// @Summary Delete a device
// @Description Delete a device by ID
// @Tags devices
// @Param id path string true "Device ID"
// @Success 204 "no content"
// @Failure 400 {object} map[string]string "business rule error"
// @Failure 404 {object} map[string]string "not found"
// @Failure 500 {object} map[string]string "internal error"
// @Router /devices/{id} [delete]
func (h *Handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.svc.Delete(r.Context(), id)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "cannot delete") {
			writeError(w, http.StatusBadRequest, msg)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	// Se não existia, retornamos 404 para ser mais explícito
	device, _ := h.svc.GetByID(r.Context(), id)
	if device == nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

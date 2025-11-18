package api

import "github.com/go-chi/chi/v5"


func (h *Handler) Routes() *chi.Mux {
    r := chi.NewRouter()

    r.Post("/devices", h.CreateDevice)
    r.Get("/devices/{id}", h.GetDeviceByID)
    r.Get("/devices", h.ListDevices)
    r.Patch("/devices/{id}", h.UpdateDevice)
    r.Delete("/devices/{id}", h.DeleteDevice)

    return r
}


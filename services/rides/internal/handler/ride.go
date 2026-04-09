package handler

import (
	"copo/rides/internal/model"
	"copo/rides/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RideHandler struct {
	svc *service.RideService
}

func NewRideHandler(svc *service.RideService) *RideHandler {
	return &RideHandler{svc: svc}
}

func (h *RideHandler) Create(w http.ResponseWriter, r *http.Request) {
	driverID := r.Header.Get("X-User-ID")

	var req model.CreateRideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "request invalido", http.StatusBadRequest)
	}

	ride, err := h.svc.Create(r.Context(), driverID, &req)
	if err != nil {
		http.Error(w, "error al crear viaje", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ride)
}

func (h *RideHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rides, err := h.svc.GetAll(r.Context())
	if err != nil {
		http.Error(w, "error al obtener viajes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rides)
}

func (h *RideHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ride, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "viaje no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ride)
}

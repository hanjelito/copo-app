package handler

import (
	"copo/bookings/internal/model"
	"copo/bookings/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type BookingHandler struct {
	svc *service.BookingService
}

func NewBookingHandler(svc *service.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}
func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	var req model.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "request invalid", http.StatusBadRequest)
		return
	}

	booking, err := h.svc.Create(r.Context(), userID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)

}

func (h *BookingHandler) GetMyBookings(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	booking, err := h.svc.GetMyBookings(r.Context(), userID)
	if err != nil {
		http.Error(w, "error retrieving reservations", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}
func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	id := chi.URLParam(r, "id")

	if err := h.svc.Cancel(r.Context(), id, userID); err != nil {
		http.Error(w, "error cancelling reservation", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

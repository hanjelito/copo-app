package handler

import (
	"copo/auth/internal/middleware"
	"copo/auth/internal/model"
	"copo/auth/internal/service"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(middleware.ContextEmail).(string)
	user, err := h.svc.GetMe(r.Context(), email)
	if err != nil {
		http.Error(w, "usuario no encontrado", http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(middleware.ContextEmail).(string)
	var req model.UPdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "request invalido", http.StatusBadRequest)
		return
	}
	user, err := h.svc.UpdateMe(r.Context(), email, &req)
	if err != nil {
		http.Error(w, "error al actualizar", http.StatusInternalServerError)
	}
	w.Header().Set("Context-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

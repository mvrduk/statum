package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	appuser "github.com/mvrduk/statum/services/user-service/internal/application/service/user"
	domain "github.com/mvrduk/statum/services/user-service/internal/domain/user"
	"github.com/mvrduk/statum/services/user-service/internal/transport/http/gen"
	"github.com/mvrduk/statum/services/user-service/internal/transport/http/mapper"
)

var _ gen.ServerInterface = (*UserHandler)(nil)

type UserHandler struct {
	service appuser.Service
}

func NewUserHandler(service appuser.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req gen.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, gen.ErrorResponse{
			Error: "invalid JSON",
		})
		return
	}

	result, err := h.service.Register(r.Context(), appuser.RegisterInput{Email: req.Email, Password: req.Password})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, mapper.ToUserResponse(result))
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request, id string) {
	result, err := h.service.GetByID(r.Context(), id)

	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mapper.ToUserResponse(result))
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidEmail), errors.Is(err, domain.ErrInvalidPassword):
		writeJSON(w, http.StatusBadRequest, gen.ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrEmailAlreadyUsed):
		writeJSON(w, http.StatusConflict, gen.ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrUserNotFound):
		writeJSON(w, http.StatusNotFound, gen.ErrorResponse{Error: err.Error()})
	default:
		log.Println("internal error: ", err)
		writeJSON(w, http.StatusInternalServerError, gen.ErrorResponse{Error: err.Error()})
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("write response: ", err)
	}
}

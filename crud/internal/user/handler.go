package user

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/httpresponse"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	slog.Info("API handler started", "action", "API_HANDLER_START", "method", r.Method, "path", r.URL.Path)
	defer func() {
		slog.Info("API handler completed", "action", "API_HANDLER_END", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds())
	}()

	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		slog.Error("Failed to fetch users", "action", "USERS_FETCH_ERROR", "method", r.Method, "path", r.URL.Path, "error", err)
		httpresponse.Error(w, http.StatusInternalServerError, "USERS_FETCH_FAILED", "Failed to fetch users")
		return
	}

	httpresponse.Success(w, http.StatusOK, "Users fetched successfully", users)
	slog.Info("Users fetched successfully", "action", "USERS_FETCH_SUCCESS", "count", len(users))
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	slog.Info("API handler started", "action", "API_HANDLER_START", "method", r.Method, "path", r.URL.Path)
	defer func() {
		slog.Info("API handler completed", "action", "API_HANDLER_END", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds())
	}()

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var request CreateUserRequest
	if err := decoder.Decode(&request); err != nil {
		slog.Error("Failed to decode create user request", "action", "USER_CREATE_REQUEST_ERROR", "error", err)
		httpresponse.Error(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Request body must contain a valid JSON object")
		return
	}

	var trailingData any
	if err := decoder.Decode(&trailingData); err != io.EOF {
		slog.Error("Create user request contains trailing data", "action", "USER_CREATE_REQUEST_ERROR", "error", err)
		httpresponse.Error(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Request body must contain one JSON object")
		return
	}

	createdUser, err := h.service.CreateUser(r.Context(), request)
	if err != nil {
		slog.Error("Failed to create user", "action", "USER_CREATE_ERROR", "method", r.Method, "path", r.URL.Path, "error", err)
		httpresponse.Error(w, http.StatusInternalServerError, "USER_CREATE_FAILED", "Failed to create user")
		return
	}

	httpresponse.Success(w, http.StatusCreated, "User created successfully", createdUser)
}

package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/joao-vitor-felix/workout-api/internal/store"
	"github.com/joao-vitor-felix/workout-api/internal/utils"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio,omitempty"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore,
		logger,
	}
}

func (h *UserHandler) validateRegisterUserRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if len(req.Email) > 255 {
		return errors.New("email must not exceed 255 characters")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	if len(req.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(req.Username) > 20 {
		return errors.New("username must not exceed 20 characters")
	}
	return nil
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("ERROR: invalid body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	if err := h.validateRegisterUserRequest(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	existing, err := h.userStore.GetByUsername(req.Username)
	if err != nil {
		h.logger.Printf("ERROR: checking existing user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if existing != nil {
		utils.WriteJSON(w, http.StatusConflict, utils.Envelope{"error": "username already taken"})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
		Bio:      req.Bio,
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: setting password hash: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	created, err := h.userStore.Create(user)
	if err != nil {
		h.logger.Printf("ERROR: creating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"data": created})
}

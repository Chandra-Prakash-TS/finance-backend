package handler

import (
	"errors"
	"finance-backend/internal/domain"
	"finance-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input service.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondValidationError(c, []FieldError{{Field: "body", Message: err.Error()}})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateEmail) {
			respondError(c, http.StatusConflict, "CONFLICT", "Email already exists")
			return
		}
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register user")
		return
	}

	respondSuccess(c, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondValidationError(c, []FieldError{{Field: "body", Message: err.Error()}})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid email or password")
			return
		}
		if errors.Is(err, domain.ErrUserInactive) {
			respondError(c, http.StatusForbidden, "FORBIDDEN", "User account is inactive")
			return
		}
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to login")
		return
	}

	respondSuccess(c, http.StatusOK, resp)
}

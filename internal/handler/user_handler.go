package handler

import (
	"errors"
	"finance-backend/internal/domain"
	"finance-backend/internal/service"
	"finance-backend/pkg/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		respondError(c, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}
	respondSuccess(c, http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	var input service.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondValidationError(c, []FieldError{{Field: "body", Message: err.Error()}})
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, input)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update profile")
		return
	}
	respondSuccess(c, http.StatusOK, user)
}

func (h *UserHandler) List(c *gin.Context) {
	var params pagination.Params
	if err := c.ShouldBindQuery(&params); err != nil {
		respondValidationError(c, []FieldError{{Field: "query", Message: err.Error()}})
		return
	}
	params.Defaults()

	users, total, err := h.userService.List(c.Request.Context(), params.Page, params.PageSize)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list users")
		return
	}

	respondPaginated(c, http.StatusOK, users, pagination.NewMeta(params.Page, params.PageSize, total))
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(c, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get user")
		return
	}
	respondSuccess(c, http.StatusOK, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var input service.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondValidationError(c, []FieldError{{Field: "body", Message: err.Error()}})
		return
	}

	user, err := h.userService.Update(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(c, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update user")
		return
	}
	respondSuccess(c, http.StatusOK, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(c, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete user")
		return
	}
	respondSuccess(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}

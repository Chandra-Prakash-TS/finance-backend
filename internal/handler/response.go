package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination,omitempty"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func respondSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, SuccessResponse{Data: data})
}

func respondPaginated(c *gin.Context, status int, data interface{}, pagination interface{}) {
	c.JSON(status, SuccessResponse{Data: data, Pagination: pagination})
}

func respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func respondValidationError(c *gin.Context, details []FieldError) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: ErrorBody{
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Details: details,
		},
	})
}

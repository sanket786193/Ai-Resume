package response

import (
	"net/http"

	domainerrors "resume/internal/domain/errors"

	"github.com/gin-gonic/gin"
)

// Standard API response envelope.
type Body struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo for API errors.
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// JSON sends a success JSON response.
func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Body{Success: true, Data: data})
}

// Error maps domain errors to HTTP status and sends error body.
func Error(c *gin.Context, err error) {
	code, message := "INTERNAL_ERROR", err.Error()
	status := http.StatusInternalServerError

	switch e := err.(type) {
	case *domainerrors.NotFoundError:
		status = http.StatusNotFound
		code = "NOT_FOUND"
		message = e.Error()
	case *domainerrors.ValidationError:
		status = http.StatusBadRequest
		code = "VALIDATION_ERROR"
		message = e.Error()
	case *domainerrors.ConflictError:
		status = http.StatusConflict
		code = "CONFLICT"
		message = e.Message
	case *domainerrors.UnauthorizedError:
		status = http.StatusUnauthorized
		code = "UNAUTHORIZED"
		message = e.Message
	case *domainerrors.ForbiddenError:
		status = http.StatusForbidden
		code = "FORBIDDEN"
		message = e.Message
	}

	c.JSON(status, Body{
		Success: false,
		Error:   &ErrorInfo{Code: code, Message: message},
	})
}

// ValidationError sends 400 for binding/validation errors.
func ValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Body{
		Success: false,
		Error:   &ErrorInfo{Code: "VALIDATION_ERROR", Message: err.Error()},
	})
}

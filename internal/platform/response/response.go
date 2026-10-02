package response

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourname/gin-template/internal/platform/apperror"
)

type Body struct {
	Success bool            `json:"success"`
	Data    any             `json:"data,omitempty"`
	Meta    any             `json:"meta,omitempty"`
	Error   *apperror.Error `json:"error,omitempty"`
}

type PageMeta struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Success: true, Data: data})
}

func OKWithMeta(c *gin.Context, data, meta any) {
	c.JSON(http.StatusOK, Body{Success: true, Data: data, Meta: meta})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Body{Success: true, Data: data})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Fail converts any error to a JSON response. Unknown errors become 500s.
func Fail(c *gin.Context, err error) {
	var appErr *apperror.Error
	if !errors.As(err, &appErr) {
		appErr = apperror.Internal(err)
	}
	if appErr.Status >= http.StatusInternalServerError {
		slog.ErrorContext(c.Request.Context(), "request failed",
			"error", err, "request_id", c.GetString("request_id"))
	}
	c.AbortWithStatusJSON(appErr.Status, Body{Success: false, Error: appErr})
}

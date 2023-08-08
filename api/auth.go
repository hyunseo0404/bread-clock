package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type authHandler struct {
}

type loginRequest struct {
	Code     string `json:"code"`
	Provider string `json:"provider"`
}

func (h *authHandler) login(c *gin.Context) {
	// TODO

	c.Status(http.StatusOK)
}

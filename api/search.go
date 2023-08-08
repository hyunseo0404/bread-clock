package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type searchHandler struct {
}

func (h *searchHandler) searchBakeries(c *gin.Context) {
	// TODO

	c.Status(http.StatusOK)
}

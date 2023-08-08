package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type bakeriesHandler struct {
}

type openingHours struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}

type bread struct {
	ID        int  `json:"id"`
	Available bool `json:"available"`
}

type breadDetail struct {
	bread
	AvailableHours []string `json:"available_hours"`
	PhotoURL       string   `json:"photo_url"`
}

type breadList struct {
	Breads []bread `json:"breads"`
}

type bakery struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Coordinates  string        `json:"coordinates"`
	Favorite     bool          `json:"favorite"`
	BreadDetails []breadDetail `json:"breads"`
	PhotoURLs    []string      `json:"photo_urls"`
}

type bakeryDetail struct {
	bakery
	Address      string         `json:"address"`
	OpeningHours []openingHours `json:"opening_hours"`
}

type bakeryList struct {
	Bakeries []*bakery `json:"bakeries"`
}

func (h *bakeriesHandler) listBakeries(c *gin.Context) {
	// TODO

	c.JSON(http.StatusOK, &bakeryList{})
}

func (h *bakeriesHandler) getBakery(c *gin.Context) {
	// TODO

	c.JSON(http.StatusOK, &bakeryDetail{})
}

func (h *bakeriesHandler) markBakeryAsFavorite(c *gin.Context) {
	// TODO

	c.Status(http.StatusOK)
}

func (h *bakeriesHandler) unmarkBakeryAsFavorite(c *gin.Context) {
	// TODO

	c.Status(http.StatusOK)
}

func (h *bakeriesHandler) updateBreadAvailabilities(c *gin.Context) {
	// TODO

	c.Status(http.StatusOK)
}

package api

import (
	"bread-clock/db"
	"bread-clock/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type bakeriesHandler struct {
	bakeryRepository db.BakeryRepository
}

type OpeningHours struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}

type Bread struct {
	ID        int  `json:"id"`
	Available bool `json:"available"`
}

type BreadDetail struct {
	Bread
	AvailableHours []string `json:"availableHours"`
	PhotoURL       string   `json:"photoUrl"`
}

type BreadList struct {
	Breads []Bread `json:"breads"`
}

type Bakery struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Coordinates  string        `json:"coordinates"`
	Favorite     bool          `json:"favorite"`
	BreadDetails []BreadDetail `json:"breads"`
	PhotoURLs    []string      `json:"photoUrls"`
}

type BakeryDetail struct {
	Bakery
	Address      string         `json:"address"`
	OpeningHours []OpeningHours `json:"openingHours"`
}

type BakeryList struct {
	Bakeries []Bakery `json:"bakeries"`
}

type listBakeriesRequest struct {
	Sort     string `form:"sort"`
	Size     int    `form:"size"`
	Offset   int    `form:"offset"`
	Filter   string `form:"filter"`
	Location string `form:"loc"`
}

type updateBreadAvailability struct {
	ID        int  `json:"id"`
	Available bool `json:"available"`
}

type updateBreadAvailabilitiesRequest struct {
	BakeryID            int                       `form:"bakeryId"`
	BreadAvailabilities []updateBreadAvailability `json:"breads"`
}

// listBakeries godoc
// @Summary		빵집 목록 조회
// @Description 전체 빵집 목록 조회 (sort param 미지정 시 이름순 정렬)
// @Tags		Bakeries
// @Produce		json
// @Param		sort query string false "정렬 옵션 (name|distance)"
// @Param		size query string false "조회 개수"
// @Param		offset query string false "조회 offset"
// @Param		filter query string false "필터 옵션 (favorites)"
// @Param		loc query string false "현재 위치 좌표값 (위도,경도)"
// @Success		200 {object} BakeryList
// @Failure		400
// @Failure		500
// @Router		/bakeries [GET]
func (h *bakeriesHandler) listBakeries(c *gin.Context) {
	userID := 0 // FIXME: get current user ID

	var req listBakeriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	var sortOption db.SortOption
	switch req.Sort {
	case "name":
		sortOption = db.SortByName
	case "distance":
		sortOption = db.SortByDistance
	default:
		sortOption = db.SortByName
	}

	if req.Size < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid size value: %d", req.Size)})
		return
	}

	if req.Offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid offset value: %d", req.Offset)})
		return
	}

	if req.Size == 0 {
		req.Size = 10
	}

	var latitude, longitude float64
	var needsDistance bool
	if sortOption == db.SortByDistance {
		if len(req.Location) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint("coordinates required if using the 'distance' filter option")})
			return
		}

		_, err := fmt.Sscanf(req.Location+"~", "%f,%f~", &latitude, &longitude)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid coordinates: %s", req.Location)})
			return
		}

		needsDistance = true
	}

	bakeries, err := h.bakeryRepository.List(c, sortOption, req.Size, req.Offset, latitude, longitude, needsDistance, userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, bakeries)
}

// getBakery godoc
// @Summary		빵집 상세 조회
// @Description 특정 빵집에 대한 상세 정보 조회
// @Tags		Bakeries
// @Produce		json
// @Param		bakeryId path int true "빵집 ID"
// @Param		loc query string false "현재 위치 좌표값 (위도,경도)"
// @Success		200 {object} BakeryDetail
// @Failure		400
// @Failure		404
// @Failure		500
// @Router		/bakeries/:bakeryId [GET]
func (h *bakeriesHandler) getBakery(c *gin.Context) {
	userID := 0 // FIXME: get current user ID

	bakeryID, err := strconv.Atoi(c.Param("bakeryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var latitude, longitude float64
	location, needsDistance := c.GetQuery("loc")
	if needsDistance {
		_, err := fmt.Sscanf(location+"~", "%f,%f~", &latitude, &longitude)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid coordinates: %s", location)})
			return
		}
	}

	bakery, err := h.bakeryRepository.Get(c, bakeryID, latitude, longitude, needsDistance, userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &bakery)
}

// markBakeryAsFavorite godoc
// @Summary		빵집 즐겨찾기 추가
// @Description 특정 빵집을 즐겨찾기에 추가
// @Tags		Bakeries
// @Produce		json
// @Param		bakeryId path int true "빵집 ID"
// @Success		200
// @Failure		400
// @Failure		404
// @Failure		500
// @Router		/bakeries/:bakeryId/favorite [PUT]
func (h *bakeriesHandler) markBakeryAsFavorite(c *gin.Context) {
	userID := 0 // FIXME: get current user ID

	bakeryID, err := strconv.Atoi(c.Param("bakeryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err = h.bakeryRepository.MarkAsFavorite(c, bakeryID, userID); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// unmarkBakeryAsFavorite godoc
// @Summary		빵집 즐겨찾기 해제
// @Description 특정 빵집을 즐겨찾기에서 해제
// @Tags		Bakeries
// @Produce		json
// @Param		bakeryId path int true "빵집 ID"
// @Success		200
// @Failure		400
// @Failure		404
// @Failure		500
// @Router		/bakeries/:bakeryId/favorite [DELETE]
func (h *bakeriesHandler) unmarkBakeryAsFavorite(c *gin.Context) {
	userID := 0 // FIXME: get current user ID

	bakeryID, err := strconv.Atoi(c.Param("bakeryId"))

	if err = h.bakeryRepository.UnmarkAsFavorite(c, bakeryID, userID); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// updateBreadAvailabilities godoc
// @Summary		빵 매진 정보 갱신
// @Description 특정 빵집의 빵 종류별 매진 정보 갱신
// @Tags		Bakeries
// @Produce		json
// @Param		bakeryId path int true "빵집 ID"
// @Param		BreadList body BreadList true "빵 정보 리스트"
// @Success		200
// @Failure		400
// @Failure		404
// @Failure		500
// @Router		/bakeries/:bakeryId/breads/availability [PUT]
func (h *bakeriesHandler) updateBreadAvailabilities(c *gin.Context) {
	var req updateBreadAvailabilitiesRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var err error
	req.BakeryID, err = strconv.Atoi(c.Param("bakeryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var breadAvailabilities []models.BreadAvailability
	for _, availability := range req.BreadAvailabilities {
		breadAvailabilities = append(breadAvailabilities, models.BreadAvailability{
			BakeryID:  req.BakeryID,
			BreadID:   availability.ID,
			Available: availability.Available,
		})
	}

	if err = h.bakeryRepository.UpdateBreadAvailabilities(c, breadAvailabilities); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

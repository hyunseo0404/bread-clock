package api

import (
	"bread-clock/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type searchHandler struct {
	bakeryRepository db.BakeryRepository
}

type searchBakeriesRequest struct {
	Query    string `form:"q"`
	Sort     string `form:"sort"`
	Size     int    `form:"size"`
	Offset   int    `form:"offset"`
	Filter   string `form:"filter"`
	Location string `form:"loc"`
}

// searchBakeries godoc
// @Summary		빵 이름으로 빵집 검색
// @Description 특정 빵 종류가 제공되는 빵집들을 검색 (sort param 미지정 시 이름순 정렬)
// @Tags		Search
// @Produce		json
// @Param		q query string true "검색 문자열"
// @Param		sort query string false "정렬 옵션 (name|distance)"
// @Param		size query string false "조회 개수"
// @Param		offset query string false "조회 offset"
// @Param		filter query string false "필터 옵션 (favorites)"
// @Param		loc query string false "현재 위치 좌표값 (위도,경도)"
// @Success		200 {object} BakeryList
// @Failure		400
// @Failure		500
// @Router		/search [GET]
func (h *searchHandler) searchBakeries(c *gin.Context) {
	userID := 0 // FIXME: get current user ID

	var req searchBakeriesRequest
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
	if req.Sort == db.SortByDistance {
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

	bakeries, err := h.bakeryRepository.ListForBreads(c, req.Query, sortOption, req.Size, req.Offset, latitude, longitude, needsDistance, userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, bakeries)
}

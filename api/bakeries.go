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

// listBakeries godoc
// @Summary		빵집 목록 조회
// @Description 전체 빵집 목록 조회 (sort param 미지정 시 이름순 정렬)
// @Tags		Bakeries
// @Produce		json
// @Param		sort query string false "정렬 옵션 (name|distance)"
// @Param		loc query string false "현재 위치 좌표값 (위도,경도)"
// @Param		filter query string false "필터 옵션 (favorites)"
// @Success		200 {object} bakeryList
// @Failure		400
// @Failure		500
// @Router		/bakeries [GET]
func (h *bakeriesHandler) listBakeries(c *gin.Context) {
	// TODO

	c.JSON(http.StatusOK, &bakeryList{})
}

// getBakery godoc
// @Summary		빵집 상세 조회
// @Description 특정 빵집에 대한 상세 정보 조회
// @Tags		Bakeries
// @Produce		json
// @Param		bakeryId path int true "빵집 ID"
// @Success		200 {object} bakeryDetail
// @Failure		400
// @Failure		404
// @Failure		500
// @Router		/bakeries/:bakeryId [GET]
func (h *bakeriesHandler) getBakery(c *gin.Context) {
	// TODO

	c.JSON(http.StatusOK, &bakeryDetail{})
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
	// TODO

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
	// TODO

	c.Status(http.StatusOK)
}

// updateBreadAvailabilities godoc
// @Summary		빵 매진 정보 갱신
// @Description 특정 빵집의 빵 종류별 매진 정보 갱신
// @Tags		Bakeries
// @Produce		json
// @Param		bakeryId path int true "빵집 ID"
// @Param		breadList body breadList true "빵 정보 리스트"
// @Success		200
// @Failure		400
// @Failure		404
// @Failure		500
// @Router		/bakeries/:bakeryId/breads/availability [PUT]
func (h *bakeriesHandler) updateBreadAvailabilities(c *gin.Context) {
	// TODO

	c.Status(http.StatusOK)
}

package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type searchHandler struct {
}

// searchBakeries godoc
// @Summary		빵 이름으로 빵집 검색
// @Description 특정 빵 종류가 제공되는 빵집들을 검색 (sort param 미지정 시 이름순 정렬)
// @Tags		Search
// @Produce		json
// @Param		q query string true "검색 문자열"
// @Param		sort query string false "정렬 옵션 (name|distance)"
// @Param		loc query string false "현재 위치 좌표값 (위도,경도)"
// @Success		200 {object} bakeryList
// @Failure		400
// @Failure		500
// @Router		/search [GET]
func (h *searchHandler) searchBakeries(c *gin.Context) {
	// TODO

	c.Status(http.StatusOK)
}

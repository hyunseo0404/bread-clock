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

// login godoc
// @Summary		로그인
// @Description OAuth2로 발급 받은 코드를 이용한 토큰 발급 및 로그인 처리
// @Tags		Authentication
// @Produce		json
// @Param		loginRequest body loginRequest true "로그인 요청 정보"
// @Success		200
// @Failure		400
// @Failure		500
// @Router		/auth/login [POST]
func (h *authHandler) login(c *gin.Context) {
	// TODO

	c.Status(http.StatusOK)
}

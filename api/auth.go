package api

import (
	"bread-clock/db"
	e "bread-clock/error"
	"bread-clock/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authHandler struct {
	userRepository db.UserRepository
}

type loginRequest struct {
	AccessToken string `json:"accessToken"`
	Provider    string `form:"provider"`
}

type loginResponse struct {
	AccessToken string `json:"accessToken"`
	Provider    string `json:"provider"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatarUrl"`
}

// login godoc
// @Summary		로그인
// @Description OAuth2로 발급 받은 코드를 이용한 토큰 발급 및 로그인 처리
// @Tags		Authentication
// @Produce		json
// @Param		provider path string true "OAuth2 provider"
// @Param		loginRequest body loginRequest true "로그인 요청 정보"
// @Success		200
// @Failure		400
// @Failure		500
// @Router		/auth/login/:provider [POST]
func (h *authHandler) login(c *gin.Context) {
	req := loginRequest{
		Provider: c.Param("provider"),
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userInfo, err := util.GetUserInfo(c, req.Provider, req.AccessToken)
	if err != nil {
		switch {
		case errors.Is(err, e.ErrAuthInvalidToken):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		case errors.Is(err, e.ErrAuthInvalidProvider):
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid provider: %s", req.Provider)})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	user, err := h.userRepository.FindOrCreate(c, userInfo.ID, req.Provider, userInfo.Email, userInfo.Picture)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	accessToken, err := util.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &loginResponse{
		AccessToken: accessToken,
		Provider:    user.Provider,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
	})
}

// getUser godoc
// @Summary		사용자 정보 조회
// @Description 현재 로그인된 사용자 정보 조회
// @Tags		Authentication
// @Produce		json
// @Success		200 {object} models.User
// @Failure		401
// @Failure		404
// @Failure		500
// @Router		/auth/me [GET]
func (h *authHandler) getUser(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not logged in"})
		return
	}

	user, err := h.userRepository.FindByUserID(c, userID)
	if err != nil {
		if errors.Is(err, e.ErrDBNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &user)
}

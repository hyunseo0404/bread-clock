package middlewares

import (
	"bread-clock/util"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := util.GetToken(c.Query("token"), c.Request.Header.Get("Authorization"))
		if err == nil {
			token, err := util.ValidateToken(tokenString)
			if err == nil {
				if userID, err := util.ExtractUserID(token); err == nil && userID > 0 {
					c.Set("user_id", userID)
				}
			}
		}
		c.Next()
	}
}

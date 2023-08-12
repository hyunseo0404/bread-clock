package api

import (
	"bread-clock/db"
	"bread-clock/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, userRepository db.UserRepository, bakeryRepository db.BakeryRepository) {
	g := r.Group("/api/v1")

	ah := authHandler{
		userRepository: userRepository,
	}
	authRouter := g.Group("/auth")
	authRouter.POST("/login/:provider", ah.login)

	g.Use(middlewares.AuthMiddleware())

	sh := searchHandler{
		bakeryRepository: bakeryRepository,
	}
	searchRouter := g.Group("/search")
	searchRouter.GET("", sh.searchBakeries)

	bh := bakeriesHandler{
		bakeryRepository: bakeryRepository,
	}
	bakeriesRouter := g.Group("/bakeries")
	bakeriesRouter.GET("", bh.listBakeries)
	bakeriesRouter.GET("/:bakeryId", bh.getBakery)
	bakeriesRouter.PUT("/:bakeryId/favorite", bh.markBakeryAsFavorite)
	bakeriesRouter.DELETE("/:bakeryId/favorite", bh.unmarkBakeryAsFavorite)
	bakeriesRouter.PUT("/:bakeryId/breads/availability", bh.updateBreadAvailabilities)
}

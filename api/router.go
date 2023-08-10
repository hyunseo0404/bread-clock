package api

import (
	"bread-clock/db"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, bakeryRepository db.BakeryRepository) {
	g := r.Group("/api/v1")

	ah := authHandler{}
	authRouter := g.Group("/auth")
	authRouter.POST("/login", ah.login)

	sh := searchHandler{
		bakeryRepository: bakeryRepository,
	}
	searchRouter := g.Group("/search")
	searchRouter.GET("/", sh.searchBakeries)

	bh := bakeriesHandler{
		bakeryRepository: bakeryRepository,
	}
	bakeriesRouter := g.Group("/bakeries")
	bakeriesRouter.GET("/", bh.listBakeries)
	bakeriesRouter.GET("/:bakeryId", bh.getBakery)
	bakeriesRouter.PUT("/:bakeryId/favorite", bh.markBakeryAsFavorite)
	bakeriesRouter.DELETE("/:bakeryId/favorite", bh.unmarkBakeryAsFavorite)
	bakeriesRouter.PUT("/:bakeryId/breads/availability", bh.updateBreadAvailabilities)
}

package v1

import (
	"github.com/Homyakadze14/AFFARM_tz/internal/usecase"

	"github.com/gin-gonic/gin"
)

type cryptocurrencyRoutes struct {
	h *usecase.CryptocurrencyService
}

func NewHellotRoutes(handler *gin.RouterGroup, h *usecase.CryptocurrencyService) {
	r := &cryptocurrencyRoutes{h}

	g := handler.Group("")
	{
		g.GET("/hello", r.get)
	}
}

// @Summary     Test
// @Description Test
// @ID          Test
// @Tags  	    Test
// @Produce     json
// @Success     200
// @Router      /hello [get]
func (r *cryptocurrencyRoutes) get(c *gin.Context) {
}

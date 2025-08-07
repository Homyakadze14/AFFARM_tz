package v1

import (
	"log/slog"
	"net/http"

	"github.com/Homyakadze14/AFFARM_tz/internal/common"
	"github.com/Homyakadze14/AFFARM_tz/internal/dto"
	"github.com/Homyakadze14/AFFARM_tz/internal/entity"
	"github.com/Homyakadze14/AFFARM_tz/internal/usecase"

	"github.com/gin-gonic/gin"
)

type cryptocurrencyRoutes struct {
	log *slog.Logger
	h   *usecase.CryptocurrencyService
}

func NewHellotRoutes(log *slog.Logger, handler *gin.RouterGroup, h *usecase.CryptocurrencyService) {
	r := &cryptocurrencyRoutes{log, h}

	g := handler.Group("currency")
	{
		g.POST("/add", r.add)
		g.POST("/remove", r.remove)
		g.POST("/price", r.price)
	}
}

func handlErr(c *gin.Context, log *slog.Logger, err error) {
	log.Error(err.Error())
	status, err := common.ParseErr(err)
	c.JSON(status, gin.H{"error": err.Error()})
}

// @Summary     Add cryptocurrency
// @Description Add cryptocurrency
// @ID          AddCryptocurrency
// @Tags  	    AddCryptocurrency
// @Accept      json
// @Param 		cryptocurrency body dto.AddCryptocurrencyRequest false "Cryptocurrency add data"
// @Success     200
// @Failure     400
// @Failure     404
// @Failure     500
// @Router      /currency/add [post]
func (r *cryptocurrencyRoutes) add(c *gin.Context) {
	const op = "cryptocurrencyRoutes.add"
	log := r.log.With(
		slog.String("op", op),
	)

	var req *dto.AddCryptocurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlErr(c, log, err)
		return
	}

	cr := &entity.Cryptocurrency{
		Symbol: req.Symbol,
	}
	err := r.h.Add(c.Request.Context(), cr)
	if err != nil {
		handlErr(c, log, err)
		return
	}

	c.JSON(http.StatusOK, "")
}

// @Summary     Remove cryptocurrency
// @Description Remove cryptocurrency
// @ID          RemoveCryptocurrency
// @Tags  	    RemoveCryptocurrency
// @Accept      json
// @Param 		cryptocurrency body dto.RemoveCryptocurrencyRequest false "Cryptocurrency remove data"
// @Success     200
// @Failure     400
// @Failure     404
// @Failure     500
// @Router      /currency/remove [post]
func (r *cryptocurrencyRoutes) remove(c *gin.Context) {
	const op = "cryptocurrencyRoutes.remove"
	log := r.log.With(
		slog.String("op", op),
	)

	var req *dto.RemoveCryptocurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlErr(c, log, err)
		return
	}

	cr := &entity.Cryptocurrency{
		Symbol: req.Symbol,
	}
	err := r.h.Remove(c.Request.Context(), cr)
	if err != nil {
		handlErr(c, log, err)
		return
	}

	c.JSON(http.StatusOK, "")
}

// @Summary     Get price
// @Description Get pice
// @ID          GetPriceCryptocurrency
// @Tags  	    GetPriceCryptocurrency
// @Accept      json
// @Param 		price body dto.PriceRequest false "Get price data"
// @Produce     json
// @Success     200 {object} dto.PriceResponse
// @Failure     400
// @Failure     404
// @Failure     500
// @Router      /currency/price [post]
func (r *cryptocurrencyRoutes) price(c *gin.Context) {
	const op = "cryptocurrencyRoutes.price"
	log := r.log.With(
		slog.String("op", op),
	)

	var req *dto.PriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlErr(c, log, err)
		return
	}

	hist, err := r.h.Price(c.Request.Context(), req.Symbol, req.Timestamp)
	if err != nil {
		handlErr(c, log, err)
		return
	}

	resp := &dto.PriceResponse{
		Price: hist.Price,
	}
	c.JSON(http.StatusOK, resp)
}

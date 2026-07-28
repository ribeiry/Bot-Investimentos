package handler

import (
	"net/http"
	"portifolio-api/internal/usecase/market"

	"github.com/gin-gonic/gin"
)

type MarketHandler struct {
	getClosed market.GetCloseUseCase
	getPrice  market.GetPricesUseCase
}

func NewMarketHandler(getClose market.GetCloseUseCase, getPrices market.GetPricesUseCase) *MarketHandler {
	return &MarketHandler{
		getClosed: getClose,
		getPrice:  getPrices,
	}
}

func (m MarketHandler) GetCloseMarket(c *gin.Context) {
	userID := c.GetInt64("userID")
	mkt := c.Query("market")
	quotes, err := m.getClosed.Execute(userID, mkt)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusOK, quotes)
}

func (m MarketHandler) GetPriceMarket(c *gin.Context) {
	userID := c.GetInt64("userID")
	quotes, err := m.getPrice.Execute(userID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusOK, quotes)
}

package handler

import (
	"net/http"
	"portifolio-api/internal/usecase/market"

	"github.com/gin-gonic/gin"
)

type MarketHandler struct {
	// quatro usecases aqui
	getClosed market.GetCloseUseCase
	getPrice  market.GetPricesUseCase
}

func NewMarketHandler(getClose market.GetCloseUseCase,
	getPrices market.GetPricesUseCase,
) *MarketHandler {

	return &MarketHandler{
		getClosed: getClose,
		getPrice:  getPrices,
	}
}

func (m MarketHandler) GetCloseMarket(c *gin.Context) {
	market := c.Query("market")
	quotes, error := m.getClosed.Execute(market)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, quotes)
}

func (m MarketHandler) GetPriceMarket(c *gin.Context) {
	quotesPrices, error := m.getPrice.Execute()
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, quotesPrices)
}

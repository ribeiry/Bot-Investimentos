package handler

import (
	"net/http"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/usecase/portifolio"

	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
	// quatro usecases aqui
	upsert         portifolio.UpsertAssetUseCase
	delete         portifolio.DeleteAssetUseCase
	getAsset       portifolio.GetAssetUseCase
	getSummary     portifolio.GetSummaryUseCase
	getPerformance portifolio.GetPerformanceUseCase
}

func NewPortfolioHandler(upsert portifolio.UpsertAssetUseCase,
	delete portifolio.DeleteAssetUseCase,
	getAsset portifolio.GetAssetUseCase,
	getSummary portifolio.GetSummaryUseCase,
	getPerformance portifolio.GetPerformanceUseCase,
) *PortfolioHandler {

	return &PortfolioHandler{
		upsert:         upsert,
		delete:         delete,
		getAsset:       getAsset,
		getSummary:     getSummary,
		getPerformance: getPerformance}
}

func (h PortfolioHandler) GetAssets(c *gin.Context) {
	assets, error := h.getAsset.Execute()
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, assets)
}

func (h PortfolioHandler) GetSummaryAsset(c *gin.Context) {
	mode := c.Query("mode")
	assets, error := h.getSummary.Execute(mode)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, assets)
}

func (h PortfolioHandler) UpsertAsset(c *gin.Context) {
	var asset domain.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	error := h.upsert.Execute(asset)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ativo salvo com sucesso"})
}

func (h PortfolioHandler) DeleteAsset(c *gin.Context) {
	ticker := c.Param("ticker")
	error := h.delete.Execute(ticker)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ativo deletado com sucesso"})
}

func (h PortfolioHandler) GetPerformance(c *gin.Context) {
	performance, error := h.getPerformance.Execute()

	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}

	c.JSON(http.StatusOK, performance)
}

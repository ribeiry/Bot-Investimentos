package handler

import (
	"net/http"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/usecase/portifolio"

	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
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
		getPerformance: getPerformance,
	}
}

func (h PortfolioHandler) GetAssets(c *gin.Context) {
	userID := c.GetInt64("userID")
	assets, err := h.getAsset.Execute(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assets)
}

func (h PortfolioHandler) GetSummaryAsset(c *gin.Context) {
	userID := c.GetInt64("userID")
	mode := c.Query("mode")
	assets, err := h.getSummary.Execute(userID, mode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assets)
}

func (h PortfolioHandler) UpsertAsset(c *gin.Context) {
	userID := c.GetInt64("userID")
	var asset domain.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.upsert.Execute(userID, asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ativo salvo com sucesso"})
}

func (h PortfolioHandler) DeleteAsset(c *gin.Context) {
	userID := c.GetInt64("userID")
	ticker := c.Param("ticker")
	if err := h.delete.Execute(userID, ticker); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ativo deletado com sucesso"})
}

func (h PortfolioHandler) GetPerformance(c *gin.Context) {
	userID := c.GetInt64("userID")
	performance, err := h.getPerformance.Execute(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, performance)
}

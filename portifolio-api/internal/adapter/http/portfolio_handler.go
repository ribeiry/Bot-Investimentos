package handler

import (
	"net/http"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/usecase/portifolio"

	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
	upsert           portifolio.UpsertAssetUseCase
	delete           portifolio.DeleteAssetUseCase
	getAsset         portifolio.GetAssetUseCase
	getSummary       portifolio.GetSummaryUseCase
	getPerformance   portifolio.GetPerformanceUseCase
	getPeriodSummary portifolio.GetPeriodSummaryUseCase
	getBenchmark     portifolio.GetBenchmarkUseCase
	getAllocation    portifolio.GetAllocationUseCase
	updateSector     portifolio.UpdateSectorUseCase
}

func NewPortfolioHandler(
	upsert portifolio.UpsertAssetUseCase,
	delete portifolio.DeleteAssetUseCase,
	getAsset portifolio.GetAssetUseCase,
	getSummary portifolio.GetSummaryUseCase,
	getPerformance portifolio.GetPerformanceUseCase,
	getPeriodSummary portifolio.GetPeriodSummaryUseCase,
	getBenchmark portifolio.GetBenchmarkUseCase,
	getAllocation portifolio.GetAllocationUseCase,
	updateSector portifolio.UpdateSectorUseCase,
) *PortfolioHandler {
	return &PortfolioHandler{
		upsert:           upsert,
		delete:           delete,
		getAsset:         getAsset,
		getSummary:       getSummary,
		getPerformance:   getPerformance,
		getPeriodSummary: getPeriodSummary,
		getBenchmark:     getBenchmark,
		getAllocation:    getAllocation,
		updateSector:     updateSector,
	}
}

func (h PortfolioHandler) GetAssets(c *gin.Context) {
	userID := c.GetInt64("userID")
	assets, err := h.getAsset.Execute(userID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusOK, assets)
}

func (h PortfolioHandler) GetSummaryAsset(c *gin.Context) {
	userID := c.GetInt64("userID")
	mode := c.Query("mode")
	groupBy := c.DefaultQuery("group_by", "market")
	summary, err := h.getSummary.Execute(userID, mode, groupBy)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusOK, summary)
}

func (h PortfolioHandler) UpsertAsset(c *gin.Context) {
	userID := c.GetInt64("userID")
	var asset domain.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}
	if err := h.upsert.Execute(userID, asset); err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondMessage(c, http.StatusOK, "ativo salvo com sucesso")
}

func (h PortfolioHandler) DeleteAsset(c *gin.Context) {
	userID := c.GetInt64("userID")
	ticker := c.Param("ticker")
	if err := h.delete.Execute(userID, ticker); err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondMessage(c, http.StatusOK, "ativo deletado com sucesso")
}

func (h PortfolioHandler) GetPerformance(c *gin.Context) {
	userID := c.GetInt64("userID")
	performance, err := h.getPerformance.Execute(userID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusOK, performance)
}

func (h PortfolioHandler) UpdateSector(c *gin.Context) {
	userID := c.GetInt64("userID")
	ticker := c.Param("ticker")
	var input portifolio.UpdateSectorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}
	if err := h.updateSector.Execute(userID, ticker, input); err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondMessage(c, http.StatusOK, "setor atualizado com sucesso")
}

func (h PortfolioHandler) GetAllocation(c *gin.Context) {
	userID := c.GetInt64("userID")
	result, err := h.getAllocation.Execute(userID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusOK, result)
}

func (h PortfolioHandler) GetBenchmark(c *gin.Context) {
	userID := c.GetInt64("userID")
	period := c.Query("period")
	result, err := h.getBenchmark.Execute(userID, period)
	if err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}
	respond(c, http.StatusOK, result)
}

func (h PortfolioHandler) GetPeriodSummary(c *gin.Context) {
	userID := c.GetInt64("userID")
	period := c.Query("period")
	summary, err := h.getPeriodSummary.Execute(userID, period)
	if err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}
	respond(c, http.StatusOK, summary)
}

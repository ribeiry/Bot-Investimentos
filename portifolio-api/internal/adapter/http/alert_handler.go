package handler

import (
	"net/http"
	usecasealert "portifolio-api/internal/usecase/alert"

	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	upsert  usecasealert.UpsertAlertUseCase
	delete  usecasealert.DeleteAlertUseCase
	get     usecasealert.GetAlertsUseCase
	check   usecasealert.CheckAlertsUseCase
}

func NewAlertHandler(
	upsert usecasealert.UpsertAlertUseCase,
	delete usecasealert.DeleteAlertUseCase,
	get usecasealert.GetAlertsUseCase,
	check usecasealert.CheckAlertsUseCase,
) *AlertHandler {
	return &AlertHandler{upsert: upsert, delete: delete, get: get, check: check}
}

func (h AlertHandler) UpsertAlert(c *gin.Context) {
	userID := c.GetInt64("userID")
	var input usecasealert.UpsertAlertInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.upsert.Execute(userID, input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "alerta salvo com sucesso"})
}

func (h AlertHandler) DeleteAlert(c *gin.Context) {
	userID := c.GetInt64("userID")
	ticker := c.Param("ticker")
	if err := h.delete.Execute(userID, ticker); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "alerta removido com sucesso"})
}

func (h AlertHandler) GetAlerts(c *gin.Context) {
	userID := c.GetInt64("userID")
	alerts, err := h.get.Execute(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alerts)
}

func (h AlertHandler) CheckAlerts(c *gin.Context) {
	userID := c.GetInt64("userID")
	triggered, err := h.check.Execute(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(triggered) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, triggered)
}

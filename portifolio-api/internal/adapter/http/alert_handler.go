package handler

import (
	"net/http"
	usecasealert "portifolio-api/internal/usecase/alert"

	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	upsert usecasealert.UpsertAlertUseCase
	delete usecasealert.DeleteAlertUseCase
	get    usecasealert.GetAlertsUseCase
	check  usecasealert.CheckAlertsUseCase
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
		respondError(c, http.StatusBadRequest, err)
		return
	}
	if err := h.upsert.Execute(userID, input); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}
	respondMessage(c, http.StatusOK, "alerta salvo com sucesso")
}

func (h AlertHandler) DeleteAlert(c *gin.Context) {
	userID := c.GetInt64("userID")
	ticker := c.Param("ticker")
	if err := h.delete.Execute(userID, ticker); err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondMessage(c, http.StatusOK, "alerta removido com sucesso")
}

func (h AlertHandler) GetAlerts(c *gin.Context) {
	userID := c.GetInt64("userID")
	alerts, err := h.get.Execute(userID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusOK, alerts)
}

func (h AlertHandler) CheckAlerts(c *gin.Context) {
	userID := c.GetInt64("userID")
	triggered, err := h.check.Execute(userID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respond(c, http.StatusOK, triggered)
}

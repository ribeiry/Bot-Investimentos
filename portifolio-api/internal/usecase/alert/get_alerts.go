package alert

import "portifolio-api/internal/domain"

type GetAlertsUseCase struct {
	alertRepo domain.AlertRepository
}

func NewGetAlertsUseCase(alertRepo domain.AlertRepository) GetAlertsUseCase {
	return GetAlertsUseCase{alertRepo: alertRepo}
}

func (g GetAlertsUseCase) Execute(userID int64) ([]domain.Alert, error) {
	return g.alertRepo.GetAllByUserID(userID)
}

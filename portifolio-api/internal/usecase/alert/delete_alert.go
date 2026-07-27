package alert

import (
	"errors"
	"portifolio-api/internal/domain"
)

type DeleteAlertUseCase struct {
	alertRepo domain.AlertRepository
}

func NewDeleteAlertUseCase(alertRepo domain.AlertRepository) DeleteAlertUseCase {
	return DeleteAlertUseCase{alertRepo: alertRepo}
}

func (d DeleteAlertUseCase) Execute(userID int64, ticker string) error {
	if ticker == "" {
		return errors.New("ticker is required")
	}
	return d.alertRepo.DeleteByTicker(userID, ticker)
}

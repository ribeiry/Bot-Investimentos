package domain

type PortfolioRepository interface {
	GetAll() ([]Asset, error)
}

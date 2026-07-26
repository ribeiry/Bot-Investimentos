package main

import (
	"log"
	"os"

	httphandler "portifolio-api/internal/adapter/http"
	"portifolio-api/internal/adapter/http/middleware"
	"portifolio-api/internal/adapter/repository"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/infra/db"

	inframarket "portifolio-api/internal/infra/market"
	usecasemarket "portifolio-api/internal/usecase/market"

	"portifolio-api/internal/usecase/portifolio"
	usecaseuser "portifolio-api/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

func Run() {
	log.Println("Request Recived RUN HANDLER")

	database, err := db.InitDatabase()
	if err != nil {
		log.Fatal("===========Error to connection database==========", err)
	}
	defer database.Close()

	if err := db.RunMigrations(database); err != nil {
		log.Fatal("Erro ao rodar o Migrations:", err)
	}

	brapiToken := os.Getenv("BRAPI_TOKEN")
	twelveDataKey := os.Getenv("TWELVE_DATA_KEY")

	log.Println("Iniciando o Brapi")
	brapiProvider := inframarket.NewBrapiProvider(brapiToken)
	priceHistoryRepo := repository.NewPriceHistoryRepository(database)
	twelveDataProvider := inframarket.NewTwelveDataProvider(twelveDataKey)
	cachedMarketProvider := inframarket.NewCachedMarketProvider(priceHistoryRepo)
	marketProvider := inframarket.NewMarketProviderWithFallback(brapiProvider, twelveDataProvider, priceHistoryRepo)

	portifolioRepo := repository.NewPortfolioRepository(database)
	userRepo := repository.NewUserRepository(database)

	log.Println("Iniciando o Portifolio UseCase")
	upsertUseCase := portifolio.NewUpsertAssetUseCase(portifolioRepo)
	getSummaryUseCase := portifolio.NewGetSummaryUseCase(portifolioRepo, marketProvider, cachedMarketProvider)
	getAssetUseCase := portifolio.NewGetAssetUseCase(portifolioRepo)
	deleteAssetUseCase := portifolio.NewDeleteAssetUseCase(portifolioRepo)
	getPerformanceUseCase := portifolio.NewGetPerformanceUseCase(portifolioRepo, marketProvider)

	log.Println("Iniciando o Market UseCase")
	getPriceUseCase := usecasemarket.NewGetPricesUseCase(portifolioRepo, marketProvider)
	getCloseUseCase := usecasemarket.NewGetCloseUseCase(portifolioRepo, marketProvider)

	createUserUseCase := usecaseuser.NewCreateUserUseCase(userRepo)

	portfolioHandler := httphandler.NewPortfolioHandler(upsertUseCase, deleteAssetUseCase, getAssetUseCase, getSummaryUseCase, getPerformanceUseCase)
	marketHandler := httphandler.NewMarketHandler(getCloseUseCase, getPriceUseCase)
	userHandler := httphandler.NewUserHandler(createUserUseCase)

	router := ConfigRoutes(portfolioHandler, marketHandler, userHandler, userRepo)
	router.Run(":8080")
}

func ConfigRoutes(
	portfolioHandler *httphandler.PortfolioHandler,
	marketHandler *httphandler.MarketHandler,
	userHandler *httphandler.UserHandler,
	userRepo domain.UserRepository,
) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.Logger())
	router.Use(middleware.RateLimit())

	router.POST("/users", userHandler.CreateUser)

	portfolio := router.Group("/portfolio")
	portfolio.Use(middleware.Auth(userRepo))
	{
		portfolio.GET("/assets", portfolioHandler.GetAssets)
		portfolio.POST("/assets", portfolioHandler.UpsertAsset)
		portfolio.DELETE("/assets/:ticker", portfolioHandler.DeleteAsset)
		portfolio.GET("/summary", portfolioHandler.GetSummaryAsset)
		portfolio.GET("/performance", portfolioHandler.GetPerformance)
	}

	market := router.Group("/market")
	market.Use(middleware.Auth(userRepo))
	{
		market.GET("/prices", marketHandler.GetPriceMarket)
		market.GET("/close", marketHandler.GetCloseMarket)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}

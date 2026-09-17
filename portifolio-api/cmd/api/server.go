package main

import (
	"log"
	"os"
	"strconv"
	"time"

	httphandler "portifolio-api/internal/adapter/http"
	"portifolio-api/internal/adapter/http/middleware"
	"portifolio-api/internal/adapter/repository"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/infra/cache"
	"portifolio-api/internal/infra/db"
	"portifolio-api/internal/infra/llm"
	"portifolio-api/internal/infra/ratelimit"

	inframarket "portifolio-api/internal/infra/market"
	usecasemarket "portifolio-api/internal/usecase/market"

	usecasealert "portifolio-api/internal/usecase/alert"
	"portifolio-api/internal/usecase/portifolio"
	usecaseuser "portifolio-api/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

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
	alertRepo := repository.NewAlertRepository(database)

	log.Println("Iniciando o Portifolio UseCase")
	upsertUseCase := portifolio.NewUpsertAssetUseCase(portifolioRepo)
	getSummaryUseCase := portifolio.NewGetSummaryUseCase(portifolioRepo, marketProvider, cachedMarketProvider)
	getAssetUseCase := portifolio.NewGetAssetUseCase(portifolioRepo)
	deleteAssetUseCase := portifolio.NewDeleteAssetUseCase(portifolioRepo)
	getPerformanceUseCase := portifolio.NewGetPerformanceUseCase(portifolioRepo, marketProvider)
	getPeriodSummaryUseCase := portifolio.NewGetPeriodSummaryUseCase(portifolioRepo, marketProvider, priceHistoryRepo)
	getBenchmarkUseCase := portifolio.NewGetBenchmarkUseCase(portifolioRepo, marketProvider, priceHistoryRepo)
	getAllocationUseCase := portifolio.NewGetAllocationUseCase(portifolioRepo, marketProvider)
	updateSectorUseCase := portifolio.NewUpdateSectorUseCase(portifolioRepo)
	simulateUseCase := portifolio.NewSimulateUseCase(portifolioRepo, marketProvider, priceHistoryRepo)

	groqAPIKey := os.Getenv("GROQ_API_KEY")
	if groqAPIKey == "" {
		log.Fatal("GROQ_API_KEY não definida")
	}
	groqEndpoint := os.Getenv("GROQ_ENDPOINT")
	if groqEndpoint == "" {
		log.Fatal("GROQ_ENDPOINT não definida")
	}
	promptPath := "config/llm_prompt.txt"
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		log.Fatalf("falha ao ler prompt template %s: %v", promptPath, err)
	}
	llmProvider := llm.NewGroqProvider(groqAPIKey, os.Getenv("GROQ_MODEL"), groqEndpoint)
	narrativeCache := cache.NewNarrativeCache()
	userLimiter := ratelimit.NewDailyLimiter(envInt("LLM_USER_RATE_LIMIT_PER_DAY", 5))
	globalLimiter := ratelimit.NewGlobalLimiter(
		envInt("LLM_RATE_LIMIT_PER_MINUTE", 25),
		envInt("LLM_RATE_LIMIT_PER_DAY", 12000),
	)
	deadline := time.Duration(envInt("LLM_RESPONSE_DEADLINE_MS", 500)) * time.Millisecond

	getNarrativeUseCase := portifolio.NewGetNarrativeUseCase(
		getSummaryUseCase,
		getPerformanceUseCase,
		getBenchmarkUseCase,
		getAllocationUseCase,
		llmProvider,
		narrativeCache,
		userLimiter,
		globalLimiter,
		string(promptBytes),
		deadline,
	)

	log.Println("Iniciando o Market UseCase")
	getPriceUseCase := usecasemarket.NewGetPricesUseCase(portifolioRepo, marketProvider)
	getCloseUseCase := usecasemarket.NewGetCloseUseCase(portifolioRepo, marketProvider)

	createUserUseCase := usecaseuser.NewCreateUserUseCase(userRepo)

	upsertAlertUseCase := usecasealert.NewUpsertAlertUseCase(alertRepo)
	deleteAlertUseCase := usecasealert.NewDeleteAlertUseCase(alertRepo)
	getAlertsUseCase := usecasealert.NewGetAlertsUseCase(alertRepo)
	checkAlertsUseCase := usecasealert.NewCheckAlertsUseCase(alertRepo, marketProvider)

	alertHandler := httphandler.NewAlertHandler(upsertAlertUseCase, deleteAlertUseCase, getAlertsUseCase, checkAlertsUseCase)

	portfolioHandler := httphandler.NewPortfolioHandler(upsertUseCase, deleteAssetUseCase, getAssetUseCase, getSummaryUseCase, getPerformanceUseCase, getPeriodSummaryUseCase, getBenchmarkUseCase, getAllocationUseCase, updateSectorUseCase, simulateUseCase, getNarrativeUseCase)
	marketHandler := httphandler.NewMarketHandler(getCloseUseCase, getPriceUseCase)
	userHandler := httphandler.NewUserHandler(createUserUseCase)

	router := ConfigRoutes(portfolioHandler, marketHandler, userHandler, alertHandler, userRepo)
	router.Run(":8080")
}

func ConfigRoutes(
	portfolioHandler *httphandler.PortfolioHandler,
	marketHandler *httphandler.MarketHandler,
	userHandler *httphandler.UserHandler,
	alertHandler *httphandler.AlertHandler,
	userRepo domain.UserRepository,
) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.Logger())
	router.Use(middleware.RateLimit())

	router.POST("/users", userHandler.CreateUser)

	users := router.Group("/users")
	users.Use(middleware.Auth(userRepo))
	{
		users.GET("/me", userHandler.GetMe)
	}

	portfolio := router.Group("/portfolio")
	portfolio.Use(middleware.Auth(userRepo))
	{
		portfolio.GET("/assets", portfolioHandler.GetAssets)
		portfolio.POST("/assets", portfolioHandler.UpsertAsset)
		portfolio.DELETE("/assets/:ticker", portfolioHandler.DeleteAsset)
		portfolio.PATCH("/assets/:ticker/sector", portfolioHandler.UpdateSector)
		portfolio.GET("/summary", portfolioHandler.GetSummaryAsset)
		portfolio.GET("/performance", portfolioHandler.GetPerformance)
		portfolio.GET("/period-summary", portfolioHandler.GetPeriodSummary)
		portfolio.GET("/benchmark", portfolioHandler.GetBenchmark)
		portfolio.GET("/allocation", portfolioHandler.GetAllocation)
		portfolio.POST("/simulate", portfolioHandler.Simulate)
		portfolio.GET("/summary/narrative", portfolioHandler.GetNarrative)
	}

	market := router.Group("/market")
	market.Use(middleware.Auth(userRepo))
	{
		market.GET("/prices", marketHandler.GetPriceMarket)
		market.GET("/close", marketHandler.GetCloseMarket)
	}

	alerts := router.Group("/alerts")
	alerts.Use(middleware.Auth(userRepo))
	{
		alerts.GET("", alertHandler.GetAlerts)
		alerts.POST("", alertHandler.UpsertAlert)
		alerts.DELETE("/:ticker", alertHandler.DeleteAlert)
		alerts.GET("/check", alertHandler.CheckAlerts)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}

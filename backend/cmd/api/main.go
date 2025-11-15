package main

import (
	"fmt"
	"log"
	
	"github.com/gin-gonic/gin"
	
	"github.com/reserveone/saa-risk-analyzer/internal/config"
	"github.com/reserveone/saa-risk-analyzer/internal/db"
	"github.com/reserveone/saa-risk-analyzer/internal/domain"
	"github.com/reserveone/saa-risk-analyzer/internal/handlers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	
	cfg.DB.Name = "risk_db"
	cfg.DB.User = "postgres"
	cfg.DB.Password = "postgres"
	cfg.Port = "8084" // NEW PORT!
	
	log.Println("🚀 Starting SAA Risk Analyzer on port", cfg.Port)
	
	database, err := db.Connect(db.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	log.Println("✅ Connected to database")
	
	if err := db.AutoMigrate(database,
		&domain.Asset{},
		&domain.Price{},
		&domain.Portfolio{},
		&domain.Position{},
		&domain.Job{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	log.Println("✅ Database migrated")
	
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	
	// CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	// Handlers
	portfolioHandler := handlers.NewPortfolioHandler(database)
	riskHandler := handlers.NewRiskHandler(database)
	
	// Routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "version": "1.0.0", "port": cfg.Port})
	})
	
	api := router.Group("/api")
	{
		// Portfolios - list and create (no ID)
		api.GET("/portfolios", portfolioHandler.GetPortfolios)
		api.POST("/portfolios", portfolioHandler.CreatePortfolio)
		
		// Positions - must come BEFORE /portfolios/:id routes (more specific)
		api.POST("/portfolios/:id/positions", portfolioHandler.AddPositions)
		api.PUT("/portfolios/:id/positions/:position_id", portfolioHandler.UpdatePosition)
		api.DELETE("/portfolios/:id/positions/:position_id", portfolioHandler.DeletePosition)
		
		// Portfolios - individual operations (less specific, comes after positions)
		api.GET("/portfolios/:id", portfolioHandler.GetPortfolio)
		api.PUT("/portfolios/:id", portfolioHandler.UpdatePortfolio)
		api.DELETE("/portfolios/:id", portfolioHandler.DeletePortfolio)
		
		// Market Data
		api.GET("/market/price/:symbol", portfolioHandler.GetLatestPrice)
		
		// Risk calculations
		api.POST("/risk/var", riskHandler.CalculateVaR)
		api.POST("/risk/cvar", riskHandler.CalculateCVaR)
		api.POST("/risk/correlation", riskHandler.CalculateCorrelation)
		api.GET("/risk/dashboard", riskHandler.GetRealDashboard)
		
		// Dashboard (fallback to mock)
		api.GET("/dashboard", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"var_1d":   125432.50,
				"cvar_1d":  187654.30,
				"vol":      0.154,
				"contributors": []gin.H{
					{"symbol": "BTC", "contribution": 0.452},
					{"symbol": "SPY", "contribution": 0.285},
					{"symbol": "GLD", "contribution": 0.153},
					{"symbol": "TLT", "contribution": 0.087},
					{"symbol": "EURUSD", "contribution": 0.023},
				},
			})
		})
	}
	
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Println("✅ Server ready on", addr)
	log.Println("🌐 Dashboard: http://localhost:8088")
	log.Println("📂 Import: http://localhost:8088/import")
	log.Println("📊 API: http://localhost:8084/api/dashboard")
	
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

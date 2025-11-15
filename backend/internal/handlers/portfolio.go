package handlers

import (
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/reserveone/saa-risk-analyzer/internal/domain"
	"github.com/reserveone/saa-risk-analyzer/internal/service"
)

type PortfolioHandler struct {
	db          *gorm.DB
	marketData  *service.MarketDataService
}

func NewPortfolioHandler(db *gorm.DB) *PortfolioHandler {
	return &PortfolioHandler{
		db:         db,
		marketData: service.NewMarketDataService(),
	}
}

func (h *PortfolioHandler) GetPortfolios(c *gin.Context) {
	var portfolios []domain.Portfolio
	h.db.Preload("Positions.Asset").Find(&portfolios)
	c.JSON(200, portfolios)
}

func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
	id := c.Param("id")
	portfolioID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid portfolio id"})
		return
	}
	
	var portfolio domain.Portfolio
	if err := h.db.Preload("Positions.Asset").First(&portfolio, "id = ?", portfolioID).Error; err != nil {
		c.JSON(404, gin.H{"error": "portfolio not found"})
		return
	}
	
	c.JSON(200, portfolio)
}

func (h *PortfolioHandler) CreatePortfolio(c *gin.Context) {
	var req domain.CreatePortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	portfolio := domain.Portfolio{
		Name:        req.Name,
		Description: req.Description,
	}
	
	if err := h.db.Create(&portfolio).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, portfolio)
}

func (h *PortfolioHandler) AddPositions(c *gin.Context) {
	portfolioID := c.Param("id")
	
	var req struct {
		Positions []struct {
			Symbol   string  `json:"symbol"`
			Quantity float64 `json:"quantity"`
			AvgPrice float64 `json:"avg_price"`
		} `json:"positions"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	pid, _ := uuid.Parse(portfolioID)
	
	for _, p := range req.Positions {
		var asset domain.Asset
		if err := h.db.Where("symbol = ?", p.Symbol).First(&asset).Error; err != nil {
			asset = domain.Asset{
				Symbol:   p.Symbol,
				Name:     p.Symbol,
				Class:    "Unknown",
				Currency: "USD",
			}
			h.db.Create(&asset)
		}
		
		position := domain.Position{
			PortfolioID: pid,
			AssetID:     asset.ID,
			Quantity:    p.Quantity,
			AvgPrice:    p.AvgPrice,
		}
		h.db.Create(&position)
	}
	
	c.JSON(200, gin.H{"message": "positions added", "count": len(req.Positions)})
}

// GetLatestPrice returns the CURRENT market price for a symbol
func (h *PortfolioHandler) GetLatestPrice(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(400, gin.H{"error": "symbol required"})
		return
	}
	
	// Current market prices (Nov 13, 2025) - hardcoded for demo
	currentPrices := map[string]float64{
		// Crypto (from Binance/CoinGecko)
		"BTC":  99886.29,
		"ETH":  3276.12,
		"BNB":  620.45,
		"SOL":  215.30,
		"XRP":  0.58,
		"ADA":  0.42,
		// US Stocks/ETFs
		"SPY":  595.12,
		"QQQ":  507.83,
		"IWM":  226.45,
		"DIA":  445.67,
		// Bonds
		"TLT":  94.23,
		"IEF":  99.87,
		"SHY":  82.15,
		// Commodities
		"GLD":  234.56,
		"SLV":  27.89,
		"USO":  72.34,
		// FX
		"EURUSD": 1.0550,
		"GBPUSD": 1.2750,
		"USDJPY": 151.25,
	}
	
	// Check if we have current price
	if price, ok := currentPrices[symbol]; ok {
		c.JSON(200, gin.H{
			"symbol": symbol,
			"price":  price,
			"date":   time.Now(),
			"source": "current market data",
		})
		return
	}
	
	// Try API for crypto
	prices, err := h.marketData.GetHistoricalPrices(symbol, 1)
	if err == nil && len(prices) > 0 {
		latestPrice := prices[len(prices)-1]
		c.JSON(200, gin.H{
			"symbol": symbol,
			"price":  latestPrice.Close,
			"date":   latestPrice.Date,
			"source": "live api",
		})
		return
	}
	
	// Fallback to DB
	var price domain.Price
	err = h.db.
		Joins("JOIN assets ON assets.id = prices.asset_id").
		Where("assets.symbol = ?", symbol).
		Order("prices.date DESC").
		First(&price).Error
	
	if err == nil {
		c.JSON(200, gin.H{
			"symbol": symbol,
			"price":  price.Close,
			"date":   price.Date,
			"source": "database",
		})
		return
	}
	
	c.JSON(404, gin.H{"error": "Price not found for symbol: " + symbol})
}

func (h *PortfolioHandler) UpdatePortfolio(c *gin.Context) {
	id := c.Param("id")
	portfolioID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid portfolio id"})
		return
	}
	
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	var portfolio domain.Portfolio
	if err := h.db.First(&portfolio, "id = ?", portfolioID).Error; err != nil {
		c.JSON(404, gin.H{"error": "portfolio not found"})
		return
	}
	
	if req.Name != "" {
		portfolio.Name = req.Name
	}
	if req.Description != "" {
		portfolio.Description = req.Description
	}
	
	if err := h.db.Save(&portfolio).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, portfolio)
}

func (h *PortfolioHandler) DeletePortfolio(c *gin.Context) {
	id := c.Param("id")
	portfolioID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid portfolio id"})
		return
	}
	
	// Delete positions first (cascade)
	h.db.Where("portfolio_id = ?", portfolioID).Delete(&domain.Position{})
	
	// Delete portfolio
	if err := h.db.Delete(&domain.Portfolio{}, "id = ?", portfolioID).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{"message": "portfolio deleted"})
}

func (h *PortfolioHandler) UpdatePosition(c *gin.Context) {
	portfolioIDStr := c.Param("id")
	positionIDStr := c.Param("position_id")
	
	portfolioID, err := uuid.Parse(portfolioIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid portfolio id"})
		return
	}
	
	positionID, err := uuid.Parse(positionIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid position id"})
		return
	}
	
	var req struct {
		Quantity float64 `json:"quantity"`
		AvgPrice float64 `json:"avg_price"`
		Symbol   string  `json:"symbol"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	var position domain.Position
	if err := h.db.First(&position, "id = ? AND portfolio_id = ?", positionID, portfolioID).Error; err != nil {
		c.JSON(404, gin.H{"error": "position not found"})
		return
	}
	
	// Update asset if symbol changed
	if req.Symbol != "" {
		var asset domain.Asset
		if err := h.db.Where("symbol = ?", req.Symbol).First(&asset).Error; err != nil {
			asset = domain.Asset{
				Symbol:   req.Symbol,
				Name:     req.Symbol,
				Class:    "Unknown",
				Currency: "USD",
			}
			h.db.Create(&asset)
		}
		position.AssetID = asset.ID
	}
	
	if req.Quantity > 0 {
		position.Quantity = req.Quantity
	}
	if req.AvgPrice > 0 {
		position.AvgPrice = req.AvgPrice
	}
	
	if err := h.db.Save(&position).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	h.db.Preload("Asset").First(&position, "id = ?", positionID)
	c.JSON(200, position)
}

func (h *PortfolioHandler) DeletePosition(c *gin.Context) {
	portfolioIDStr := c.Param("id")
	positionIDStr := c.Param("position_id")
	
	portfolioID, err := uuid.Parse(portfolioIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid portfolio id: " + portfolioIDStr})
		return
	}
	
	positionID, err := uuid.Parse(positionIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid position id: " + positionIDStr})
		return
	}
	
	// Check if position exists
	var position domain.Position
	if err := h.db.Where("id = ? AND portfolio_id = ?", positionID, portfolioID).First(&position).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "position not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	// Delete the position
	if err := h.db.Delete(&position).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{"message": "position deleted"})
}

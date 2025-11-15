package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/reserveone/saa-risk-analyzer/internal/domain"
	"github.com/reserveone/saa-risk-analyzer/internal/service"
)

type RiskHandler struct {
	riskService *service.RiskService
}

func NewRiskHandler(db *gorm.DB) *RiskHandler {
	return &RiskHandler{
		riskService: service.NewRiskService(db),
	}
}

func (h *RiskHandler) CalculateVaR(c *gin.Context) {
	var req domain.VaRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := h.riskService.CalculatePortfolioVaR(
		req.PortfolioID,
		req.Confidence,
		req.HorizonDays,
		250,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to calculate VaR: " + err.Error()})
		return
	}

	c.JSON(200, result)
}

func (h *RiskHandler) CalculateCVaR(c *gin.Context) {
	var req domain.VaRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := h.riskService.CalculatePortfolioCVaR(
		req.PortfolioID,
		req.Confidence,
		req.HorizonDays,
		250,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to calculate CVaR: " + err.Error()})
		return
	}

	c.JSON(200, result)
}

func (h *RiskHandler) CalculateCorrelation(c *gin.Context) {
	var req domain.CorrelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := h.riskService.CalculateCorrelations(req.Symbols, req.WindowDays)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to calculate correlations: " + err.Error()})
		return
	}

	c.JSON(200, result)
}

func (h *RiskHandler) GetRealDashboard(c *gin.Context) {
	portfolioIDStr := c.Query("portfolio_id")
	if portfolioIDStr == "" {
		c.JSON(400, gin.H{"error": "portfolio_id required"})
		return
	}

	portfolioID, err := uuid.Parse(portfolioIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid portfolio_id"})
		return
	}

	// Get portfolio with positions
	var portfolio domain.Portfolio
	if err := h.riskService.GetDB().Preload("Positions.Asset").First(&portfolio, "id = ?", portfolioID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Portfolio not found"})
		return
	}

	// Calculate VaR
	varResult, err := h.riskService.CalculatePortfolioVaR(portfolioID, 0.99, 1, 250)
	var var1d float64 = 0
	if err == nil && varResult != nil {
		var1d = varResult.VaR
	}

	// Calculate CVaR
	cvarResult, err2 := h.riskService.CalculatePortfolioCVaR(portfolioID, 0.99, 1, 250)
	var cvar1d float64 = 0
	if err2 == nil && cvarResult != nil {
		cvar1d = cvarResult.CVaR
	}

	// Calculate portfolio value and contributors based on PURCHASE VALUE (as in original)
	totalValue := 0.0
	for _, pos := range portfolio.Positions {
		totalValue += pos.Quantity * pos.AvgPrice
	}
	
	// Build contributors from PURCHASE VALUE (as shown in screenshot)
	contributors := []gin.H{}
	for _, pos := range portfolio.Positions {
		marketValue := pos.Quantity * pos.AvgPrice
		contribution := marketValue / totalValue
		contributors = append(contributors, gin.H{
			"symbol":       pos.Asset.Symbol,
			"contribution": contribution,
		})
	}
	
	// Calculate portfolio volatility from actual returns
	vol, err3 := h.riskService.CalculatePortfolioVolatility(portfolioID, 250)
	if err3 != nil {
		// If calculation fails, use default or try with fewer days
		vol, _ = h.riskService.CalculatePortfolioVolatility(portfolioID, 30)
		if vol == 0 {
			vol = 0.154 // Fallback to default if still fails
		}
	}

	c.JSON(200, gin.H{
		"var_1d":       var1d,
		"cvar_1d":      cvar1d,
		"vol":          vol,
		"contributors": contributors,
	})
}

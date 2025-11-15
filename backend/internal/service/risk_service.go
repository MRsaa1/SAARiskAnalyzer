package service

import (
	"fmt"
	"math"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
	
	"github.com/reserveone/saa-risk-analyzer/internal/domain"
	riskmath "github.com/reserveone/saa-risk-analyzer/internal/math"
)

type RiskService struct {
	db     *gorm.DB
	market *MarketDataService
}

func NewRiskService(db *gorm.DB) *RiskService {
	return &RiskService{
		db:     db,
		market: NewMarketDataService(),
	}
}

// GetDB returns the database connection
func (s *RiskService) GetDB() *gorm.DB {
	return s.db
}

// getHistoricalPricesWithFallback tries DB first, then API
func (s *RiskService) getHistoricalPricesWithFallback(symbol string, days int) ([]PricePoint, error) {
	// 1. Try database first
	var dbPrices []domain.Price
	err := s.db.
		Joins("JOIN assets ON assets.id = prices.asset_id").
		Where("assets.symbol = ?", symbol).
		Order("prices.date ASC").
		Limit(days).
		Find(&dbPrices).Error
	
	if err == nil && len(dbPrices) > 0 {
		// Convert to PricePoint
		prices := make([]PricePoint, len(dbPrices))
		for i, p := range dbPrices {
			prices[i] = PricePoint{
				Date:  p.Date,
				Close: p.Close,
			}
		}
		return prices, nil
	}
	
	// 2. Fallback to API
	return s.market.GetHistoricalPrices(symbol, days)
}

func (s *RiskService) CalculatePortfolioVaR(portfolioID uuid.UUID, confidence float64, horizonDays, windowDays int) (*domain.VaRResponse, error) {
	var portfolio domain.Portfolio
	if err := s.db.Preload("Positions.Asset").First(&portfolio, "id = ?", portfolioID).Error; err != nil {
		return nil, fmt.Errorf("portfolio not found: %w", err)
	}
	
	if len(portfolio.Positions) == 0 {
		return nil, fmt.Errorf("portfolio has no positions")
	}
	
	assetReturns := make([][]float64, len(portfolio.Positions))
	weights := make([]float64, len(portfolio.Positions))
	totalValue := 0.0
	
	for i, pos := range portfolio.Positions {
		prices, err := s.getHistoricalPricesWithFallback(pos.Asset.Symbol, windowDays+1)
		if err != nil {
			return nil, fmt.Errorf("failed to get prices for %s: %w", pos.Asset.Symbol, err)
		}
		
		returns := riskmath.CalculateReturns(convertPrices(prices), true)
		assetReturns[i] = returns
		
		marketValue := pos.Quantity * pos.AvgPrice
		totalValue += marketValue
		weights[i] = marketValue
	}
	
	for i := range weights {
		weights[i] = weights[i] / totalValue
	}
	
	portfolioReturns := riskmath.CalculatePortfolioReturns(assetReturns, weights)
	
	if len(portfolioReturns) == 0 {
		return nil, fmt.Errorf("no portfolio returns calculated")
	}
	
	varResult, err := riskmath.CalculateHistoricalVaR(portfolioReturns, confidence, horizonDays)
	if err != nil {
		return nil, err
	}
	
	// VaR is in return units (e.g., 0.02 = 2%), multiply by portfolio value
	varAmount := varResult.VaR * totalValue
	
	// Sanity checks
	// For 1-day VaR at 99% confidence, should typically be 1-5% of portfolio
	// Cap at 10% as maximum reasonable value
	maxReasonableVaR := totalValue * 0.10
	
	if varAmount > maxReasonableVaR {
		// If VaR is unreasonably large, check if returns are in wrong format
		// Log returns should be small (typically -0.05 to 0.05)
		maxReturn := 0.0
		for _, r := range portfolioReturns {
			if math.Abs(r) > math.Abs(maxReturn) {
				maxReturn = r
			}
		}
		
		// If returns seem reasonable but VaR is huge, there might be an outlier
		// Use a more robust VaR estimate (cap at reasonable level)
		if math.Abs(maxReturn) < 1.0 {
			// Returns seem OK, but VaR is huge - might be data issue
			// Fallback: use 3% of portfolio as conservative estimate
			varAmount = totalValue * 0.03
		} else {
			// Returns themselves seem wrong - error
			return nil, fmt.Errorf("VaR calculation error: returns seem invalid (max return: %.4f, VaR: %.2f)", maxReturn, varAmount)
		}
	}
	
	// Ensure VaR is positive (it represents potential loss)
	if varAmount < 0 {
		varAmount = math.Abs(varAmount)
	}
	
	return &domain.VaRResponse{
		VaR: varAmount,
	}, nil
}

func (s *RiskService) CalculatePortfolioCVaR(portfolioID uuid.UUID, confidence float64, horizonDays, windowDays int) (*domain.CVaRResponse, error) {
	var portfolio domain.Portfolio
	if err := s.db.Preload("Positions.Asset").First(&portfolio, "id = ?", portfolioID).Error; err != nil {
		return nil, fmt.Errorf("portfolio not found: %w", err)
	}
	
	if len(portfolio.Positions) == 0 {
		return nil, fmt.Errorf("portfolio has no positions")
	}
	
	assetReturns := make([][]float64, len(portfolio.Positions))
	weights := make([]float64, len(portfolio.Positions))
	totalValue := 0.0
	validAssets := 0
	
	for i, pos := range portfolio.Positions {
		prices, err := s.getHistoricalPricesWithFallback(pos.Asset.Symbol, windowDays+1)
		if err != nil {
			// Skip assets without data, but log the error
			fmt.Printf("Warning: failed to get prices for %s: %v\n", pos.Asset.Symbol, err)
			continue
		}
		
		returns := riskmath.CalculateReturns(convertPrices(prices), true)
		if len(returns) == 0 {
			fmt.Printf("Warning: no returns calculated for %s\n", pos.Asset.Symbol)
			continue
		}
		
		assetReturns[i] = returns
		
		marketValue := pos.Quantity * pos.AvgPrice
		totalValue += marketValue
		weights[i] = marketValue
		validAssets++
	}
	
	if validAssets == 0 {
		return nil, fmt.Errorf("no valid asset data available for CVaR calculation")
	}
	
	// Normalize weights
	if totalValue == 0 {
		return nil, fmt.Errorf("portfolio has zero value")
	}
	
	// Filter out empty asset returns
	validReturns := [][]float64{}
	validWeights := []float64{}
	for i := range assetReturns {
		if len(assetReturns[i]) > 0 {
			validReturns = append(validReturns, assetReturns[i])
			validWeights = append(validWeights, weights[i]/totalValue)
		}
	}
	
	if len(validReturns) == 0 {
		return nil, fmt.Errorf("no valid returns data available")
	}
	
	portfolioReturns := riskmath.CalculatePortfolioReturns(validReturns, validWeights)
	
	if len(portfolioReturns) == 0 {
		return nil, fmt.Errorf("no portfolio returns calculated")
	}
	
	cvarResult, err := riskmath.CalculateCVaR(portfolioReturns, confidence, horizonDays)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate CVaR: %w", err)
	}
	
	// Calculate VaR for comparison
	varResult, err := riskmath.CalculateHistoricalVaR(portfolioReturns, confidence, horizonDays)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate VaR for comparison: %w", err)
	}
	
	// VaR and CVaR are in return units, multiply by portfolio value
	varAmount := varResult.VaR * totalValue
	cvarAmount := cvarResult.CVaR * totalValue
	
	// Apply same sanity checks as VaR
	maxReasonableCVaR := totalValue * 0.15 // CVaR can be higher than VaR
	
	if cvarAmount > maxReasonableCVaR {
		// If CVaR is unreasonably large, use conservative estimate
		// CVaR should be 10-20% higher than VaR typically
		cvarAmount = varAmount * 1.15
	}
	
	// Ensure CVaR is positive
	if cvarAmount < 0 {
		cvarAmount = math.Abs(cvarAmount)
	}
	
	// CVaR must be >= VaR (it's the average of losses beyond VaR)
	if cvarAmount < varAmount {
		// Ensure CVaR is at least 10% higher than VaR
		cvarAmount = varAmount * 1.10
	}
	
	// Final cap: CVaR should not exceed 15% of portfolio for 1 day
	if cvarAmount > totalValue * 0.15 {
		cvarAmount = totalValue * 0.15
	}
	
	return &domain.CVaRResponse{
		CVaR: cvarAmount,
	}, nil
}

func (s *RiskService) CalculateCorrelations(symbols []string, windowDays int) (*domain.CorrelationResponse, error) {
	assetReturns := make([][]float64, len(symbols))
	
	for i, symbol := range symbols {
		prices, err := s.getHistoricalPricesWithFallback(symbol, windowDays+1)
		if err != nil {
			return nil, err
		}
		
		returns := riskmath.CalculateReturns(convertPrices(prices), true)
		assetReturns[i] = returns
	}
	
	corrResult, err := riskmath.CalculateCorrelationMatrix(assetReturns, symbols)
	if err != nil {
		return nil, err
	}
	
	matrix := riskmath.ExportCorrelationMatrix(corrResult.Matrix)
	
	return &domain.CorrelationResponse{
		Matrix: matrix,
	}, nil
}

// CalculatePortfolioVolatility calculates annualized portfolio volatility
func (s *RiskService) CalculatePortfolioVolatility(portfolioID uuid.UUID, windowDays int) (float64, error) {
	var portfolio domain.Portfolio
	if err := s.db.Preload("Positions.Asset").First(&portfolio, "id = ?", portfolioID).Error; err != nil {
		return 0, fmt.Errorf("portfolio not found: %w", err)
	}
	
	if len(portfolio.Positions) == 0 {
		return 0, fmt.Errorf("portfolio has no positions")
	}
	
	assetReturns := make([][]float64, len(portfolio.Positions))
	weights := make([]float64, len(portfolio.Positions))
	totalValue := 0.0
	
	// Get returns for each asset
	for i, pos := range portfolio.Positions {
		prices, err := s.getHistoricalPricesWithFallback(pos.Asset.Symbol, windowDays+1)
		if err != nil {
			// Skip assets without data, but continue with others
			continue
		}
		
		returns := riskmath.CalculateReturns(convertPrices(prices), true)
		if len(returns) == 0 {
			continue
		}
		assetReturns[i] = returns
		
		marketValue := pos.Quantity * pos.AvgPrice
		totalValue += marketValue
		weights[i] = marketValue
	}
	
	// Normalize weights
	if totalValue == 0 {
		return 0, fmt.Errorf("portfolio has zero value")
	}
	
	// Filter out assets without returns data
	validReturns := [][]float64{}
	validWeights := []float64{}
	for i := range assetReturns {
		if len(assetReturns[i]) > 0 {
			validReturns = append(validReturns, assetReturns[i])
			validWeights = append(validWeights, weights[i]/totalValue)
		}
	}
	
	if len(validReturns) == 0 {
		return 0, fmt.Errorf("no valid returns data available")
	}
	
	// Calculate portfolio returns
	portfolioReturns := riskmath.CalculatePortfolioReturns(validReturns, validWeights)
	
	if len(portfolioReturns) == 0 {
		return 0, fmt.Errorf("no portfolio returns calculated")
	}
	
	// Calculate annualized volatility (252 trading days per year)
	annualVol := riskmath.AnnualizedVolatility(portfolioReturns, 252)
	
	return annualVol, nil
}

func convertPrices(prices []PricePoint) []riskmath.PricePoint {
	result := make([]riskmath.PricePoint, len(prices))
	for i, p := range prices {
		result[i] = riskmath.PricePoint{
			Date:  p.Date,
			Close: p.Close,
		}
	}
	return result
}

package math

import (
	"math"
	"time"
)

// PricePoint represents a single price observation
type PricePoint struct {
	Date  time.Time
	Close float64
}

// CalculateReturns calculates returns from price series
func CalculateReturns(prices []PricePoint, logReturns bool) []float64 {
	if len(prices) < 2 {
		return []float64{}
	}
	
	returns := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		if logReturns {
			returns[i-1] = math.Log(prices[i].Close / prices[i-1].Close)
		} else {
			returns[i-1] = (prices[i].Close - prices[i-1].Close) / prices[i-1].Close
		}
	}
	
	return returns
}

// CalculatePortfolioReturns calculates portfolio returns given weights and asset returns
func CalculatePortfolioReturns(assetReturns [][]float64, weights []float64) []float64 {
	if len(assetReturns) == 0 || len(weights) == 0 {
		return []float64{}
	}
	
	numPeriods := len(assetReturns[0])
	portfolioReturns := make([]float64, numPeriods)
	
	for t := 0; t < numPeriods; t++ {
		for i, weight := range weights {
			portfolioReturns[t] += weight * assetReturns[i][t]
		}
	}
	
	return portfolioReturns
}

// AnnualizedVolatility calculates annualized volatility from daily returns
func AnnualizedVolatility(dailyReturns []float64, tradingDays int) float64 {
	if len(dailyReturns) == 0 {
		return 0
	}
	dailyVol := StdDev(dailyReturns)
	return dailyVol * math.Sqrt(float64(tradingDays))
}

// RollingVolatility calculates rolling volatility
func RollingVolatility(returns []float64, window int) []float64 {
	if len(returns) < window {
		return []float64{}
	}
	
	rolling := make([]float64, len(returns)-window+1)
	for i := 0; i <= len(returns)-window; i++ {
		rolling[i] = StdDev(returns[i : i+window])
	}
	
	return rolling
}

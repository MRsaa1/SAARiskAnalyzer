package math

import (
	"time"
)

type StressScenarioResult struct {
	Name        string
	DeltaNAV    float64
	DeltaVaR    float64
	AssetImpact map[string]float64
}

type StressTestResult struct {
	Scenarios []StressScenarioResult
}

func ApplyHistoricalStress(
	positions map[string]float64,
	prices map[string][]PricePoint,
	startDate, endDate time.Time,
) (*StressScenarioResult, error) {
	
	totalDelta := 0.0
	assetImpact := make(map[string]float64)
	
	for symbol, marketValue := range positions {
		priceHistory, ok := prices[symbol]
		if !ok {
			continue
		}
		
		var startPrice, endPrice float64
		for _, p := range priceHistory {
			if p.Date.Equal(startDate) || (p.Date.After(startDate) && startPrice == 0) {
				startPrice = p.Close
			}
			if p.Date.Equal(endDate) || (p.Date.Before(endDate) && endPrice == 0) {
				endPrice = p.Close
			}
		}
		
		if startPrice > 0 && endPrice > 0 {
			returnPct := (endPrice - startPrice) / startPrice
			impact := marketValue * returnPct
			assetImpact[symbol] = impact
			totalDelta += impact
		}
	}
	
	return &StressScenarioResult{
		Name:        "Historical",
		DeltaNAV:    totalDelta,
		DeltaVaR:    0,
		AssetImpact: assetImpact,
	}, nil
}

func ApplyCustomStress(
	positions map[string]float64,
	assetClasses map[string]string,
	shocks map[string]float64,
) (*StressScenarioResult, error) {
	
	totalDelta := 0.0
	assetImpact := make(map[string]float64)
	
	for symbol, marketValue := range positions {
		class, ok := assetClasses[symbol]
		if !ok {
			continue
		}
		
		shock, ok := shocks[class]
		if !ok {
			shock = 0
		}
		
		impact := marketValue * shock
		assetImpact[symbol] = impact
		totalDelta += impact
	}
	
	return &StressScenarioResult{
		Name:        "Custom",
		DeltaNAV:    totalDelta,
		DeltaVaR:    0,
		AssetImpact: assetImpact,
	}, nil
}

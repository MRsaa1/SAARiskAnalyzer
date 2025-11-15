package math

import (
	"fmt"
	"math"
)

// BacktestResult contains backtesting results
type BacktestResult struct {
	Exceedances  int
	KupiecLR     float64
	KupiecPValue float64
	ChristLR     float64
	ChristPValue float64
}

// BacktestVaR performs VaR backtesting
func BacktestVaR(
	portfolioReturns []float64,
	varEstimates []float64, // rolling VaR estimates
	confidence float64,
) (*BacktestResult, error) {
	
	if len(portfolioReturns) != len(varEstimates) {
		return nil, fmt.Errorf("returns and VaR estimates length mismatch")
	}
	
	// Count exceedances (when loss exceeds VaR)
	exceedances := 0
	violations := make([]bool, len(portfolioReturns))
	
	for i := 0; i < len(portfolioReturns); i++ {
		loss := -portfolioReturns[i]
		if loss > varEstimates[i] {
			exceedances++
			violations[i] = true
		}
	}
	
	// Kupiec POF (Proportion of Failures) test
	n := float64(len(portfolioReturns))
	x := float64(exceedances)
	p := 1 - confidence // expected failure rate
	
	// LR_uc = -2 * ln[(p^x * (1-p)^(n-x)) / ((x/n)^x * (1-x/n)^(n-x))]
	var kupiecLR float64
	if x > 0 && x < n {
		pHat := x / n
		kupiecLR = -2 * (x*math.Log(p) + (n-x)*math.Log(1-p) - 
			x*math.Log(pHat) - (n-x)*math.Log(1-pHat))
	}
	
	// P-value from chi-square distribution with 1 df
	kupiecPValue := chiSquarePValue(kupiecLR, 1)
	
	// Christoffersen independence test (simplified)
	// Count transitions: 00, 01, 10, 11
	n00, n01, n10, n11 := 0, 0, 0, 0
	for i := 1; i < len(violations); i++ {
		if !violations[i-1] && !violations[i] {
			n00++
		} else if !violations[i-1] && violations[i] {
			n01++
		} else if violations[i-1] && !violations[i] {
			n10++
		} else {
			n11++
		}
	}
	
	// Independence test statistic
	var christLR float64
	if n01 > 0 && n11 > 0 && n00 > 0 && n10 > 0 {
		p01 := float64(n01) / float64(n00+n01)
		p11 := float64(n11) / float64(n10+n11)
		p2 := float64(n01+n11) / float64(n00+n01+n10+n11)
		
		christLR = -2 * (
			float64(n00+n01)*math.Log(1-p2) + float64(n01+n11)*math.Log(p2) -
			float64(n00)*math.Log(1-p01) - float64(n01)*math.Log(p01) -
			float64(n10)*math.Log(1-p11) - float64(n11)*math.Log(p11))
	}
	
	christPValue := chiSquarePValue(christLR, 1)
	
	return &BacktestResult{
		Exceedances:  exceedances,
		KupiecLR:     kupiecLR,
		KupiecPValue: kupiecPValue,
		ChristLR:     christLR,
		ChristPValue: christPValue,
	}, nil
}

// chiSquarePValue approximates chi-square p-value
func chiSquarePValue(x, df float64) float64 {
	// Simplified approximation
	if x <= 0 {
		return 1.0
	}
	// Using normal approximation for large x
	if x > 30 {
		return 0.0
	}
	// Very rough approximation
	return math.Exp(-x / 2)
}

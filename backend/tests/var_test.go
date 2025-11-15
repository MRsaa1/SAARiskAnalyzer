package tests

import (
	"testing"
	
	riskmath "github.com/reserveone/saa-risk-analyzer/internal/math"
)

func TestHistoricalVaR(t *testing.T) {
	// Sample returns data (daily)
	returns := []float64{
		-0.02, 0.01, 0.015, -0.01, 0.005,
		0.02, -0.015, 0.01, -0.025, 0.03,
		-0.01, 0.008, 0.012, -0.018, 0.022,
	}
	
	confidence := 0.95
	horizonDays := 1
	
	result, err := riskmath.CalculateHistoricalVaR(returns, confidence, horizonDays)
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if result.VaR <= 0 {
		t.Errorf("Expected positive VaR, got %f", result.VaR)
	}
	
	if result.Confidence != confidence {
		t.Errorf("Expected confidence %f, got %f", confidence, result.Confidence)
	}
	
	t.Logf("Historical VaR (95%%): %f", result.VaR)
}

func TestParametricVaR(t *testing.T) {
	// Sample returns data
	returns := []float64{
		0.01, -0.02, 0.015, -0.01, 0.02,
		-0.015, 0.01, 0.008, -0.012, 0.018,
	}
	
	confidence := 0.99
	horizonDays := 10
	
	result, err := riskmath.CalculateParametricVaR(returns, confidence, horizonDays)
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if result.VaR <= 0 {
		t.Errorf("Expected positive VaR, got %f", result.VaR)
	}
	
	if result.Method != "parametric_normal" {
		t.Errorf("Expected method parametric_normal, got %s", result.Method)
	}
	
	t.Logf("Parametric VaR (99%%, 10-day): %f", result.VaR)
}

func TestMeanStdDev(t *testing.T) {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	
	mean := riskmath.Mean(data)
	if mean != 3.0 {
		t.Errorf("Expected mean 3.0, got %f", mean)
	}
	
	stdDev := riskmath.StdDev(data)
	if stdDev < 1.4 || stdDev > 1.6 {
		t.Errorf("Expected stdDev around 1.5, got %f", stdDev)
	}
	
	t.Logf("Mean: %f, StdDev: %f", mean, stdDev)
}

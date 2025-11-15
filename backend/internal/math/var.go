package math

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

type VaRConfig struct {
	Confidence    float64
	HorizonDays   int
	Method        string
	Simulations   int
	UseLogReturns bool
	StudentDF     float64
}

type VaRResult struct {
	VaR           float64
	Method        string
	Confidence    float64
	HorizonDays   int
	Distribution  []float64
}

func CalculateHistoricalVaR(portfolioReturns []float64, confidence float64, horizonDays int) (*VaRResult, error) {
	if len(portfolioReturns) == 0 {
		return nil, fmt.Errorf("no returns data")
	}
	
	scaledReturns := make([]float64, len(portfolioReturns))
	scaleFactor := math.Sqrt(float64(horizonDays))
	for i, r := range portfolioReturns {
		scaledReturns[i] = r * scaleFactor
	}
	
	alpha := 1 - confidence
	varValue := -Quantile(scaledReturns, alpha)
	
	return &VaRResult{
		VaR:          varValue,
		Method:       "historical",
		Confidence:   confidence,
		HorizonDays:  horizonDays,
		Distribution: scaledReturns,
	}, nil
}

func CalculateParametricVaR(portfolioReturns []float64, confidence float64, horizonDays int) (*VaRResult, error) {
	if len(portfolioReturns) == 0 {
		return nil, fmt.Errorf("no returns data")
	}
	
	mu := Mean(portfolioReturns)
	sigma := StdDev(portfolioReturns)
	
	muScaled := mu * float64(horizonDays)
	sigmaScaled := sigma * math.Sqrt(float64(horizonDays))
	
	normal := distuv.Normal{Mu: 0, Sigma: 1}
	zScore := normal.Quantile(1 - confidence)
	
	varValue := -(muScaled + zScore*sigmaScaled)
	
	return &VaRResult{
		VaR:         varValue,
		Method:      "parametric_normal",
		Confidence:  confidence,
		HorizonDays: horizonDays,
	}, nil
}

func CalculateMonteCarloVaR(
	assetReturns [][]float64,
	weights []float64,
	confidence float64,
	horizonDays int,
	simulations int,
) (*VaRResult, error) {
	if len(assetReturns) == 0 || len(weights) == 0 {
		return nil, fmt.Errorf("invalid input")
	}
	
	numAssets := len(assetReturns)
	numPeriods := len(assetReturns[0])
	
	returnsMatrix := mat.NewDense(numPeriods, numAssets, nil)
	for i := 0; i < numAssets; i++ {
		for j := 0; j < numPeriods; j++ {
			returnsMatrix.Set(j, i, assetReturns[i][j])
		}
	}
	
	cov := Covariance(returnsMatrix)
	means := make([]float64, numAssets)
	for i := 0; i < numAssets; i++ {
		means[i] = Mean(assetReturns[i])
	}
	
	chol, err := Cholesky(cov)
	if err != nil {
		return nil, err
	}
	
	// Get L matrix directly
	L := mat.NewTriDense(numAssets, mat.Lower, nil)
	for i := 0; i < numAssets; i++ {
		for j := 0; j <= i; j++ {
			L.SetTri(i, j, chol.At(i, j))
		}
	}
	
	simulatedReturns := make([]float64, simulations)
	normal := distuv.Normal{Mu: 0, Sigma: 1}
	
	for sim := 0; sim < simulations; sim++ {
		z := mat.NewVecDense(numAssets, nil)
		for i := 0; i < numAssets; i++ {
			z.SetVec(i, normal.Rand())
		}
		
		var sampledReturns mat.VecDense
		sampledReturns.MulVec(L, z)
		
		portfolioReturn := 0.0
		for i := 0; i < numAssets; i++ {
			assetReturn := means[i] + sampledReturns.AtVec(i)
			portfolioReturn += weights[i] * assetReturn
		}
		
		scaleFactor := math.Sqrt(float64(horizonDays))
		simulatedReturns[sim] = portfolioReturn * scaleFactor
	}
	
	sort.Float64s(simulatedReturns)
	alpha := 1 - confidence
	varValue := -Quantile(simulatedReturns, alpha)
	
	return &VaRResult{
		VaR:          varValue,
		Method:       "monte_carlo",
		Confidence:   confidence,
		HorizonDays:  horizonDays,
		Distribution: simulatedReturns,
	}, nil
}

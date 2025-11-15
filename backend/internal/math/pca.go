package math

import (
	"fmt"
	"sort"
	
	"gonum.org/v1/gonum/mat"
)

// PCAResult contains PCA analysis results
type PCAResult struct {
	ExplainedVariance     []float64
	CumulativeVariance    []float64
	Components            *mat.Dense
	NumComponents         int
}

// CalculatePCA performs Principal Component Analysis
func CalculatePCA(assetReturns [][]float64, numComponents int) (*PCAResult, error) {
	if len(assetReturns) == 0 {
		return nil, fmt.Errorf("no returns data")
	}
	
	numAssets := len(assetReturns)
	numPeriods := len(assetReturns[0])
	
	if numComponents > numAssets {
		numComponents = numAssets
	}
	
	// Build matrix
	returnsMatrix := mat.NewDense(numPeriods, numAssets, nil)
	for i := 0; i < numAssets; i++ {
		for j := 0; j < numPeriods; j++ {
			returnsMatrix.Set(j, i, assetReturns[i][j])
		}
	}
	
	// Calculate covariance matrix
	cov := Covariance(returnsMatrix)
	
	// Eigenvalue decomposition
	var eigen mat.EigenSym
	ok := eigen.Factorize(cov, true)
	if !ok {
		return nil, fmt.Errorf("eigenvalue decomposition failed")
	}
	
	// Get eigenvalues and eigenvectors
	values := eigen.Values(nil)
	var vectors mat.Dense
	eigen.VectorsTo(&vectors)
	
	// Sort by eigenvalues (descending)
	type eigenPair struct {
		value  float64
		vector []float64
	}
	
	pairs := make([]eigenPair, len(values))
	for i := range values {
		vec := make([]float64, numAssets)
		mat.Col(vec, i, &vectors)
		pairs[i] = eigenPair{value: values[i], vector: vec}
	}
	
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].value > pairs[j].value
	})
	
	// Calculate explained variance
	totalVariance := 0.0
	for _, val := range values {
		totalVariance += val
	}
	
	explainedVar := make([]float64, numComponents)
	cumulativeVar := make([]float64, numComponents)
	cumulative := 0.0
	
	for i := 0; i < numComponents; i++ {
		explainedVar[i] = pairs[i].value / totalVariance
		cumulative += explainedVar[i]
		cumulativeVar[i] = cumulative
	}
	
	// Extract top components
	components := mat.NewDense(numComponents, numAssets, nil)
	for i := 0; i < numComponents; i++ {
		for j := 0; j < numAssets; j++ {
			components.Set(i, j, pairs[i].vector[j])
		}
	}
	
	return &PCAResult{
		ExplainedVariance:  explainedVar,
		CumulativeVariance: cumulativeVar,
		Components:         components,
		NumComponents:      numComponents,
	}, nil
}

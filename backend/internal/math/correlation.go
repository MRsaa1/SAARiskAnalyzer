package math

import (
	"fmt"
	
	"gonum.org/v1/gonum/mat"
)

type CorrelationResult struct {
	Matrix  *mat.SymDense
	Symbols []string
}

func CalculateCorrelationMatrix(assetReturns [][]float64, symbols []string) (*CorrelationResult, error) {
	if len(assetReturns) == 0 {
		return nil, fmt.Errorf("no returns data")
	}
	
	numAssets := len(assetReturns)
	numPeriods := len(assetReturns[0])
	
	returnsMatrix := mat.NewDense(numPeriods, numAssets, nil)
	for i := 0; i < numAssets; i++ {
		for j := 0; j < numPeriods; j++ {
			returnsMatrix.Set(j, i, assetReturns[i][j])
		}
	}
	
	corrMatrix := Correlation(returnsMatrix)
	
	return &CorrelationResult{
		Matrix:  corrMatrix,
		Symbols: symbols,
	}, nil
}

func ExportCorrelationMatrix(corr *mat.SymDense) [][]float64 {
	n, _ := corr.Dims()
	result := make([][]float64, n)
	for i := 0; i < n; i++ {
		result[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			result[i][j] = corr.At(i, j)
		}
	}
	return result
}

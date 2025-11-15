package math

import (
	"errors"
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

func Quantile(data []float64, q float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sorted := make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)
	
	return stat.Quantile(q, stat.Empirical, sorted, nil)
}

func Mean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	return stat.Mean(data, nil)
}

func StdDev(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	return stat.StdDev(data, nil)
}

func Covariance(data *mat.Dense) *mat.SymDense {
	r, c := data.Dims()
	cov := mat.NewSymDense(c, nil)
	
	for i := 0; i < c; i++ {
		for j := i; j < c; j++ {
			col1 := make([]float64, r)
			col2 := make([]float64, r)
			mat.Col(col1, i, data)
			mat.Col(col2, j, data)
			covVal := stat.Covariance(col1, col2, nil)
			cov.SetSym(i, j, covVal)
		}
	}
	
	return cov
}

func Correlation(data *mat.Dense) *mat.SymDense {
	r, c := data.Dims()
	corr := mat.NewSymDense(c, nil)
	
	for i := 0; i < c; i++ {
		for j := i; j < c; j++ {
			col1 := make([]float64, r)
			col2 := make([]float64, r)
			mat.Col(col1, i, data)
			mat.Col(col2, j, data)
			
			if i == j {
				corr.SetSym(i, j, 1.0)
			} else {
				corrVal := stat.Correlation(col1, col2, nil)
				corr.SetSym(i, j, corrVal)
			}
		}
	}
	
	return corr
}

func Cholesky(cov *mat.SymDense) (*mat.Cholesky, error) {
	var chol mat.Cholesky
	if ok := chol.Factorize(cov); !ok {
		return nil, errors.New("cholesky factorization failed")
	}
	return &chol, nil
}

func PortfolioVariance(weights []float64, cov *mat.SymDense) float64 {
	w := mat.NewVecDense(len(weights), weights)
	var result mat.VecDense
	result.MulVec(cov, w)
	return mat.Dot(w, &result)
}

func PortfolioStdDev(weights []float64, cov *mat.SymDense) float64 {
	return math.Sqrt(PortfolioVariance(weights, cov))
}

var (
	ErrCholeskyFailed = errors.New("cholesky failed")
	ErrInvalidInput   = errors.New("invalid input")
)

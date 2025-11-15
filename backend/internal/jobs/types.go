package jobs

import (
	"github.com/google/uuid"
)

const (
	StatusQueued    = "queued"
	StatusRunning   = "running"
	StatusSucceeded = "succeeded"
	StatusFailed    = "failed"
)

const (
	TypeVaR               = "var"
	TypeCVaR              = "cvar"
	TypeCorrelation       = "correlation"
	TypePCA               = "pca"
	TypeStress            = "stress"
	TypeBacktest          = "backtest"
	TypeRiskContribution  = "risk_contribution"
)

type JobFunc func(jobID uuid.UUID, progress chan<- int) (map[string]interface{}, error)

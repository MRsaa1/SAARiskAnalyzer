package domain

import (
	"github.com/google/uuid"
	"time"
)

// Auth DTOs
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

// Portfolio DTOs
type CreatePortfolioRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type ImportPositionsRequest struct {
	Positions []PositionImport `json:"positions" binding:"required"`
}

type PositionImport struct {
	Symbol   string  `json:"symbol" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required"`
	AvgPrice float64 `json:"avg_price" binding:"required"`
}

// Risk calculation DTOs
type VaRRequest struct {
	PortfolioID    uuid.UUID `json:"portfolio_id" binding:"required"`
	HorizonDays    int       `json:"horizon_days" binding:"required,min=1"`
	Confidence     float64   `json:"confidence" binding:"required,min=0,max=1"`
	Method         string    `json:"method" binding:"required"` // historical, parametric_normal, parametric_student, monte_carlo
	WindowDays     int       `json:"window_days"`
	Simulations    int       `json:"simulations"`
	UseLogReturns  bool      `json:"use_log_returns"`
}

type CVaRRequest struct {
	PortfolioID    uuid.UUID `json:"portfolio_id" binding:"required"`
	HorizonDays    int       `json:"horizon_days" binding:"required,min=1"`
	Confidence     float64   `json:"confidence" binding:"required,min=0,max=1"`
	Method         string    `json:"method" binding:"required"`
	WindowDays     int       `json:"window_days"`
	Simulations    int       `json:"simulations"`
	UseLogReturns  bool      `json:"use_log_returns"`
}

type CorrelationRequest struct {
	Symbols    []string `json:"symbols" binding:"required"`
	WindowDays int      `json:"window_days" binding:"required,min=10"`
}

type PCARequest struct {
	Symbols    []string `json:"symbols" binding:"required"`
	Components int      `json:"components" binding:"required,min=1"`
	WindowDays int      `json:"window_days" binding:"required,min=10"`
}

type StressTestRequest struct {
	PortfolioID       uuid.UUID        `json:"portfolio_id" binding:"required"`
	Scenarios         []StressScenario `json:"scenarios" binding:"required"`
	CorrelationRegime string           `json:"correlation_regime"` // tight, loose, current
}

type StressScenario struct {
	Name   string                 `json:"name" binding:"required"`
	Type   string                 `json:"type" binding:"required"` // historical, custom
	Window *TimeWindow            `json:"window,omitempty"`
	Shocks map[string]float64     `json:"shocks,omitempty"`
}

type TimeWindow struct {
	From string `json:"from" binding:"required"` // YYYY-MM-DD
	To   string `json:"to" binding:"required"`   // YYYY-MM-DD
}

type BacktestVaRRequest struct {
	PortfolioID uuid.UUID `json:"portfolio_id" binding:"required"`
	Confidence  float64   `json:"confidence" binding:"required,min=0,max=1"`
	WindowDays  int       `json:"window_days" binding:"required,min=10"`
	Method      string    `json:"method" binding:"required"` // historical, parametric
}

type RiskContributionRequest struct {
	PortfolioID uuid.UUID `json:"portfolio_id" binding:"required"`
	Confidence  float64   `json:"confidence" binding:"required,min=0,max=1"`
	WindowDays  int       `json:"window_days" binding:"required,min=10"`
}

// Response DTOs
type VaRResponse struct {
	JobID uuid.UUID `json:"job_id"`
	VaR   float64   `json:"var,omitempty"`
}

type CVaRResponse struct {
	JobID uuid.UUID `json:"job_id"`
	CVaR  float64   `json:"cvar,omitempty"`
}

type CorrelationResponse struct {
	JobID  uuid.UUID   `json:"job_id"`
	Matrix [][]float64 `json:"matrix,omitempty"`
}

type PCAResponse struct {
	JobID              uuid.UUID   `json:"job_id"`
	ExplainedVariance  []float64   `json:"explained_variance,omitempty"`
	Components         [][]float64 `json:"components,omitempty"`
}

type StressTestResponse struct {
	JobID    uuid.UUID              `json:"job_id"`
	Scenarios []ScenarioResult      `json:"scenarios,omitempty"`
}

type ScenarioResult struct {
	Name        string             `json:"name"`
	DeltaNAV    float64            `json:"delta_nav"`
	DeltaVaR    float64            `json:"delta_var"`
	AssetImpact []AssetImpact      `json:"asset_impact"`
}

type AssetImpact struct {
	Symbol string  `json:"symbol"`
	Impact float64 `json:"impact"`
}

type BacktestResult struct {
	JobID        uuid.UUID `json:"job_id"`
	Exceedances  int       `json:"exceedances,omitempty"`
	KupiecLR     float64   `json:"kupiec_lr,omitempty"`
	KupiecPValue float64   `json:"kupiec_p_value,omitempty"`
	ChristLR     float64   `json:"christ_lr,omitempty"`
	ChristPValue float64   `json:"christ_p_value,omitempty"`
}

type RiskContributionResponse struct {
	JobID         uuid.UUID            `json:"job_id"`
	Contributions []AssetContribution  `json:"contributions,omitempty"`
}

type AssetContribution struct {
	Symbol      string  `json:"symbol"`
	Component   float64 `json:"component_var"`
	Marginal    float64 `json:"marginal_var"`
	Percentage  float64 `json:"percentage"`
}

// Job response
type JobResponse struct {
	ID        uuid.UUID              `json:"id"`
	Type      string                 `json:"type"`
	Status    string                 `json:"status"`
	Progress  int                    `json:"progress"`
	Result    map[string]interface{} `json:"result,omitempty"`
	Error     string                 `json:"error,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// Health check
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

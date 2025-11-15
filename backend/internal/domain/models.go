package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a system user
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"not null;default:'viewer'" json:"role"` // admin, analyst, viewer
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// Asset represents a financial instrument
type Asset struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Symbol    string    `gorm:"uniqueIndex;not null" json:"symbol"`
	Name      string    `gorm:"not null" json:"name"`
	Class     string    `gorm:"not null" json:"class"` // Equity, Bond, FX, Commodities, Crypto
	Currency  string    `gorm:"not null;default:'USD'" json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Asset) TableName() string {
	return "assets"
}

// Price represents historical price data
type Price struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AssetID   uuid.UUID `gorm:"type:uuid;not null;index" json:"asset_id"`
	Asset     Asset     `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	Date      time.Time `gorm:"not null;index" json:"date"`
	Close     float64   `gorm:"not null" json:"close"`
	CreatedAt time.Time `json:"created_at"`
}

func (Price) TableName() string {
	return "prices"
}

// Portfolio represents a collection of positions
type Portfolio struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `json:"description"`
	Positions   []Position `gorm:"foreignKey:PortfolioID" json:"positions,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Portfolio) TableName() string {
	return "portfolios"
}

// Position represents a holding in a portfolio
type Position struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PortfolioID uuid.UUID `gorm:"type:uuid;not null;index" json:"portfolio_id"`
	AssetID     uuid.UUID `gorm:"type:uuid;not null;index" json:"asset_id"`
	Asset       Asset     `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	Quantity    float64   `gorm:"not null" json:"quantity"`
	AvgPrice    float64   `gorm:"not null" json:"avg_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Position) TableName() string {
	return "positions"
}

// Scenario represents a stress test scenario
type Scenario struct {
	ID          uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string                 `gorm:"not null" json:"name"`
	Description string                 `json:"description"`
	Kind        string                 `gorm:"not null" json:"kind"` // historical, custom, regime
	Payload     map[string]interface{} `gorm:"type:jsonb" json:"payload"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

func (Scenario) TableName() string {
	return "scenarios"
}

// Job represents an async computation job
type Job struct {
	ID        uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Type      string                 `gorm:"not null;index" json:"type"` // var, cvar, stress, pca, etc.
	Status    string                 `gorm:"not null;index" json:"status"` // queued, running, succeeded, failed
	Progress  int                    `gorm:"default:0" json:"progress"` // 0-100
	Result    map[string]interface{} `gorm:"type:jsonb" json:"result,omitempty"`
	Error     string                 `json:"error,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func (Job) TableName() string {
	return "jobs"
}

// BeforeCreate hooks to ensure UUIDs
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (a *Asset) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (p *Price) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (p *Portfolio) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (p *Position) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (s *Scenario) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}

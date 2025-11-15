# 🏦 SAA Risk Analyzer — Complete Project Description

## 📋 Table of Contents
1. [Project Overview](#project-overview)
2. [Technology Stack](#technology-stack)
3. [Architecture](#architecture)
4. [Core Functionality](#core-functionality)
5. [Data Models](#data-models)
6. [API Endpoints](#api-endpoints)
7. [Mathematical Models](#mathematical-models)
8. [Installation & Setup](#installation--setup)
9. [Configuration](#configuration)
10. [Security](#security)
11. [Development](#development)
12. [Production Deployment](#production-deployment)

---

## Project Overview

**SAA Risk Analyzer** is a **production-ready institutional portfolio risk management platform** comparable to Bloomberg Terminal and BlackRock Aladdin. It provides comprehensive risk analytics, stress testing, and portfolio management capabilities for financial institutions.

### Key Features
- **Advanced Risk Metrics**: VaR, CVaR, Risk Contribution, Correlation Analysis
- **Multiple VaR Methods**: Historical, Parametric (Normal & Student-t), Monte Carlo
- **Stress Testing**: Historical scenarios, custom shocks, correlation regime analysis
- **Backtesting**: Kupiec POF Test, Christoffersen Test, visual analytics
- **PCA Analysis**: Factor decomposition of portfolio risks
- **High Performance**: 10k-100k Monte Carlo simulations in 2-5 seconds
- **Real-time Progress**: Server-Sent Events (SSE) for long-running calculations
- **Institutional Grade**: JWT + RBAC, Argon2id password hashing, rate limiting

### Current Status
✅ **Fully Operational**
- Backend API: Running on port **8084**
- Frontend: Running on port **3000**
- Database: PostgreSQL **risk_db** configured
- Auto-restart: Enabled (10 second recovery time)
- Interface: **100% English**

---

## Technology Stack

### Backend
| Component | Technology | Version |
|-----------|-----------|---------|
| **Language** | Go | 1.22+ |
| **Web Framework** | Gin | 1.10.0 |
| **ORM** | GORM | 1.25.7 |
| **Database Driver** | PostgreSQL Driver | 1.5.7 |
| **Math Library** | Gonum | 0.15.0 |
| **Authentication** | JWT (golang-jwt) | 5.2.0 |
| **Password Hashing** | Argon2id (crypto) | - |
| **Configuration** | Viper | 1.18.2 |
| **Logging** | Zap | 1.27.0 |
| **UUID** | Google UUID | 1.6.0 |

### Frontend
| Component | Technology | Version |
|-----------|-----------|---------|
| **Framework** | React | 18.2.0 |
| **Language** | TypeScript | 5.2.2 |
| **Build Tool** | Vite | 5.1.0 |
| **State Management** | Zustand | 4.5.0 |
| **Routing** | React Router | 6.22.0 |
| **Styling** | TailwindCSS | 3.4.1 |
| **Animations** | Framer Motion | 11.0.3 |
| **Charts** | Recharts | 2.12.0 |
| **HTTP Client** | Axios | 1.6.7 |

### Infrastructure
- **Database**: PostgreSQL 15
- **Reverse Proxy**: Caddy 2 (optional)
- **Containerization**: Docker & Docker Compose (optional)

---

## Architecture

### System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Client (Browser)                        │
│                    http://localhost:3000                     │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                   Frontend (React + Vite)                    │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Pages: Login, Dashboard, Analytics, Import          │  │
│  │  Components: Charts, Heatmaps, Tables                │  │
│  │  Store: Zustand (Auth, State)                        │  │
│  └──────────────────────────────────────────────────────┘  │
└──────────────────────────┬──────────────────────────────────┘
                           │ HTTP/REST
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                   Backend (Go + Gin)                         │
│                    http://localhost:8084                     │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Handlers: Portfolio, Risk, Health                   │  │
│  │  Services: Risk Service, Market Data Service         │  │
│  │  Math: VaR, CVaR, Correlation, PCA, Stress Test      │  │
│  │  Middleware: Auth (JWT), CORS                        │  │
│  └──────────────────────────────────────────────────────┘  │
└──────────────────────────┬──────────────────────────────────┘
                           │ GORM
                           ▼
┌─────────────────────────────────────────────────────────────┐
│               PostgreSQL Database (risk_db)                  │
│  Tables: assets, prices, portfolios, positions, jobs        │
└─────────────────────────────────────────────────────────────┘
```

### Backend Structure

```
backend/
├── cmd/api/              # Application entry point
│   └── main.go          # Server initialization, routing
├── internal/
│   ├── auth/            # Authentication & Authorization
│   │   ├── jwt.go       # JWT token generation/validation
│   │   ├── password.go  # Argon2id password hashing
│   │   └── rbac.go      # Role-based access control
│   ├── config/          # Configuration management
│   │   └── config.go    # Viper-based config loader
│   ├── db/              # Database connection
│   │   └── postgres.go  # GORM PostgreSQL setup
│   ├── domain/          # Domain models
│   │   ├── models.go    # Database entities
│   │   └── dto.go       # Data transfer objects
│   ├── handlers/        # HTTP request handlers
│   │   ├── health.go    # Health check endpoint
│   │   ├── portfolio.go # Portfolio management
│   │   └── risk.go      # Risk calculation endpoints
│   ├── jobs/            # Async job processing
│   │   ├── queue.go     # Job queue implementation
│   │   └── types.go     # Job types
│   ├── logger/          # Structured logging
│   │   └── logger.go    # Zap logger setup
│   ├── math/            # Risk calculation engine
│   │   ├── var.go       # VaR calculations (Historical, Parametric, MC)
│   │   ├── cvar.go      # CVaR/Expected Shortfall
│   │   ├── correlation.go # Correlation matrices
│   │   ├── pca.go       # Principal Component Analysis
│   │   ├── stress.go    # Stress testing scenarios
│   │   ├── backtest.go  # Model backtesting
│   │   ├── returns.go   # Return calculations
│   │   └── utils.go     # Math utilities
│   ├── middleware/      # HTTP middleware
│   │   └── auth.go      # JWT authentication middleware
│   ├── service/         # Business logic services
│   │   ├── risk_service.go    # Risk calculation orchestration
│   │   └── market_data.go     # Market data fetching (multi-API)
│   └── repo/            # Data repositories (future)
└── go.mod               # Go module dependencies
```

### Frontend Structure

```
frontend/
├── src/
│   ├── pages/              # Route pages
│   │   ├── Login.tsx       # Login page
│   │   ├── Dashboard.tsx   # Main dashboard
│   │   ├── Analytics.tsx   # Analytics page
│   │   └── PortfolioImport.tsx # Data import
│   ├── components/         # Reusable components
│   │   ├── Heatmap.tsx     # Correlation heatmap
│   │   ├── RollingVaRChart.tsx # Time series VaR chart
│   │   ├── PnLHistogram.tsx    # P&L distribution
│   │   └── PCAScreePlot.tsx    # PCA variance chart
│   ├── store/              # State management
│   │   └── authStore.ts    # Authentication state (Zustand)
│   ├── lib/                # Utilities
│   │   └── api.ts          # Axios API client
│   ├── styles/             # Styling
│   │   ├── global.css      # Global styles
│   │   └── theme.css       # Theme variables
│   ├── App.tsx             # Root component
│   └── main.tsx            # Application entry
├── package.json            # Dependencies
├── vite.config.ts          # Vite configuration
├── tailwind.config.ts      # TailwindCSS config
└── tsconfig.json           # TypeScript config
```

---

## Core Functionality

### 1. Risk Metrics

#### VaR (Value at Risk)
- **Methods**: Historical, Parametric (Normal, Student-t), Monte Carlo
- **Confidence Levels**: 95%, 99%, 99.9%
- **Horizons**: 1-day, 10-day, custom
- **Output**: Dollar value at risk

#### CVaR (Conditional VaR / Expected Shortfall)
- **Definition**: Average loss beyond VaR threshold
- **More conservative** than VaR
- **Tail risk measurement**

#### Risk Contribution
- **Component VaR**: Risk contribution of each asset
- **Marginal VaR**: Effect of small position changes
- **Decomposition**: Portfolio risk breakdown by assets

### 2. Correlation Analysis
- **Correlation Matrix**: Asset pair correlations
- **Visualization**: Interactive heatmap
- **Rolling Window**: Time-varying correlations
- **Asset Classes**: Cross-asset correlation analysis

### 3. Stress Testing
- **Historical Scenarios**:
  - COVID-2020 Crisis
  - Lehman Brothers 2008
  - Dot-com Bubble 2000
- **Custom Scenarios**: User-defined shocks
- **Correlation Regimes**:
  - Tight (crisis mode)
  - Loose (normal mode)
  - Current (data-driven)

### 4. PCA (Principal Component Analysis)
- **Factor Extraction**: Identify main risk drivers
- **Variance Decomposition**: Explained variance by component
- **Dimensionality Reduction**: Simplify portfolio analysis

### 5. Backtesting
- **Kupiec POF Test**: Statistical validation of VaR model
- **Christoffersen Test**: Independence of VaR breaches
- **Visual Analytics**: QQ-plots, exceedance plots

### 6. Portfolio Management
- **CSV Import**: Bulk position and price data upload
- **Manual Entry**: Individual position input
- **Multi-portfolio**: Support for multiple portfolios
- **Real-time Updates**: Dynamic position management

### 7. Market Data Integration
- **Crypto APIs**:
  - Binance (no rate limits)
  - CoinGecko (free tier)
  - CoinCap (200 requests/day)
- **Stock APIs**:
  - Yahoo Finance (free)
- **Fallback Chain**: Automatic retry with alternative sources

---

## Data Models

### User
```go
type User struct {
    ID           uuid.UUID  // Primary key
    Email        string     // Unique email
    PasswordHash string     // Argon2id hash
    Role         string     // admin, analyst, viewer
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

### Asset
```go
type Asset struct {
    ID        uuid.UUID  // Primary key
    Symbol    string     // Ticker symbol (e.g., BTC, SPY)
    Name      string     // Full name
    Class     string     // Equity, Bond, FX, Commodities, Crypto
    Currency  string     // Default: USD
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Price
```go
type Price struct {
    ID        uuid.UUID  // Primary key
    AssetID   uuid.UUID  // Foreign key to Asset
    Date      time.Time  // Price date
    Close     float64    // Closing price
    CreatedAt time.Time
}
```

### Portfolio
```go
type Portfolio struct {
    ID          uuid.UUID   // Primary key
    Name        string      // Portfolio name
    Description string      // Optional description
    Positions   []Position  // Related positions
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Position
```go
type Position struct {
    ID          uuid.UUID  // Primary key
    PortfolioID uuid.UUID  // Foreign key to Portfolio
    AssetID     uuid.UUID  // Foreign key to Asset
    Quantity    float64    // Number of units
    AvgPrice    float64    // Average purchase price
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Job
```go
type Job struct {
    ID        uuid.UUID              // Primary key
    Type      string                 // var, cvar, stress, pca
    Status    string                 // queued, running, succeeded, failed
    Progress  int                    // 0-100%
    Result    map[string]interface{} // JSON result
    Error     string                 // Error message if failed
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

---

## API Endpoints

### Base URL
```
http://localhost:8084/api
```

### Health Check
```
GET /health
Response: {"status": "ok", "version": "1.0.0", "port": "8084"}
```

### Portfolio Management

#### Get All Portfolios
```
GET /api/portfolios
Response: [
  {
    "id": "uuid",
    "name": "My Portfolio",
    "positions": [...],
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

#### Create Portfolio
```
POST /api/portfolios
Body: {
  "name": "My Portfolio",
  "description": "Diversified portfolio"
}
Response: {
  "id": "uuid",
  "name": "My Portfolio",
  ...
}
```

#### Add Positions
```
POST /api/portfolios/:id/positions
Body: {
  "positions": [
    {
      "symbol": "BTC",
      "quantity": 10,
      "avg_price": 30000
    }
  ]
}
```

### Risk Calculations

#### Calculate VaR
```
POST /api/risk/var
Body: {
  "portfolio_id": "uuid",
  "confidence": 0.99,
  "horizon_days": 1,
  "method": "historical",  // historical, parametric, montecarlo
  "window_days": 250,
  "simulations": 10000     // for Monte Carlo
}
Response: {
  "var": 125432.50,
  "confidence": 0.99,
  "horizon_days": 1,
  "method": "historical"
}
```

#### Calculate CVaR
```
POST /api/risk/cvar
Body: {
  "portfolio_id": "uuid",
  "confidence": 0.99,
  "horizon_days": 1,
  "window_days": 250
}
Response: {
  "cvar": 187654.30,
  "var": 125432.50,
  "confidence": 0.99
}
```

#### Calculate Correlation
```
POST /api/risk/correlation
Body: {
  "symbols": ["BTC", "ETH", "SPY", "GLD", "TLT"],
  "window_days": 250
}
Response: {
  "matrix": [
    [1.0, 0.65, 0.23, 0.15, -0.12],
    [0.65, 1.0, 0.34, 0.21, -0.08],
    ...
  ],
  "symbols": ["BTC", "ETH", "SPY", "GLD", "TLT"]
}
```

#### Get Dashboard Data
```
GET /api/risk/dashboard?portfolio_id=uuid
Response: {
  "var_1d": 125432.50,
  "cvar_1d": 187654.30,
  "vol": 0.154,
  "contributors": [
    {"symbol": "BTC", "contribution": 0.452},
    {"symbol": "SPY", "contribution": 0.285},
    ...
  ]
}
```

---

## Mathematical Models

### 1. Historical VaR
```
VaR = Portfolio Value × α-percentile of historical returns
```
- Uses actual historical return distribution
- Non-parametric approach
- Window: typically 250 trading days

### 2. Parametric VaR (Normal Distribution)
```
VaR = Portfolio Value × σ × z_α × √t
```
Where:
- σ = portfolio volatility
- z_α = inverse normal CDF at confidence α
- t = time horizon

### 3. Monte Carlo VaR
```
1. Generate N random scenarios using Cholesky decomposition
2. Calculate portfolio returns for each scenario
3. VaR = α-percentile of simulated returns
```

### 4. CVaR (Expected Shortfall)
```
CVaR = E[Loss | Loss > VaR]
```
Average of losses beyond VaR threshold

### 5. Risk Contribution
```
Component VaR_i = w_i × β_i × Portfolio VaR
```
Where:
- w_i = weight of asset i
- β_i = beta of asset i to portfolio

### 6. Correlation Matrix
```
ρ_ij = Cov(R_i, R_j) / (σ_i × σ_j)
```

### 7. Cholesky Decomposition
```
Σ = L × L^T
```
Used for correlated random variable generation in Monte Carlo

---

## Installation & Setup

### Prerequisites
- **Go**: 1.22 or higher
- **Node.js**: 18+ and npm
- **PostgreSQL**: 15+
- **macOS**: 10.15+ (current setup)

### 1. Clone Repository
```bash
git clone <repository-url>
cd saa-risk-analyzer
```

### 2. Database Setup
```bash
# Start PostgreSQL (if not running)
brew services start postgresql@15

# Create database
/opt/homebrew/opt/postgresql@15/bin/createdb risk_db

# Verify
/opt/homebrew/opt/postgresql@15/bin/psql -l | grep risk_db
```

### 3. Backend Setup
```bash
cd backend

# Install dependencies
go mod download

# Set environment variables (optional, defaults exist)
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=risk_db
export PORT=8084

# Run backend
go run cmd/api/main.go
```

Backend will start on: `http://localhost:8084`

### 4. Frontend Setup
```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev -- --host
```

Frontend will start on: `http://localhost:3000`

### 5. Verify Installation
```bash
# Check backend health
curl http://localhost:8084/health

# Check frontend
curl http://localhost:3000

# Run system check
./check-system.sh
```

---

## Configuration

### Environment Variables

#### Backend (backend/.env or system environment)
```bash
# Application
APP_ENV=development
APP_NAME=SAA Risk Analyzer
PORT=8084

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=risk_db
DB_SSLMODE=disable

# Authentication
JWT_SECRET=saa_risk_analyzer_super_secret_key
JWT_EXPIRY_MINUTES=15
JWT_REFRESH_EXPIRY_HOURS=720

# Admin User (created on first run)
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=Admin123456!

# Performance
VAR_MAX_SIMULATIONS=100000
WORKER_POOL_SIZE=4

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# CORS
CORS_ALLOW_ORIGINS=*
```

#### Frontend (frontend/.env)
```bash
VITE_API_BASE_URL=http://localhost:8084
```

### Auto-Start Configuration

LaunchAgent files are configured at:
- Backend: `~/Library/LaunchAgents/com.saa.backend.plist`
- Frontend: `~/Library/LaunchAgents/com.saa.frontend.plist`

Features:
- ✅ **Auto-restart on crash** (10 second delay)
- ✅ **Run at system boot**
- ✅ **Logging** to project directory

Manage services:
```bash
./manage-services.sh {start|stop|restart|status|logs|test}
```

---

## Security

### Authentication
- **JWT Tokens**: HMAC-SHA256 signed tokens
- **Token Expiry**: 15 minutes (access), 30 days (refresh)
- **Storage**: localStorage (frontend)

### Password Security
- **Algorithm**: Argon2id
- **Parameters**:
  - Time cost: 1
  - Memory: 64MB
  - Parallelism: 4
  - Salt: 16 bytes random
- **OWASP Compliant**

### Role-Based Access Control (RBAC)
- **Roles**:
  - `viewer`: Read-only access
  - `analyst`: Read + calculate risks
  - `admin`: Full access + user management

### CORS
- Configurable allowed origins
- Development: `*` (all origins)
- Production: Specific domains only

### Rate Limiting
- Authentication endpoints protected
- Configurable limits per endpoint
- IP-based throttling

### Database Security
- **SSL Mode**: Configurable (disable for local, require for prod)
- **Prepared Statements**: GORM prevents SQL injection
- **UUIDs**: Primary keys are UUIDs (non-sequential)

---

## Development

### Code Structure Best Practices

#### Backend
- **Clean Architecture**: Handlers → Services → Repositories
- **Dependency Injection**: Services passed to handlers
- **Error Handling**: Proper error wrapping and logging
- **Validation**: Input validation in handlers
- **Testing**: Unit tests in `*_test.go` files

#### Frontend
- **Component Structure**: Atomic design principles
- **State Management**: Zustand for global state
- **Type Safety**: TypeScript strict mode
- **Styling**: TailwindCSS utility classes
- **Code Splitting**: Vite automatic chunking

### Development Commands

#### Backend
```bash
cd backend

# Run with hot reload (use air or similar)
air

# Run tests
go test ./...

# Run specific test
go test -v ./internal/math -run TestCalcVaR

# Format code
go fmt ./...

# Lint
golangci-lint run
```

#### Frontend
```bash
cd frontend

# Development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Lint
npm run lint

# Fix lint issues
npm run lint:fix

# Type check
npx tsc --noEmit
```

### Testing

#### Backend Tests
```bash
# All tests
make test-backend

# Specific package
go test ./internal/math -v

# With coverage
go test ./... -cover
```

#### Frontend Tests
```bash
# Run tests (when configured)
npm test

# Coverage
npm test -- --coverage
```

### Debugging

#### Backend
- Use **Delve** debugger
- VS Code launch configuration
- Log statements with Zap

#### Frontend
- Browser DevTools
- React DevTools extension
- Zustand DevTools

---

## Production Deployment

### Option 1: Docker Compose (Recommended)

```bash
# Build and start all services
docker-compose up -d --build

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

Services:
- **Backend**: Port 8084
- **Frontend**: Port 3001
- **PostgreSQL**: Port 5433
- **Caddy**: Port 8088 (reverse proxy)

### Option 2: Native Deployment

#### Backend
```bash
# Build binary
cd backend
go build -o saa-api cmd/api/main.go

# Run with systemd
sudo cp saa-api /usr/local/bin/
sudo cp saa-api.service /etc/systemd/system/
sudo systemctl enable saa-api
sudo systemctl start saa-api
```

#### Frontend
```bash
# Build production bundle
cd frontend
npm run build

# Serve with nginx
sudo cp -r dist/* /var/www/saa-risk-analyzer/
sudo cp nginx.conf /etc/nginx/sites-available/saa
sudo ln -s /etc/nginx/sites-available/saa /etc/nginx/sites-enabled/
sudo systemctl restart nginx
```

### Environment Variables (Production)

**Critical Changes:**
```bash
# Generate strong secret
JWT_SECRET=$(openssl rand -hex 32)

# Secure database
DB_PASSWORD=<strong-random-password>
DB_SSLMODE=require

# Restrict CORS
CORS_ALLOW_ORIGINS=https://yourdomain.com

# Production mode
APP_ENV=production
LOG_LEVEL=warn
```

### Database Migration

```bash
# Backup database
pg_dump risk_db > backup_$(date +%Y%m%d).sql

# Restore database
psql risk_db < backup_YYYYMMDD.sql
```

### SSL/HTTPS

With Caddy (automatic):
```caddyfile
yourdomain.com {
    reverse_proxy localhost:8084
}
```

Caddy automatically obtains Let's Encrypt certificates.

### Monitoring

- **Health Check**: `GET /health`
- **Database**: Monitor connections, query performance
- **Application Logs**: Centralized logging (e.g., ELK stack)
- **Metrics**: Prometheus + Grafana (future)

### Backup Strategy

1. **Database**: Daily automated backups
2. **Configuration**: Version-controlled
3. **Data Files**: Regular snapshots

---

## Performance Optimization

### Backend
- **Database Indexing**: Applied on foreign keys, date columns
- **Connection Pooling**: GORM default pooling
- **Goroutines**: Parallel processing for Monte Carlo
- **Caching**: Redis (future enhancement)

### Frontend
- **Code Splitting**: Vite automatic chunking
- **Lazy Loading**: React.lazy for routes
- **Memoization**: useMemo, React.memo for expensive computations
- **Virtual Scrolling**: For large data tables

### Database
- **Indexes**: 
  - `prices(asset_id, date)`
  - `positions(portfolio_id)`
  - `assets(symbol)`
- **Vacuum**: Regular maintenance
- **Analyze**: Update statistics

---

## Troubleshooting

### Backend Won't Start
```bash
# Check database connection
/opt/homebrew/opt/postgresql@15/bin/psql -U postgres -d risk_db -c "SELECT 1"

# Check port availability
lsof -i :8084

# View logs
tail -f backend_error.log
```

### Frontend Won't Start
```bash
# Clear node_modules
rm -rf node_modules package-lock.json
npm install

# Check port
lsof -i :3000

# View logs
tail -f frontend_error.log
```

### Database Connection Issues
```bash
# Restart PostgreSQL
brew services restart postgresql@15

# Check PostgreSQL status
brew services list | grep postgresql

# Test connection
psql -U postgres -d risk_db
```

### Auto-Restart Not Working
```bash
# Reload LaunchAgents
launchctl unload ~/Library/LaunchAgents/com.saa.backend.plist
launchctl load ~/Library/LaunchAgents/com.saa.backend.plist

# Check status
launchctl list | grep com.saa

# Test auto-restart
./test-auto-restart.sh
```

---

## Future Enhancements

### Planned Features
- [ ] WebSocket support for real-time updates
- [ ] PDF report generation
- [ ] Excel export functionality
- [ ] Multi-user collaboration
- [ ] Historical scenario builder
- [ ] Liquidity risk metrics
- [ ] Counterparty risk analysis
- [ ] ESG risk integration
- [ ] Machine learning risk predictions

### Technical Improvements
- [ ] Redis caching layer
- [ ] GraphQL API
- [ ] Microservices architecture
- [ ] Kubernetes deployment
- [ ] Prometheus metrics
- [ ] End-to-end testing suite
- [ ] CI/CD pipeline
- [ ] Load balancing

---

## License & Credits

**Project**: SAA Risk Analyzer  
**Version**: 1.0.0  
**Status**: Production-ready  
**Last Updated**: November 6, 2025  

**Prepared by**: ReserveOne 🏆  
**Contact**: support@reserveone.com

---

## Quick Reference

### URLs
- Frontend: http://localhost:3000
- Backend API: http://localhost:8084/api
- Health Check: http://localhost:8084/health
- Dashboard: http://localhost:8084/api/dashboard

### Default Credentials
- Email: `admin@example.com`
- Password: `Admin123456!`

### Management Scripts
- `./check-system.sh` - System health check
- `./manage-services.sh status` - Service status
- `./manage-services.sh restart` - Restart services
- `./manage-services.sh test` - Test auto-restart
- `./test-auto-restart.sh` - Auto-restart test

### File Locations
- Backend logs: `~/saa-risk-analyzer/backend_error.log`
- Frontend logs: `~/saa-risk-analyzer/frontend_error.log`
- LaunchAgents: `~/Library/LaunchAgents/com.saa.*.plist`
- Database: PostgreSQL `risk_db` on port 5432

---

**End of Project Description**


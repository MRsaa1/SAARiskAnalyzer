# SAA Risk Analyzer

**Institutional-Grade Portfolio Risk Management Platform**

[![Live Demo](https://img.shields.io/badge/Demo-Live-brightgreen)](http://104.248.70.69/risk-analyzer)
[![Tech Stack](https://img.shields.io/badge/Stack-Go%20%2B%20React-blue)](./backend/go.mod)
[![License](https://img.shields.io/badge/License-Proprietary-red)](LICENSE)

---

## 📋 Overview

**SAA Risk Analyzer** is a production-ready portfolio risk analysis platform at the level of Bloomberg Terminal / BlackRock Aladdin. The system provides comprehensive risk metrics, stress testing, and backtesting capabilities for institutional investors.

### Key Features

- **Risk Calculation**
  - **VaR (Value at Risk)**: Historical, Parametric (Normal & Student-t), Monte Carlo
  - **CVaR (Conditional VaR / Expected Shortfall)**: Average in tail distribution
  - **Risk Contribution**: Portfolio risk decomposition by assets (Component & Marginal VaR)
  - **Correlation Analysis**: Correlation matrices with heatmap visualization
  - **PCA (Principal Component Analysis)**: Factor decomposition of risks

- **Stress Testing**
  - Historical scenarios: COVID-2020, Lehman 2008, Tech Bubble 2000
  - Custom scenarios: User-defined shocks by asset classes
  - Correlation modes: tight (crisis), loose (normal), current
  - Impact Analysis: Influence on adjacent assets and industries

- **Backtesting**
  - Kupiec POF Test: Statistical validation of VaR model
  - Christoffersen Test: Independence of violations check
  - Visual Analytics: QQ-plots, exceedance plots

- **Performance**
  - 10k–100k Monte Carlo VaR simulations in 2–5 seconds
  - Streaming progress indication via Server-Sent Events (SSE)
  - Asynchronous processing of long calculations via job queue

---

## 🏗️ Architecture

### Backend (Go + Gin)

```
backend/
├── cmd/api/              # Entry point
├── internal/
│   ├── math/             # Risk engine (VaR, CVaR, PCA, Stress)
│   ├── handlers/         # HTTP handlers
│   ├── service/          # Business services
│   ├── repo/             # Data access layer
│   └── auth/             # JWT, RBAC
└── api/                  # OpenAPI spec
```

**Technologies:**
- Go 1.22+ (Gin, GORM, Gonum)
- PostgreSQL 15
- JWT + RBAC (Argon2id for passwords)
- OpenAPI 3.1 + Swagger UI
- Zap (structured logging)

### Frontend (React + TypeScript)

```
frontend/
└── src/
    ├── pages/            # Application pages
    ├── components/       # UI components
    └── store/           # Zustand stores
```

**Technologies:**
- React 18 + TypeScript
- Vite (fast build)
- Zustand (state management)
- TailwindCSS + Framer Motion
- Recharts (visualizations)

---

## 🚀 Live Demo

**Production URL:** [http://104.248.70.69/risk-analyzer](http://104.248.70.69/risk-analyzer)

### Features Demonstrated

1. **Dashboard**
   - Overview of VaR, CVaR, correlations
   - Top risk contributors
   - Real-time metrics

2. **Portfolio Management**
   - Position management
   - CSV import
   - Portfolio composition

3. **Risk Lab**
   - VaR/CVaR calculations with different methods
   - Monte Carlo simulations
   - Historical analysis

4. **Stress Testing**
   - Scenario constructor
   - Impact analysis
   - Custom scenarios

5. **Factor Analysis (PCA)**
   - Factor decomposition of risks
   - Principal component visualization

---

## 📸 Screenshots

![Dashboard](./screenshots/dashboard.png)
![Risk Lab](./screenshots/risk-lab.png)
![Stress Testing](./screenshots/stress-test.png)

---

## 🛠️ Installation & Setup

### Prerequisites

- Docker 24+ and Docker Compose
- Make (optional, for convenience)
- 4GB+ RAM

### Quick Start

```bash
# Clone repository
git clone <repository-url>
cd saa-risk-analyzer

# Copy and configure environment variables
cp .env.example .env
# Edit .env, at minimum change JWT_SECRET!

# Build and start all services
make build
make up

# Or one command
docker-compose up -d --build

# Load demo data (wait ~30 sec for DB and API to start)
make seed
```

### Access

- **Frontend:** http://localhost
- **API Health:** http://localhost/health
- **Swagger UI:** http://localhost/swagger/index.html
- **API Direct:** http://localhost:8083

**Default Login:**
- Email: `admin@example.com`
- Password: `Admin123456!`

---

## 📊 Demo Data

In `data/` directory:
- **prices.csv**: 5 assets (SPY, TLT, GLD, BTC, EUR), 3 years of history
- **positions.csv**: Demo portfolio ($1M)

---

## 🧪 Testing

```bash
# All tests
make test

# Backend only
make test-backend

# Frontend only
make test-frontend

# Linters
make lint
```

---

## 📖 API Documentation

After startup, interactive Swagger UI documentation is available: http://localhost/swagger/index.html

### Example Requests

#### Authentication
```bash
curl -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"Admin123456!"}'
```

#### Calculate VaR (Historical)
```bash
curl -X POST http://localhost/api/risk/var \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "portfolio_id": "uuid-here",
    "horizon_days": 1,
    "confidence": 0.99,
    "method": "historical",
    "window_days": 250
  }'
```

---

## 🎨 UI/UX

The platform uses a dark theme matching your website (ReserveOne):
- Color palette: black background, gold accent (#caa76a)
- Font: Inter
- Animations: Framer Motion

### Main Pages

1. **Dashboard**: Overview of VaR, CVaR, correlations, top risk contributors
2. **Portfolio**: Position management, CSV import
3. **Risk Lab**: Run VaR/CVaR calculations with different methods
4. **Stress Testing**: Scenario constructor, impact analysis
5. **Factors (PCA)**: Factor decomposition of risks
6. **Reports**: Generate PDF/CSV reports
7. **Admin**: User and role management

---

## 🔒 Security

- **Passwords**: Argon2id (OWASP standard)
- **JWT**: HMAC-SHA256, short TTL (15 min) + refresh
- **RBAC**: Three roles (viewer, analyst, admin)
- **Rate Limiting**: Protection for /auth/* endpoints
- **CORS**: Configurable origins

---

## 🛠️ Development

### Project Structure

```
saa-risk-analyzer/
├── backend/          # Go API
│   ├── cmd/api/      # Entry point
│   ├── internal/     # Business logic
│   │   ├── math/     # Risk engine
│   │   ├── handlers/ # HTTP handlers
│   │   ├── service/  # Business services
│   │   ├── repo/     # Data access
│   │   └── auth/     # JWT, RBAC
│   └── api/          # OpenAPI spec
├── frontend/         # React SPA
│   └── src/
│       ├── pages/    # Pages
│       ├── components/ # UI components
│       └── store/    # Zustand stores
├── data/             # Demo CSV
└── scripts/          # Utilities
```

### Local Development (without Docker)

**Backend:**
```bash
cd backend
cp ../.env.example .env
# Configure DB_HOST=localhost in .env
go run cmd/api/main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

---

## 🚢 Production Deployment

### With Domain (Let's Encrypt auto-certificate)

1. Edit `Caddyfile` - uncomment section with domain
2. Set correct DNS A records pointing to your server
3. Run: `docker-compose up -d`

Caddy will automatically obtain SSL certificate!

### Digital Ocean Deployment

```bash
# Build frontend
cd frontend
npm run build

# Build backend
cd backend
go build -o saa-risk-analyzer cmd/api/main.go

# Deploy with PM2
pm2 start saa-risk-analyzer --name saa-risk-analyzer

# Configure Nginx
# Location: /etc/nginx/sites-enabled/risk-analyzer
```

---

## 📄 License

Proprietary - Scientific Analytics Alliance

---

## 👥 Author

**Scientific Analytics Alliance**

Premium Research & Wealth Intelligence Platform

---

## 🔗 Related Projects

- [Crypto Analytics Portal](../crypto_reports)
- [SAA Learn Your Way](../saa-learn-your-way)
- [Liquidity Positioner](../liquidity-positioner)
- [News Analytics AI](../signal-analysis)

---

**Last Updated:** November 2025

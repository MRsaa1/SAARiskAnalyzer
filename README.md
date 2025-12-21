# SAA Risk Analyzer

**Institutional-Grade Portfolio Risk Management Platform**

[![Live Demo](https://img.shields.io/badge/Demo-Live-brightgreen)](https://risk.saa-alliance.com)
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

**Production URL:** [https://risk.saa-alliance.com](https://risk.saa-alliance.com)

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
   - VaR/CVaR calculation with multiple methods
   - Stress testing with historical scenarios
   - Backtesting and validation

4. **Analytics**
   - Risk contribution analysis
   - Correlation heatmaps
   - PCA factor decomposition

---

## 🛠️ Installation & Setup

### Prerequisites

- Go 1.22+
- Node.js 18+ and npm
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Backend Setup

```bash
cd backend
go mod download
go build -o saa-risk-analyzer cmd/api/main.go
./saa-risk-analyzer
```

### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

### Environment Variables

Create `.env` file:

```env
# Database
DATABASE_URL=postgresql://user:password@localhost:5432/risk_analyzer

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h

# Server
PORT=8080
GIN_MODE=release
```

---

## 📊 API Endpoints

### Risk Calculation

```
POST /api/v1/risk/var
POST /api/v1/risk/cvar
POST /api/v1/risk/contribution
GET /api/v1/risk/correlation
POST /api/v1/risk/pca
```

### Stress Testing

```
POST /api/v1/stress/historical
POST /api/v1/stress/custom
GET /api/v1/stress/scenarios
```

### Backtesting

```
POST /api/v1/backtest/kupiec
POST /api/v1/backtest/christoffersen
```

**Full API Documentation**: Available via Swagger UI at `/swagger/index.html`

---

## 🔧 Configuration

### VaR Parameters

- **Confidence Level**: 95%, 99%, 99.9%
- **Time Horizon**: 1 day, 1 week, 1 month
- **Methods**: Historical, Parametric, Monte Carlo
- **Monte Carlo Paths**: 10,000 - 100,000

### Stress Test Scenarios

Pre-configured historical scenarios:
- COVID-2020: Market crash March 2020
- Lehman 2008: Financial crisis September 2008
- Tech Bubble 2000: Dot-com crash

---

## 🧪 Testing

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

---

## 📈 Performance

- **VaR Calculation**: < 100ms for standard portfolios
- **Monte Carlo (10k paths)**: 2-5 seconds
- **Monte Carlo (100k paths)**: 15-30 seconds
- **Stress Testing**: 1-3 seconds per scenario
- **Database Queries**: Optimized with indexes and connection pooling

---

## 🔒 Security

- JWT-based authentication
- RBAC (Role-Based Access Control)
- Argon2id password hashing
- Input validation and sanitization
- CORS configuration
- Rate limiting

---

## 📄 License

Proprietary - Scientific Analytics Alliance

---

## 👥 Author

**Scientific Analytics Alliance**

Premium Research & Wealth Intelligence Platform

---

## 🔗 Related Projects

- [Global Risk Intelligence Platform](../Global-Risk-Intelligence-Platform)
- [Investment Dashboard](../investment-bot)
- [Crypto Analytics Portal](../CryptoAnalyticsPortal)
- [ARIN Platform](../arin-platform)

---

## 📞 Support

For questions or support, please contact: [support@saa-alliance.com](mailto:support@saa-alliance.com)

---

**Last Updated:** December 2025

**Production Domain:** https://risk.saa-alliance.com


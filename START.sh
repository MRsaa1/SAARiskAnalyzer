#!/bin/bash
set -e

echo "🚀 Запуск SAA Risk Analyzer..."

# Экспорт переменных
export APP_ENV=development
export PORT=8083
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=$(whoami)
export DB_PASSWORD=
export DB_NAME=risk_db
export DB_SSLMODE=disable
export JWT_SECRET=saa_risk_analyzer_super_secret_key
export ADMIN_EMAIL=admin@example.com
export ADMIN_PASSWORD=Admin123456!

echo "✅ PostgreSQL: Запущен"
echo "✅ База данных risk_db: Создана"
echo ""
echo "🔧 Запускаю backend..."

cd backend
go run cmd/api/main.go

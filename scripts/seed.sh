#!/bin/bash
set -e

echo "Seeding database with demo data..."

# Wait for API to be ready
echo "Waiting for API..."
until curl -s http://localhost:8083/health > /dev/null; do
  echo "API not ready yet, waiting..."
  sleep 2
done

echo "API is ready!"

# Login as admin
TOKEN=$(curl -s -X POST http://localhost:8083/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"Admin123456!"}' \
  | jq -r '.token')

echo "Logged in, token: ${TOKEN:0:20}..."

# Create portfolio
PORTFOLIO_ID=$(curl -s -X POST http://localhost:8083/api/portfolios \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Demo Portfolio","description":"Sample multi-asset portfolio"}' \
  | jq -r '.id')

echo "Created portfolio: $PORTFOLIO_ID"

echo "Seed completed!"

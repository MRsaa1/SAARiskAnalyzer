#!/bin/bash

echo ""
echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║              🔍 ПРОВЕРКА СИСТЕМЫ SAA RISK ANALYZER               ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo ""

# Backend
echo "🔹 Backend (8084):"
if curl -s http://localhost:8084/api/dashboard > /dev/null 2>&1; then
    echo "   ✅ Работает и отвечает"
    PID=$(lsof -ti:8084 2>/dev/null)
    echo "   📝 PID: $PID"
else
    echo "   ❌ Не отвечает"
fi
echo ""

# Frontend
echo "🔹 Frontend (3000):"
if curl -s http://localhost:3000 > /dev/null 2>&1; then
    echo "   ✅ Работает и отвечает"
    PID=$(lsof -ti:3000 2>/dev/null)
    echo "   📝 PID: $PID"
else
    echo "   ❌ Не отвечает"
fi
echo ""

# PostgreSQL
echo "🔹 PostgreSQL (5432):"
if /opt/homebrew/opt/postgresql@15/bin/psql -U postgres -d risk_db -c "SELECT 1" > /dev/null 2>&1; then
    echo "   ✅ Работает, база доступна"
else
    echo "   ⚠️  Проблемы с доступом"
fi
echo ""

# LaunchAgents
echo "🔹 LaunchAgents (автозапуск):"
launchctl list | grep com.saa | while read -r line; do
    pid=$(echo "$line" | awk '{print $1}')
    status=$(echo "$line" | awk '{print $2}')
    name=$(echo "$line" | awk '{print $3}')
    
    if [ "$pid" = "-" ]; then
        echo "   ⚠️  $name - не запущен (код: $status)"
    else
        echo "   ✅ $name - работает (PID: $pid)"
    fi
done
echo ""

# Логи
echo "🔹 Последние логи Backend:"
tail -3 ~/saa-risk-analyzer/backend_error.log 2>/dev/null | sed 's/^/   /' || echo "   (нет логов)"
echo ""

echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║  Доступ: http://localhost:3000 | API: :8084/api/dashboard       ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo ""

#!/bin/bash

echo "📊 Загрузка демо-данных в базу данных..."

PSQL="/opt/homebrew/opt/postgresql@15/bin/psql"
DB="risk_db"
USER="postgres"
DATA_DIR="/Users/artur220513timur110415gmail.com/saa-risk-analyzer/data"

# Проверка что файлы существуют
if [ ! -f "$DATA_DIR/prices.csv" ]; then
    echo "❌ Файл prices.csv не найден!"
    exit 1
fi

echo "1️⃣  Создаю временную таблицу для импорта..."

$PSQL -U $USER -d $DB <<EOF
-- Создаем временную таблицу
CREATE TEMP TABLE temp_prices (
    date TEXT,
    symbol TEXT,
    close TEXT
);

-- Загружаем CSV
\copy temp_prices(date,symbol,close) FROM '$DATA_DIR/prices.csv' WITH (FORMAT csv, HEADER true);

-- Создаем/обновляем активы
INSERT INTO assets (id, symbol, name, class, currency, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    symbol,
    symbol,
    'Unknown',
    'USD',
    NOW(),
    NOW()
FROM (SELECT DISTINCT symbol FROM temp_prices) s
ON CONFLICT (symbol) DO NOTHING;

-- Загружаем цены
INSERT INTO prices (id, asset_id, date, close, created_at)
SELECT 
    gen_random_uuid(),
    a.id,
    tp.date::timestamp,
    tp.close::numeric,
    NOW()
FROM temp_prices tp
JOIN assets a ON a.symbol = tp.symbol
ON CONFLICT DO NOTHING;

-- Статистика
SELECT 
    'Assets' as table_name,
    COUNT(*) as count
FROM assets
UNION ALL
SELECT 
    'Prices' as table_name,
    COUNT(*) as count
FROM prices;
EOF

echo "✅ Демо-данные загружены!"
echo ""
echo "📊 Проверьте:"
$PSQL -U $USER -d $DB -c "SELECT symbol, COUNT(*) as price_count FROM prices p JOIN assets a ON p.asset_id = a.id GROUP BY symbol ORDER BY symbol;"


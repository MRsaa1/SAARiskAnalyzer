# 🌐 Порты и сервисы

## 📊 Активные сервисы

| Сервис | Порт | URL | Статус |
|--------|------|-----|--------|
| **Backend API** | 8084 | http://localhost:8084/api/dashboard | ✅ Работает |
| **Frontend** | 3000 | http://localhost:3000 | ✅ Работает |
| **PostgreSQL** | 5432 | localhost:5432 | ✅ Работает |

## 🔧 LaunchAgents (автозапуск)

| Сервис | Label | PID | Автоперезапуск |
|--------|-------|-----|----------------|
| Backend | com.saa.backend | Динамический | ✅ Да (10 сек) |
| Frontend | com.saa.frontend | Динамический | ✅ Да (10 сек) |
| MinIO | com.saa.minio | 73178 | ✅ Да |

## 📋 API Endpoints

### Dashboard
```bash
curl http://localhost:8084/api/dashboard
```

Возвращает:
```json
{
  "contributors": [...],
  "cvar_1d": 187654.3,
  "var_1d": 125432.5,
  "vol": 0.154
}
```

### Health Check
```bash
curl http://localhost:8084/api/health
```

## 🗄️ База данных

- **Имя**: `risk_db`
- **Пользователь**: `postgres`
- **Пароль**: `postgres`
- **Хост**: `localhost`
- **Порт**: `5432`

Подключение:
```bash
/opt/homebrew/opt/postgresql@15/bin/psql -U postgres -d risk_db
```

## 📝 Логи

```bash
# Backend
tail -f ~/saa-risk-analyzer/backend_error.log

# Frontend
tail -f ~/saa-risk-analyzer/frontend_error.log
```

## ⚙️ Управление

```bash
# Статус
./manage-services.sh status

# Перезапуск
./manage-services.sh restart

# Тест автоперезапуска
./manage-services.sh test

# Логи
./manage-services.sh logs
```

---
**Обновлено**: $(date)

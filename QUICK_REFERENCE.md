# 🚀 Быстрая справка SAA Risk Analyzer

## Основные команды

```bash
# Проверить статус
./manage-services.sh status

# Перезапустить всё
./manage-services.sh restart

# Посмотреть логи
./manage-services.sh logs
```

## Доступ к приложению

- **Frontend**: http://localhost:3001
- **Backend API**: http://localhost:8084
- **Health Check**: http://localhost:8084/health

## Если что-то не работает

```bash
# Перезапустить сервисы
./manage-services.sh restart

# Проверить PostgreSQL
brew services list | grep postgresql

# Запустить PostgreSQL (если нужно)
brew services start postgresql@15
```

## Логи

```bash
# Backend
tail -f /tmp/saa-backend.error.log

# Frontend
tail -f /tmp/saa-frontend.log
```

## Управление автозапуском

```bash
# Отключить автозапуск
./manage-services.sh uninstall

# Включить автозапуск
./manage-services.sh install
```

---

📚 **Подробная документация**: [AUTOSTART.md](AUTOSTART.md)




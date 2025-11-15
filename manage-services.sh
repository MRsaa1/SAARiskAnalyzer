#!/bin/bash

# Скрипт управления сервисами SAA Risk Analyzer

BACKEND_PLIST="$HOME/Library/LaunchAgents/com.saa.backend.plist"
FRONTEND_PLIST="$HOME/Library/LaunchAgents/com.saa.frontend.plist"

case "$1" in
    start)
        echo "🚀 Запуск сервисов..."
        launchctl start com.saa.backend
        launchctl start com.saa.frontend
        echo "✅ Сервисы запущены"
        ;;
    stop)
        echo "🛑 Остановка сервисов..."
        launchctl stop com.saa.backend
        launchctl stop com.saa.frontend
        echo "✅ Сервисы остановлены"
        ;;
    restart)
        echo "🔄 Перезапуск сервисов..."
        launchctl stop com.saa.backend
        launchctl stop com.saa.frontend
        sleep 2
        launchctl start com.saa.backend
        launchctl start com.saa.frontend
        echo "✅ Сервисы перезапущены"
        ;;
    status)
        echo "📊 Статус сервисов:"
        echo ""
        launchctl list | grep com.saa
        echo ""
        echo "📡 Порты:"
        echo "Backend (8084):"
        lsof -i :8084 2>/dev/null || echo "  ❌ Не запущен"
        echo ""
        echo "Frontend (3001):"
        lsof -i :3001 2>/dev/null || echo "  ❌ Не запущен"
        echo ""
        echo "🔍 Проверка API:"
        if curl -s http://localhost:8084/health > /dev/null 2>&1; then
            echo "  ✅ Backend API работает"
        else
            echo "  ❌ Backend API не отвечает"
        fi
        if curl -s http://localhost:3001 > /dev/null 2>&1; then
            echo "  ✅ Frontend работает"
        else
            echo "  ❌ Frontend не отвечает"
        fi
        ;;
    logs)
        echo "📋 Логи Backend (stdout):"
        tail -30 ~/saa-risk-analyzer/backend_error.log 2>/dev/null || echo "Логов нет"
        echo ""
        echo "📋 Логи Frontend (errors):"
        tail -30 ~/saa-risk-analyzer/frontend_error.log 2>/dev/null || echo "Логов нет"
        ;;
    install)
        echo "📦 Установка автозапуска..."
        launchctl load "$BACKEND_PLIST"
        launchctl load "$FRONTEND_PLIST"
        echo "✅ Автозапуск установлен"
        ;;
    uninstall)
        echo "🗑️  Удаление автозапуска..."
        launchctl unload "$BACKEND_PLIST"
        launchctl unload "$FRONTEND_PLIST"
        echo "✅ Автозапуск удален"
        ;;
    test)
        echo "🧪 Запуск теста автоматического перезапуска..."
        ~/saa-risk-analyzer/test-auto-restart.sh
        ;;
    *)
        echo "Использование: $0 {start|stop|restart|status|logs|test|install|uninstall}"
        echo ""
        echo "Команды:"
        echo "  start      - Запустить сервисы"
        echo "  stop       - Остановить сервисы"
        echo "  restart    - Перезапустить сервисы"
        echo "  status     - Показать статус сервисов и проверить API"
        echo "  logs       - Показать логи"
        echo "  test       - Протестировать автоматический перезапуск"
        echo "  install    - Установить автозапуск"
        echo "  uninstall  - Удалить автозапуск"
        exit 1
        ;;
esac

exit 0




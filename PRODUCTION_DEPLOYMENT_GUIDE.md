# 🚀 Production Deployment Guide - Fixing Black Screen Issue

## 🔍 Проблема

**Симптомы:**
- ✅ Dev server (localhost:3000) работает
- ❌ Production build на сервере → черный экран
- ✅ HTML/CSS/JS файлы загружаются (HTTP 200)
- ✅ Backend API работает
- ❌ React приложение не рендерится

## 🎯 Причины черного экрана в production

1. **Жестко закодированный API URL** в `api.ts`
2. **Отсутствие переменных окружения** для production
3. **Неправильная конфигурация Vite** для production build
4. **Возможные JavaScript ошибки** (нужно проверить в консоли браузера)
5. **CORS проблемы** (если API на другом домене)

---

## 📋 Пошаговая инструкция по исправлению

### Шаг 1: Создайте файлы переменных окружения

#### Файл 1: `frontend/.env.development`
```bash
# Для локальной разработки
VITE_API_BASE_URL=http://localhost:8084/api
```

#### Файл 2: `frontend/.env.production`
```bash
# Для production сервера
VITE_API_BASE_URL=http://104.248.70.69:8087/api
```

**Как создать:**
```bash
cd frontend

# Development
echo "VITE_API_BASE_URL=http://localhost:8084/api" > .env.development

# Production
echo "VITE_API_BASE_URL=http://104.248.70.69:8087/api" > .env.production
```

---

### Шаг 2: Исправьте файл `frontend/src/lib/api.ts`

**Было (неправильно):**
```typescript
import axios from 'axios'

const api = axios.create({
  baseURL: 'http://104.248.70.69:8087/api',  // ❌ Жестко закодировано
})
```

**Должно быть (правильно):**
```typescript
import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8084/api',
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export default api
```

**Что изменилось:**
- Теперь API URL берется из переменной окружения
- Fallback на localhost для разработки

---

### Шаг 3: Обновите `frontend/vite.config.ts`

**Добавьте секцию build:**

```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  base: '/',  // ✅ Важно! Корневой путь
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: true,  // ✅ Для отладки в production
    rollupOptions: {
      output: {
        manualChunks: undefined,
      },
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8083',
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path,
      },
    },
  },
})
```

**Ключевые параметры:**
- `base: '/'` - приложение работает с корня домена
- `sourcemap: true` - поможет найти ошибки в production

---

### Шаг 4: Добавьте обработку ошибок в `App.tsx`

**Добавьте Error Boundary для отлова ошибок:**

```typescript
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { useEffect } from 'react'
import DashboardPage from './pages/Dashboard'
import PortfolioImport from './pages/PortfolioImport'
import AnalyticsPage from './pages/Analytics'

function App() {
  useEffect(() => {
    // Логирование для отладки
    console.log('🚀 App mounted')
    console.log('📡 API Base URL:', import.meta.env.VITE_API_BASE_URL)
    console.log('🏗️ Mode:', import.meta.env.MODE)
  }, [])

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<DashboardPage />} />
        <Route path="/import" element={<PortfolioImport />} />
        <Route path="/analytics" element={<AnalyticsPage />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
```

---

### Шаг 5: Соберите production build

```bash
cd frontend

# Очистите старую сборку
rm -rf dist

# Соберите для production
npm run build

# Проверьте что build создался
ls -la dist/
```

**Ожидаемый результат:**
```
dist/
├── index.html
├── assets/
│   ├── index-[hash].js
│   ├── index-[hash].css
│   └── ...
└── vite.svg
```

---

### Шаг 6: Проверьте build локально

```bash
cd frontend

# Запустите preview сервер
npm run preview

# Откройте в браузере
open http://localhost:4173
```

**Если preview работает** → значит сборка правильная, проблема на сервере.  
**Если preview НЕ работает** → проблема в коде, проверьте консоль браузера.

---

### Шаг 7: Проверьте CORS на backend

**Файл:** `backend/cmd/api/main.go`

**Убедитесь что CORS настроен правильно:**

```go
// CORS middleware
router.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }
    
    c.Next()
})
```

**Или для конкретного домена:**
```go
c.Writer.Header().Set("Access-Control-Allow-Origin", "http://104.248.70.69:3001")
```

---

## 🖥️ Deployment на сервер

### Вариант 1: Через SCP (простой способ)

```bash
# 1. Соберите frontend
cd frontend
npm run build

# 2. Скопируйте на сервер
scp -r dist/* root@104.248.70.69:/var/www/saa-frontend/

# 3. Проверьте на сервере
ssh root@104.248.70.69
ls -la /var/www/saa-frontend/
```

### Вариант 2: Через Git (рекомендуется)

```bash
# На сервере
ssh root@104.248.70.69

# Клонируйте репозиторий
cd /opt
git clone <your-repo-url> saa-risk-analyzer
cd saa-risk-analyzer/frontend

# Установите зависимости
npm install

# Соберите
npm run build

# Скопируйте в nginx директорию
cp -r dist/* /var/www/saa-frontend/
```

---

## 🔧 Конфигурация Nginx на сервере

**Файл:** `/etc/nginx/sites-available/saa-frontend`

```nginx
server {
    listen 3001;
    server_name 104.248.70.69;

    root /var/www/saa-frontend;
    index index.html;

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

    location / {
        # SPA fallback - ВСЕ несуществующие пути идут на index.html
        try_files $uri $uri/ /index.html;
        
        # Cache control
        add_header Cache-Control "no-cache, must-revalidate";
    }

    # Кэширование статических файлов
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # CORS headers (если нужны на frontend)
    add_header Access-Control-Allow-Origin "*" always;
    add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS" always;
    add_header Access-Control-Allow-Headers "Content-Type, Authorization" always;
}
```

**Активируйте конфигурацию:**
```bash
# На сервере
sudo ln -s /etc/nginx/sites-available/saa-frontend /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 🔍 Диагностика проблем

### Проблема 1: Черный экран (самое частое)

**Проверьте в консоли браузера:**
1. Откройте DevTools (F12)
2. Перейдите в Console
3. Обновите страницу (Ctrl+R)
4. Ищите **красные ошибки**

**Типичные ошибки:**

#### Ошибка: "Failed to fetch"
```
Причина: CORS или недоступный API
Решение: Проверьте CORS на backend, проверьте что API доступен
```

#### Ошибка: "Unexpected token '<'"
```
Причина: JavaScript файлы не загружаются или возвращается HTML вместо JS
Решение: Проверьте nginx конфигурацию, проверьте пути к файлам
```

#### Ошибка: "Cannot read property of undefined"
```
Причина: JavaScript ошибка в коде
Решение: Проверьте sourcemap, найдите строку с ошибкой
```

### Проблема 2: 404 на JavaScript файлах

**Проверьте:**
```bash
# На сервере
ls -la /var/www/saa-frontend/
ls -la /var/www/saa-frontend/assets/

# Права доступа
sudo chown -R www-data:www-data /var/www/saa-frontend/
sudo chmod -R 755 /var/www/saa-frontend/
```

### Проблема 3: API недоступен

**Проверьте backend:**
```bash
# На сервере
curl http://104.248.70.69:8087/api/health
curl http://104.248.70.69:8087/api/dashboard

# Проверьте что backend запущен
ps aux | grep go
netstat -tulpn | grep 8087
```

### Проблема 4: CORS ошибки

**В консоли видите:**
```
Access to fetch at 'http://104.248.70.69:8087/api/...' from origin 'http://104.248.70.69:3001' has been blocked by CORS policy
```

**Решение:**
1. Добавьте CORS headers в backend (см. Шаг 7)
2. Или добавьте proxy в nginx:

```nginx
location /api/ {
    proxy_pass http://104.248.70.69:8087/api/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

---

## 📝 Чек-лист перед deployment

### Frontend
- [ ] Созданы файлы `.env.development` и `.env.production`
- [ ] В `api.ts` используется `import.meta.env.VITE_API_BASE_URL`
- [ ] В `vite.config.ts` добавлена секция `build` с `base: '/'`
- [ ] Выполнена команда `npm run build`
- [ ] Проверена сборка локально с `npm run preview`
- [ ] Нет ошибок в консоли браузера

### Backend
- [ ] CORS настроен правильно
- [ ] Backend доступен по URL `http://104.248.70.69:8087/api`
- [ ] Endpoint `/health` отвечает
- [ ] Endpoint `/api/dashboard` возвращает данные

### Server
- [ ] Nginx установлен и запущен
- [ ] Конфигурация nginx содержит `try_files $uri $uri/ /index.html`
- [ ] Файлы frontend скопированы в `/var/www/saa-frontend/`
- [ ] Права доступа правильные (`www-data:www-data`)
- [ ] Nginx перезагружен после изменений

---

## 🧪 Тестирование после deployment

### 1. Проверьте что файлы доступны
```bash
# HTML
curl -I http://104.248.70.69:3001/

# JavaScript
curl -I http://104.248.70.69:3001/assets/index-*.js

# CSS
curl -I http://104.248.70.69:3001/assets/index-*.css
```

Все должны возвращать **200 OK**.

### 2. Проверьте API
```bash
curl http://104.248.70.69:8087/api/health
curl http://104.248.70.69:8087/api/dashboard
```

### 3. Откройте в браузере
```
http://104.248.70.69:3001/
```

### 4. Откройте консоль браузера (F12)
- Должны быть логи: `🚀 App mounted`, `📡 API Base URL:...`
- Не должно быть красных ошибок
- Проверьте Network tab - все файлы должны загружаться (200 OK)

---

## 🚨 Если все еще черный экран

### Шаг 1: Включите sourcemap
В `vite.config.ts`:
```typescript
build: {
  sourcemap: true,  // ✅ Добавьте это
}
```

Пересоберите:
```bash
npm run build
```

### Шаг 2: Добавьте console.log везде

В `main.tsx`:
```typescript
console.log('🔥 main.tsx loaded')

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)

console.log('✅ React rendered')
```

В `App.tsx`:
```typescript
function App() {
  console.log('🎯 App component loaded')
  console.log('API URL:', import.meta.env.VITE_API_BASE_URL)
  
  return (
    // ...
  )
}
```

### Шаг 3: Создайте минимальный тест

**Файл:** `frontend/test.html`
```html
<!DOCTYPE html>
<html>
<head>
    <title>Test</title>
</head>
<body>
    <h1>Test Page</h1>
    <div id="root">Loading...</div>
    <script>
        console.log('Test page loaded');
        document.getElementById('root').innerHTML = 'React Test';
        
        // Тест API
        fetch('http://104.248.70.69:8087/api/health')
            .then(r => r.json())
            .then(data => console.log('API works:', data))
            .catch(e => console.error('API error:', e));
    </script>
</body>
</html>
```

Скопируйте на сервер и откройте: `http://104.248.70.69:3001/test.html`

Если этот файл работает → проблема в React build.  
Если не работает → проблема в nginx/server конфигурации.

---

## 📞 Куда смотреть дальше

### Логи сервера
```bash
# Nginx access log
tail -f /var/log/nginx/access.log

# Nginx error log
tail -f /var/log/nginx/error.log

# Backend logs
journalctl -u saa-backend -f
```

### Консоль браузера
1. Откройте DevTools (F12)
2. Console tab - ищите ошибки
3. Network tab - проверьте что загружается
4. Sources tab - проверьте исходники с sourcemap

### Типичные места ошибок
1. ❌ Неправильный API URL
2. ❌ CORS не настроен
3. ❌ Nginx неправильно настроен
4. ❌ JavaScript ошибки в коде
5. ❌ Отсутствует файл index.html
6. ❌ Неправильные пути к assets

---

## ✅ После исправления

После того как всё заработает, **зафиксируйте изменения:**

```bash
git add .
git commit -m "Fix production build configuration"
git push
```

**Создайте deployment скрипт:**

```bash
#!/bin/bash
# deploy.sh

echo "🚀 Deploying SAA Risk Analyzer..."

# Build frontend
cd frontend
npm run build

# Copy to server
scp -r dist/* root@104.248.70.69:/var/www/saa-frontend/

# Restart nginx
ssh root@104.248.70.69 'sudo systemctl reload nginx'

echo "✅ Deployment complete!"
echo "🌐 Open: http://104.248.70.69:3001"
```

---

## 📚 Итоговый чек-лист

### Локально (перед deployment)
- [ ] Созданы `.env.development` и `.env.production`
- [ ] `api.ts` использует переменные окружения
- [ ] `vite.config.ts` настроен для production
- [ ] `npm run build` успешно выполняется
- [ ] `npm run preview` показывает работающее приложение
- [ ] Нет ошибок в консоли

### На сервере
- [ ] Backend работает на порту 8087
- [ ] CORS настроен в backend
- [ ] Nginx настроен с `try_files` для SPA
- [ ] Frontend файлы скопированы в `/var/www/saa-frontend/`
- [ ] Права доступа установлены правильно
- [ ] Nginx перезапущен

### В браузере
- [ ] Страница открывается (не 404)
- [ ] Нет красных ошибок в Console
- [ ] Все файлы загружаются (Network tab)
- [ ] API запросы работают
- [ ] Приложение рендерится

---

**Если следовать этой инструкции, черный экран должен исчезнуть! 🎉**

Если проблема останется - откройте консоль браузера и пришлите скриншот ошибок.


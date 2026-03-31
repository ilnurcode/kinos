# Kinos Frontend

Frontend приложение для Kinos - интернет-магазина электроники.

## 🛠 Технологии

- **Vue.js 3** - прогрессивный фреймворк
- **Vite** - сборщик проектов
- **Pinia** - управление состоянием
- **Vue Router 4** - навигация
- **Axios** - HTTP клиент
- **Bootstrap 5** - UI библиотека
- **SCSS** - препроцессор CSS

## 📦 Установка

```bash
# Установить зависимости
npm install

# Запустить в режиме разработки
npm run dev

# Собрать для production
npm run build

# Запустить production сборку
npm run preview
```

## 🌐 Переменные окружения

Создайте файл `.env` в корне проекта:

```env
VITE_API_BASE_URL=http://localhost:8080/api
VITE_APP_NAME=Kinos
VITE_APP_VERSION=1.0.0
```

## 📁 Структура проекта

```
frontend/
├── public/              # Статические файлы
├── src/
│   ├── api/            # API модули
│   │   ├── client.js   # Axios клиент
│   │   └── auth.js     # Auth API
│   ├── assets/         # Ресурсы (CSS, изображения)
│   ├── components/     # Vue компоненты
│   │   └── common/     # Общие компоненты
│   ├── router/         # Vue Router
│   ├── stores/         # Pinia stores
│   ├── utils/          # Утилиты
│   ├── views/          # Страницы
│   │   └── admin/      # Админ панели
│   ├── App.vue         # Корневой компонент
│   └── main.js         # Точка входа
├── index.html
├── package.json
├── vite.config.js
└── .env
```

## 🚀 Запуск с бэкендом

1. Запустить бэкенд сервисы:
```bash
cd ..
docker-compose up -d
```

2. Запустить frontend:
```bash
cd frontend
npm run dev
```

3. Открыть браузер: http://localhost:3000

## 📝 Основные возможности

### Авторизация
- Регистрация нового пользователя
- Вход с email/паролем
- JWT аутентификация
- Автоматический refresh токена
- Защита роутов

### Навигация
- Публичные страницы (Главная, Каталог, Вход, Регистрация)
- Защищённые страницы (Профиль)
- Админ панель (только для admin роли)

## 🔧 Конфигурация

### Vite

Находится в `vite.config.js`:
- Proxy для API запросов
- Aliases для импортов (@)
- Настройки сборки

### Bootstrap

Импортируется в `src/main.js`:
- CSS Bootstrap 5
- Кастомные стили в `src/assets/css/main.scss`

## 📄 Лицензия

MIT

# Blockchain Server

![Go](https://img.shields.io/badge/Go-1.20-blue)
![React](https://img.shields.io/badge/React-18.2.0-blue)
![License](https://img.shields.io/badge/License-MIT-green)

## 📖 Описание

Этот проект представляет собой сервер блокчейна, реализованный на языке Go с использованием React для фронтенда. Он поддерживает создание блоков, майнинг, регистрацию кошельков и управление транзакциями.

---

## 🚀 Основные возможности

- **Майнинг блоков**: Автоматическое добавление транзакций в блоки.
- **Регистрация кошельков**: Создание новых кошельков с уникальными адресами.
- **Управление транзакциями**: Добавление и проверка транзакций.
- **Консенсус**: Разрешение конфликтов между узлами.
- **API**: REST API для взаимодействия с блокчейном.

---

## 🛠️ Технологии

- **Backend**: Go
- **Frontend**: React
- **HTTP Framework**: Gorilla Mux
- **Криптография**: Подписи и хеширование для безопасности.

---

## 📂 Структура проекта

```plaintext
├── block/
│   ├── blockchain.go       # Логика блокчейна
│   ├── mining.go           # Майнинг и Proof-of-Work
│   ├── transaction.go      # Управление транзакциями
├── server/
│   ├── handlers/           # Обработчики API
│   ├── middleware/         # Middleware для CORS и логирования
├── struct/
│   ├── wallet/             # Логика кошельков
│   ├── utils/              # Утилиты
├── main.go                 # Точка входа
├── go.mod                  # Зависимости Go
├── package.json            # Зависимости React

📦 Установка и запуск
Backend
Установите Go.
Клонируйте репозиторий:
git clone https://github.com/ваш-аккаунт/blockchain-server.git
cd blockchain-server
Установите зависимости:
go mod tidy
Запустите сервер:
go run main.go
Frontend
Перейдите в папку фронтенда:
cd frontend
Установите зависимости:
npm install
Запустите React-приложение:
npm start
<hr></hr>
📖 API Документация
Основные эндпоинты
GET /chain: Получить текущую цепочку блоков.
POST /transactions: Добавить новую транзакцию.
GET /balance: Получить баланс кошелька.
PUT /consensus: Запустить процесс консенсуса.
Пример запроса:
curl -X POST http://localhost:5001/transactions \
-H "Content-Type: application/json" \
-d '{
  "sender": "sender_1",
  "recipient": "recipient_1",
  "amount": 10.0
}'

🧪 Тестирование
Запустите сервер.
Используйте curl или Postman для тестирования API.
Убедитесь, что транзакции добавляются, а блоки майнятся корректно.
<hr></hr>
🤝 Вклад
Будем рады вашему вкладу! Открывайте issues или создавайте pull requests.  <hr></hr>
📜 Лицензия
Этот проект распространяется под лицензией MIT.<hr></hr>
📧 Контакты
Автор: Almaz Toktassin

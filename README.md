# Blockchain Server

![Go](https://img.shields.io/badge/Go-1.20-blue)
![React](https://img.shields.io/badge/React-18.2.0-blue)
![License](https://img.shields.io/badge/License-MIT-green)

## 📖 Описание

Blockchain Server — это серверная реализация блокчейна на языке Go с фронтендом на React. Проект поддерживает создание блоков, майнинг, регистрацию кошельков и управление транзакциями с помощью REST API.

---

## 🚀 Основные возможности

- 🔗 **Майнинг блоков** — автоматическое добавление транзакций и генерация новых блоков.
- 👛 **Регистрация кошельков** — создание новых кошельков с уникальными адресами.
- 💸 **Управление транзакциями** — создание, валидация и добавление транзакций в цепочку.
- 🤝 **Консенсус** — разрешение конфликтов между узлами сети.
- 📡 **REST API** — удобное взаимодействие с сервером.

---

## 🛠️ Технологии

- **Backend**: Go
- **Frontend**: React
- **HTTP Framework**: Gorilla Mux
- **Криптография**: Хеширование и цифровые подписи для безопасности транзакций

---

## 📂 Структура проекта

```plaintext
├── block/
│   ├── blockchain.go       # Логика блокчейна
│   ├── mining.go           # Proof-of-Work
│   ├── transaction.go      # Транзакции
├── server/
│   ├── handlers/           # Обработчики API
│   ├── middleware/         # CORS и логгеры
├── struct/
│   ├── wallet/             # Кошельки
│   ├── utils/              # Вспомогательные функции
├── main.go                 # Точка входа
├── go.mod                  # Зависимости Go
├── frontend/
│   ├── package.json        # Зависимости React

```

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

    Перейдите в директорию фронтенда:

cd frontend

Установите зависимости:

npm install

Запустите приложение:

    npm start

📖 API Документация
Основные эндпоинты
Метод	Эндпоинт	Описание
GET	/chain	Получить текущую цепочку блоков
POST	/transactions	Добавить новую транзакцию
GET	/balance	Получить баланс кошелька
PUT	/consensus	Запустить процесс консенсуса

Пример запроса:
```
curl -X POST http://localhost:5001/transactions \
-H "Content-Type: application/json" \
-d '{
  "sender": "sender_1",
  "recipient": "recipient_1",
  "amount": 10.0
}'
```
🧪 Тестирование

    Запустите сервер.

    Используйте curl или Postman для тестирования REST API.

    Проверьте корректность транзакций и майнинга.

🤝 Вклад

Будем рады вашему участию!
Создавайте issues или отправляйте pull requests.
📜 Лицензия

Проект распространяется под лицензией MIT.
📧 Контакты

Автор: Almaz Toktassin
📬 Email: almaztok8@gmail.com
💻 GitHub: nighbee

Bank API

REST API для банковского сервиса на Go.

Функции:

Регистрация и вход по JWT

Создание и управление счетами в рублях

Выпуск виртуальных карт с шифрованием

Переводы между счетами

Оформление кредитов с графиком платежей

Интеграция с API Центробанка для ключевой ставки

Отправка уведомлений на email

Аналитика доходов и расходов

Прогноз баланса

Запуск:

Установить Go 1.23+, PostgreSQL 17+

Создать базу bankapi и пользователя bankuser

Выполнить миграции из папки migrations

Скопировать .env.example в .env и заполнить

Запустить go run main.go

Сервер будет доступен на порту 8080

API:

Публичные запросы:
POST /register - регистрация, тело: username, email, password
POST /login - вход, тело: username, password

Защищенные запросы, нужен заголовок Authorization: Bearer токен:
POST /accounts - создать счет, тело: currency
GET /accounts - список счетов
POST /deposit - пополнить, тело: account_id, amount
POST /transfer - перевод, тело: from_account_id, to_account_id, amount
POST /cards - выпустить карту, тело: account_id
GET /cards - список карт
POST /credits - оформить кредит, тело: amount, term_months
GET /credits - список кредитов
GET /credits/id/schedule - график платежей

Безопасность:

Пароли хешируются bcrypt

CVV карт хешируется bcrypt

Номера карт шифруются PGP

Целостность данных через HMAC-SHA256

Доступ к чужим счетам запрещен

Автор: Timur Nurzhau
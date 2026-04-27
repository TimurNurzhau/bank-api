Bank API

REST API для банковского сервиса на Go.

Функции:

- Регистрация и вход по JWT
- Создание и управление счетами в рублях
- Выпуск виртуальных карт с PGP-шифрованием
- Переводы между счетами
- Оформление кредитов с графиком платежей
- Интеграция с API Центробанка для ключевой ставки
- Отправка уведомлений на email
- Аналитика доходов и расходов
- Прогноз баланса

Запуск:

1. Установить Go 1.23+, PostgreSQL 17+
2. Создать базу bankapi и пользователя bankuser
3. Выполнить миграции из папки migrations
4. Сгенерировать PGP-ключи: go run scripts/gen_keys.go > pgp_keys.txt
5. Установить переменные окружения:
   - PGP_PUBLIC_KEY (из pgp_keys.txt)
   - PGP_PRIVATE_KEY (из pgp_keys.txt, опционально)
   - JWT_SECRET, HMAC_SECRET, SMTP_*
6. Запустить: go run main.go

Сервер будет доступен на порту 8080

API (публичные):
POST /register - регистрация (username, email, password)
POST /login - вход (username, password)

API (защищенные, нужен Bearer токен):
POST /accounts - создать счет (currency)
GET /accounts - список счетов
POST /deposit - пополнить (account_id, amount)
POST /transfer - перевод (from_account_id, to_account_id, amount)
POST /cards - выпустить карту (account_id)
GET /cards - список карт
POST /credits - оформить кредит (amount, term_months)
GET /credits - список кредитов
GET /credits/{id}/schedule - график платежей
GET /analytics - аналитика
GET /accounts/{id}/predict - прогноз баланса (?days=30)

Безопасность:
- Пароли: bcrypt
- CVV карт: bcrypt
- Номера карт: PGP-шифрование
- Целостность: HMAC-SHA256
- Доступ к чужим счетам запрещен

Автор: Timur Nurzhau
Bank API

REST API для банковского сервиса на Go с поддержкой JWT-аутентификации, выпуском виртуальных карт (PGP шифрование), кредитованием и аналитикой.

О ПРОЕКТЕ

Bank API - учебный проект, демонстрирующий создание безопасного банковского сервиса на Go. Проект реализует полный цикл операций: от регистрации пользователей до кредитования с интеграцией внешних сервисов (ЦБ РФ, SMTP).

Цель проекта: закрепить навыки разработки REST API на Go, работы с PostgreSQL, реализации безопасности (JWT, PGP, bcrypt), интеграции с внешними API и написания тестов.

ФУНКЦИОНАЛЬНОСТЬ

Core функционал:

Регистрация и аутентификация пользователей (JWT, 24 часа)

Создание банковских счетов (RUB, USD, EUR)

Пополнение баланса и переводы между счетами

История транзакций

Карты:

Выпуск виртуальных карт (Visa, 16 цифр, алгоритм Луна)

PGP-шифрование номера и срока действия карты

HMAC-SHA256 для проверки целостности

Хеширование CVV (bcrypt)

Маскирование номера карты при просмотре

Кредиты:

Оформление кредита с аннуитетными платежами

Автоматический расчёт на основе ключевой ставки ЦБ РФ

Генерация графика платежей

Досрочное погашение

Начисление штрафов за просрочку (+10%)

Аналитика:

Статистика доходов и расходов за месяц

Кредитная нагрузка

Прогноз баланса на N дней (с учётом будущих платежей)

Интеграции:

Центральный банк РФ - получение ключевой ставки (SOAP)

Email уведомления о платежах (SMTP, HTML-письма)

Фоновые задачи:

Шедулер (каждые 12 часов) для списания просроченных платежей

Автоматическое начисление штрафов

Email-напоминания о просрочках

ТЕХНОЛОГИИ

Язык: Go 1.25
Маршрутизация: gorilla/mux
База данных: PostgreSQL 17
Драйвер БД: lib/pq
Аутентификация: JWT (golang-jwt/jwt/v5)
Логирование: logrus
Валидация: go-playground/validator/v10
Шифрование: bcrypt, PGP (gopenpgp/v2), HMAC-SHA256
Email: go-mail/mail/v2
XML парсинг: beevik/etree
Тестирование: стандартный testing пакет

АРХИТЕКТУРА

Проект построен на чистой архитектуре с разделением на слои:

HTTP Layer (Handlers, Middleware, Router)
|
v
Service Layer (Auth, Account, Transfer, Credit, Card, Analytics, CBR, Email)
|
v
Repository Layer (User, Account, Card, Credit, Transaction)
|
v
Database Layer (PostgreSQL 17)

Принципы:

Инъекция зависимостей

Транзакционная целостность

Параметризованные SQL-запросы

Централизованная обработка ошибок

БЕЗОПАСНОСТЬ

Реализованные меры:

Угроза: Перехват пароля
Решение: bcrypt хеширование (cost=10)

Угроза: Кража данных карт
Решение: PGP-шифрование (RSA 2048)

Угроза: Подмена данных карты
Решение: HMAC-SHA256

Угроза: Перехват CVV
Решение: bcrypt хеширование

Угроза: Несанкционированный доступ
Решение: JWT с 24h TTL

Угроза: Доступ к чужим счетам
Решение: Проверка на уровне репозиториев (JOIN с user_id)

Угроза: SQL-инъекции
Решение: Параметризованные запросы

Угроза: Перехват JWT
Решение: HTTPS рекомендуется (не реализован в dev-версии)

Схема шифрования карт:

Card Data (4111****) -> PGP Encrypt (RSA 2048) -> Encrypted Base64
|
v
HMAC Compute
|
v
HMAC-SHA256 Integrity

ИЗВЕСТНЫЕ ОГРАНИЧЕНИЯ И ПУТИ РЕШЕНИЯ

Ограничение 1: Только RUB (для кредитов), частичная поддержка USD/EUR
Причина: Упрощение логики конвертации
Решение: Добавить слой конвертации валют через внешнее API

Ограничение 2: PGP-ключ хранится в переменной окружения
Причина: Упрощение деплоя
Решение: Интеграция с HashiCorp Vault или AWS KMS

Ограничение 3: Максимальный прогноз баланса - 365 дней
Причина: Требование ТЗ
Решение: Можно расширить до любого периода

Ограничение 4: Нет rate limiting
Причина: Не было в ТЗ
Решение: Добавить golang.org/x/time/rate

Ограничение 5: Нет кэширования
Причина: Упрощение архитектуры
Решение: Redis для курсов валют и ключевой ставки

Ограничение 6: Нет HTTPS в dev-режиме
Причина: Упрощение локальной разработки
Решение: Добавить поддержку TLS для production

Ограничение 7: Нет пагинации в списках
Причина: Упрощение реализации
Решение: Добавить limit/offset параметры

Ограничение 8: Нет WebSocket уведомлений
Причина: Требования ТЗ
Решение: Добавить SSE или WebSocket для real-time уведомлений

Ограничение 9: Нет 2FA
Причина: Дополнительная функциональность (не обязательная)
Решение: Добавить TOTP (Google Authenticator)

Ограничение 10: Нет административной панели
Причина: Дополнительная функциональность (не обязательная)
Решение: Реализовать отдельный API для администраторов

УСТАНОВКА И ЗАПУСК

Требования:

Go 1.25+

PostgreSQL 17+

SMTP сервер (для email уведомлений, опционально)

Шаги:

Клонировать репозиторий
git clone https://github.com/ваш-username/bank-api.git
cd bank-api

Установить зависимости
go mod download

Создать базу данных PostgreSQL
createdb bankapi
createuser bankuser
psql -c "ALTER USER bankuser WITH PASSWORD 'bankuser_pass_2024'"

Выполнить миграции
psql -U bankuser -d bankapi -f migrations/001_init.sql

Сгенерировать PGP-ключи (опционально)
go run scripts/gen_keys.go > pgp_keys.txt
export PGP_PUBLIC_KEY="$(cat pgp_keys.txt | grep -A 100 'PUBLIC KEY' | tail -n +2)"

Настроить переменные окружения (создать файл .env)
DB_HOST=localhost
DB_PORT=5432
DB_USER=bankuser
DB_PASSWORD=bankuser_pass_2024
DB_NAME=bankapi
DB_SSLMODE=disable
JWT_SECRET=your-secret-key-min-32-chars
HMAC_SECRET=your-hmac-secret-key
SERVER_PORT=8080
LOG_LEVEL=debug
SMTP_HOST=smtp.example.com (опционально)
SMTP_PORT=587 (опционально)
SMTP_USER=noreply@example.com (опционально)
SMTP_PASS=your-password (опционально)
PGP_PUBLIC_KEY= (опционально)

Запустить сервер
go run main.go

Сервер будет доступен на http://localhost:8080

API ENDPOINTS

Публичные (не требуют токена):

POST /register - регистрация пользователя
Body: {"username":"user","email":"user@example.com","password":"123456"}

POST /login - вход
Body: {"username":"user","password":"123456"}
Response: {"token":"jwt-token","user":{...}}

Защищенные (требуют Bearer токен в заголовке Authorization):

Счета:
POST /accounts - создать счёт
Body: {"currency":"RUB"}
GET /accounts - список счетов

Операции:
POST /transfer - перевод
Body: {"from_account_id":1,"to_account_id":2,"amount":100,"description":"test"}
POST /deposit - пополнение
Body: {"account_id":1,"amount":100}

Карты:
POST /cards - выпустить карту
Body: {"account_id":1}
GET /cards - список карт
POST /cards/pay - оплата картой
Body: {"card_id":1,"amount":100}

Кредиты:
POST /credits - оформить кредит
Body: {"amount":100000,"term_months":12}
GET /credits - список кредитов
GET /credits/{creditId}/schedule - график платежей
POST /credits/{creditId}/repay - досрочное погашение
Body: {"amount":5000}

Аналитика:
GET /analytics - статистика за месяц
GET /accounts/{accountId}/predict?days=30 - прогноз баланса

ТЕСТИРОВАНИЕ

Запуск всех тестов:
go test ./... -v

Отдельные тесты:
go test ./utils -v -run TestGenerateCVV
go test ./utils -v -run TestPGPEncryption

Покрытие кода:
go test ./... -cover

CI/CD

Проект настроен на автоматическую сборку и тестирование через GitHub Actions.

Файл конфигурации: .github/workflows/ci.yml

Что происходит при push в main:

Установка Go 1.25

Загрузка зависимостей

Запуск PostgreSQL в контейнере

Выполнение миграций

Сборка проекта

Запуск всех тестов

Для работы CI необходимо добавить secrets в GitHub:

JWT_SECRET

HMAC_SECRET

ПЛАНЫ ПО УЛУЧШЕНИЮ

Приоритет высокий:

Добавить пагинацию для всех списков (limit/offset)

Реализовать rate limiting

Добавить поддержку HTTPS (TLS)

Написать интеграционные тесты с testcontainers

Приоритет средний:

Добавить кэширование курсов валют (Redis)

Реализовать конвертацию валют при переводах

Добавить WebSocket для real-time уведомлений

Создать документацию Swagger/OpenAPI

Приоритет низкий:

Добавить 2FA (TOTP)

Реализовать административную панель

Добавить поддержку других языков (i18n)

Написать нагрузочное тестирование (pprof)

АВТОР

Timur Nurzhau
# Онлайн магазин аренды

## Описание проекта
Этот проект представляет собой онлайн-магазин аренды, построенный на основе микросервисной архитектуры. Он позволяет пользователям добавлять вещи для аренды, бронировать вещи других пользователей и получать уведомления о статусе аренды.

## Цель проекта
Основная цель проекта — практика создания микросервисной архитектуры с использованием паттернов **API Gateway**, **разделения базы данных** для каждого микросервиса, а также **работы с Kafka** как брокером сообщений. Проект направлен на изучение межсервисной коммуникации, интеграции Kafka для асинхронного обмена данными между сервисами и использования **Nginx** в качестве API-шлюза. Применение транзакций и уровней блокировок для обеспечения атамарности, корректной работе в многопоточном приложении.

## Структура проекта
Проект состоит из следующих микросервисов:
- **Сервис авторизации и регистрации**: обеспечивает регистрацию и аутентификацию пользователей, управляет доступом с использованием токенов. Использует свою собственную базу данных.
- **Сервис отправки уведомлений**: отвечает за отправку уведомлений пользователям о статусе аренды (например, напоминания об окончании аренды). Использует отдельную базу данных для хранения истории уведомлений.
- **Сервис основной логики**: реализует функции приложения, такие как добавление вещей, аренда, управление бронированиями и другие операции. Данные этого сервиса также хранятся в отдельной базе данных.
- **API-шлюз**: реализован на базе Nginx и выполняет маршрутизацию запросов между микросервисами, обеспечивая централизованный доступ к функциональности приложения.

## Технологии/Паттерны
- **Nginx**: используется как API-шлюз для маршрутизации запросов.
- **Kafka**: используется как брокер сообщений для обеспечения асинхронного взаимодействия и обмена данными между микросервисами.
- **Раздельные базы данных**: каждый микросервис имеет собственную базу данных, что позволяет достичь независимости сервисов и лучшего управления данными.
- **Транзакции и уровни блокировок**: для корректной работы приложения используются транзакции и различные уровни блокировок, что позволяет обеспечить целостность данных и минимизировать конфликтные ситуации при параллельной обработке запросов.
- Для авторизации выбраны **Cookie**, так как клиентом является браузер. В куки сохраняются:
   - **Токен аутентификации** — для подтверждения личности пользователя
   - **Email** — почтовый адрес пользователя
   - **Nickname** — псевдоним пользователя


# Микросервис регистрации/авторизации

Этот микросервис отвечает за регистрацию пользователей, аутентификацию и управление сессиями с использованием куков для аутентификации. Архитектура микросервиса является слоёной, что обеспечивает чистоту кода и лёгкость в сопровождении.

## Структура проекта

- **/cmd** - точка входа для приложения
- **/internal** - бизнес-логика и основные компоненты микросервиса
  - **/service** - бизнес-логика
  - **/repository** - взаимодействие с базой данных
  - **/controller** - эндпоинты и обработка данных
  - **/models** - модели данных

## API

### 1. Регистрация пользователя
- **Метод:** `POST`
- **URL:** `/signup`
- **Тело запроса:**
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword",
    "nickname": "Alice123"
  }
- **Транзакции:** В этом эндпоинте используется транзакция для обеспечения атомарности операции.

### 2. Авторизация пользователя
- **Метод:** `POST`
- **URL:** `/signin`
- **Тело запроса:**
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword",
    "nickname": "Alice123"
  }

### 3. Разлогинивание пользователя
- **Метод:** `POST`
- **URL:** `/logout`

### 4.Подтверждение почты пользователя
- **Метод:** `POST`
- **URL:** `/verify-email`
- **Параметры запроса:** - `token`: Токен подтверждения почты, который пользователь получает на свой email.

### 5. Аутенитификация пользователя
- **Метод:** `POST`
- **URL:** `/check-auth`
- **Транзакции:** В этом эндпоинте используется транзакция для обеспечения атомарности операции.
  
# Микросервис Notifier

`Notifier` — это микросервис, отвечающий за отправку уведомлений пользователям. Он является частью проекта **rental-store** и обрабатывает задачи по отправке уведомлений, возникающих по различным событиям, таким как окончание срока аренды.

## Основные функции

- Отправка уведомлений пользователям о предстоящих событиях, завершении аренды и других обновлениях.
- Поддержка нескольких каналов уведомлений (например, email, SMS, push-уведомления).
- Интеграция с другими микросервисами проекта **rental-store**.

## Структура проекта

- **/cmd** - Точка входа для приложения. Описывается логика получения/обработки сообщений.
- **/internal** - Содержит логику отправки сообщений

## Описание функциональности

Микросервис является консюмером в **Kafka**, который подписан на топики из которых получает сообщения от других микросервисов, уведомления отправляет на почту.

# Микросервис Shop

Микросервис **shop** реализует основную логику онлайн-магазина аренды, позволяя добавлять, удалять и арендовать товары, а также просматривать списки всех доступных товаров, товаров для аренды и товаров для покупки.

## Структура проекта

- **/cmd** - точка входа для запуска сервиса
- **/controller** - обработчики запросов (эндпоинты API)
- **/service** - бизнес-логика микросервиса
- **/repository** - взаимодействие с базой данных

## API Эндпоинты

### 1. Добавление товара
- **Метод:** `POST`
- **URL:** `/add-thing`
- **Описание:** Добавляет новый товар в каталог.
- **Тело запроса:**
  ```json
  {
    "owner": "Alice",
    "type": "Тип товара",
    "description": "Описание товара",
    "price": 100.0,
  }

### 2. Удаления товара
- **Метод:** `POST`
- **URL:** `/remuve-thing`
- **Описание:** Удаляет товар из каталога.
- **Тело запроса:**
  ```json
  {
    "thing_id":"Уникальные идентификатор вещи"
  }

### 3. Аренда товара
- **Метод:** `POST`
- **URL:** `/buy-thing`
- **Описание:** Арендует товар из каталога.
- **Тело запроса:**
  ```json
  {
    "thing_id":"Уникальные идентификатор вещи",
    "time_interval":{
         "months":1,
         "days": 13,
         "hours": 12
     }
  }
- **Транзакции:** В этом эндпоинте используется транзакция для обеспечения атомарности операции.

### 3. Показать товар
- **Метод:** `GET`
- **URL:** `/show-all-things`
- **Описание:** Возвращает список всех товаров.

### 4. Показать арендованный товара
- **Метод:** `GET`
- **URL:** `/show-rental-things`
- **Описание:** Показывает все товары арендованные пользователем.

### 5. Показать товар выставленный на аренду
- **Метод:** `GET`
- **URL:** `/show-sale-things`
- **Описание:** Показывает все товары выставленные на аренду пользователем.

## Установка и запуск
1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/Dima205502/rental-store.git

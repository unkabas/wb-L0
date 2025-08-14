Сервис заказов WB-L0
Это мой проект для работы с заказами. Бэкенд на Go, использует Kafka и PostgreSQL, а фронт на React. Бэкенд принимает заказы из Kafka, сохраняет их в базу, кэширует (до 10 заказов) и отдаёт через API. Фронт позволяет искать заказы по Order_UID и показывает краткую инфу.
Как устроен проект
В репе две папки:

wb-l0-backend — бэкенд на Go:
Принимает сообщения из Kafka, сохраняет валидные заказы в базу и кэш.
API GET /order/<order_uid> отдаёт заказ в JSON.
Кэш на 10 заказов.
В папке cmd/:
producer — отправляет 3 валидных JSON заказа в Kafka.
emulator — генерит 10 заказов (50% валидных, 50% невалидных JSON).




wb-l0-front — фронт на React:
Поиск заказа по Order_UID, показывает имя клиента, сумму, товар и город доставки.


Как запустить
1. Склонировать репу
```
git clone <ссылка_на_репу>
cd wb-l0
```

3. Настроить бэкенд
Настроить .env

Зайди в папку бэкенда:cd wb-l0-backend


Создай файл .env с таким содержимым:DB_URL="host=хост user=юзер password=твой_пароль dbname=имя_базы port=5432 sslmode=disable"

Запустить сервисы через Docker

В папке wb-l0-backend выполни:docker-compose up -d

Это поднимет:
3 брокера Kafka.
Zookeeper.
Kafka UI.
PostgreSQL.


Запустить бэкенд

Запусти сервер:go run main.go -m (для автомиграции)

Сервер работает на http://localhost:8080.
API GET /order/<order_uid> отдаёт заказ в JSON.
Заказы кэшируются (до 10), если не в кэше — тянет из базы.
Невалидные сообщения из Kafka игнорируются, валидные сохраняются.



Протестировать с продюсером или эмулятором

Продюсер:
Отправляет 3 валидных заказа в Kafka:go run cmd/producer/main.go




Эмулятор:
Генерит 10 заказов (половина валидные, половина мусор):go run cmd/emulator/main.go


Валидные заказы сохраняются в базу и кэш.
Невалидные логируются и пропускаются.



3. Запустить фронт

Зайди в папку фронта:
cd wb-l0-front


Установи зависимости:
npm install


Запусти фронт:
npm run dev


Откроется на http://localhost:5173.


Как пользоваться фронтом:

Введи Order_UID (возьми из логов эмулятора, например, 123e4567-...).
Нажми "Найти".
Показывает:
Имя клиента (delivery.name).
Сумму (payment.amount).
Название товара (items.name).
Город доставки (delivery.city).


Если заказ не найден, вылезет ошибка.

Детали API

Эндпоинт: GET http://localhost:8080/order/<order_uid>
Ответ: JSON с данными заказа (order_uid, delivery, payment, items).
CORS: Настроен (Access-Control-Allow-Origin: *), фронт с localhost:5173 работает без проблем.
Скорость:
Из кэша: ~3 мс.
Из базы: ~30 мс (кэшируется после первого запроса).



Если что-то не работает

Логи бэкенда:
Смотри логи: "Успешно сохранён заказ" (успех), "Ошибка обработки сообщения" (невалидный JSON).



Проблемы с базой:
Проверь, что PostgreSQL работает (docker ps) и DB_URL правильный.
go run main.go -m




Проблемы с Kafka:
Убедись, что брокеры запущены (docker ps).
Сбрось топик, если надо:kafka-topics.sh --bootstrap-server localhost:9092 --delete --topic order
kafka-topics.sh --bootstrap-server localhost:9092 --create --topic order --partitions 1 --replication-factor 1




Проблемы с фронтом:
Открой DevTools (F12 → Network), проверь запросы к http://localhost:8080.
Убедись, что бэкенд запущен.


Как всё проверить

Запусти docker-compose up -d.
Запустить бэк с автомиграцией: go run main.go -m.
Отправь заказы: go run cmd/emulator/main.go.
Запусти фронт: cd wb-l0-front; npm install; npm run dev.
Открой http://localhost:5173, введи Order_UID из логов эмулятора, проверь результат.

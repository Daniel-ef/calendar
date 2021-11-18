### Example commands
Start the containers: `docker-compose up -d`

Бывает, что go-контейнер стартует до поднятия базы.
Поэтому нужно вручную перезапустить контейнер 
`docker-compose stop go && docker-compose run -p 55443:55443 go`

Rebuild and start: `docker-compose up -d --build`


Важные файлы:
- postgresql/migrations - схема
- postgresql/queries - sql-запросы
- swagger - API
- views - ручки 
- restapi/*_test.go - тесты на ручки. В таком месте, т.к. костыль, чтобы запускать сервер 
  и использовать боевую БД в контейнере и не дружить go-swagger с моком постгреса
- Остальное сгенерено go-swagger

Реализовано:
* создать пользователя;
* создать встречу в календаре пользователя со списком приглашенных пользователей;
* получить детали встречи;
* принять или отклонить приглашение другого пользователя;
* найти все встречи пользователя для заданного промежутка времени;
* для заданного списка пользователей и минимальной продолжительности встречи, найти ближайшей интервал времени, в котором все эти пользователи свободны.
* У встреч в календаре должна быть возможна настройка повторов. В повторах нужно поддержать все возможности, доступные в Google-календаре, кроме Сustom.
* поддержка видимости встреч (если встреча приватная, другие пользователи могут получить только информацию о занятости пользователя, но не детали встречи);
* настройки нотификации пользователя перед встречей - только информаци

Notes:
* Swagger схема в одном файле, т.к. на Маке не работает external $ref :(
* Хорошо бы в 4xx/500 ответах ручки возвращать причину
* Для удобства можно передавать user_id и event_id при создании объектов
* workday_start, workday_end пока тоже не учитывается

Примеры запросов:

`curl --location --request POST 'localhost:55443/users/create' --header 'Content-Type: application/json' --data-raw '{
"user_id": "user_id1",
"email": "abc@mail.ru",
"phone": "+790312331222"
}'`

200

`curl --location --request POST 'localhost:55443/users/create' --header 'Content-Type: application/json' --data-raw '{
"user_id": "user_id2",
"email": "abc@mail.ru",
"phone": "+790312331222"
}'`

500 (duplication email or phone)

`curl --location --request POST 'localhost:55443/users/create' --header 'Content-Type: application/json' --data-raw '{
"user_id": "user_id2",
"email": "cde@mail.ru",
"phone": "+790312331223"
}'`

200

`curl --location --request POST 'localhost:55443/users/create' --header 'Content-Type: application/json' --data-raw '{
"user_id": "user_id3",
"email": "bcd@mail.ru",
"phone": "+790312331224",
"first_name": "Kot",
"last_name": "Matroskin",
"day_start": "12:00",
"day_end": "20:00"
}'`

200

`curl --location --request GET 'localhost:55443/users/info?user_id=user_id3'`

```json
{
  "email":"bcd@mail.ru",
  "first_name":"Kot",
  "last_name":"Matroskin",
  "phone":"+790312331224",
  "user_id":"user_id3"
}
```

`curl --location --request POST 'localhost:55443/event/create' --header 'Content-Type: application/json' --data-raw '{
"event_id": "event_id1",
"name": "Interview1",
"creator": "user_id1",
"time_start": "2021-10-08T10:30:00.000Z",
"time_end": "2021-10-08T12:00:00.000Z",
"visibility": "all"
}'`

200

`curl --location --request POST 'localhost:55443/event/create' --header 'Content-Type: application/json' --data-raw '{
"event_id": "event_id2",
"name": "Interview2",
"creator": "user_id2",
"time_start": "2021-10-08T12:30:00.000Z",
"time_end": "2021-10-08T14:00:00.000Z",
"visibility": "all"
}'`

200

`curl --location --request POST 'localhost:55443/event/room/create' --header 'Content-Type: application/json' --data-raw '{
"room_id": "shmeshariki_1",
"name": "Smeshariki"
}'`

200

`curl --location --request POST 'localhost:55443/event/create' --header 'Content-Type: application/json' --data-raw '{
"event_id": "event_id3",
"name": "Interview3",
"creator": "user_id3",
"notifications": [
{
"before_start": 60,
"step": "m",
"method": "telegram"
},
{
"before_start": 12,
"step": "h",
"method": "sms"
}
],
"time_start": "2020-08-01T18:22:44",
"time_end": "2020-08-01T19:22:44",
"visibility": "all",
"participants": [
{"user_id": "user_id1"},
{"user_id": "user_id2"}
],
"meeting_room": "shmeshariki_1",
"meeting_link": "zoom.us"
}'`

200

`curl --location --request POST 'localhost:55443/invitation/update' --header 'Content-Type: application/json' --data-raw '{
"user_id":  "user_id2",
"event_id": "event_id3",
"accepted": "maybe"
}'`

200

`curl --location --request GET 'localhost:55443/event/info?event_id=event_id3'`

```json
{
"creator":"user_id3",
"name":"Interview3",
"notifications":[
  {
    "before_start":60,
    "method":"telegram",
    "step":"m"
  },
  {
    "before_start":12,
    "method":"sms",
    "step":"h"
  }
],
"participants":[
  {
    "user_id":"user_id3"
  },
  {
    "user_id":"user_id1"
  },
  {
    "accepted":"maybe",
    "user_id":"user_id2"
  }
],
"time_end":"2020-08-01T19:22:44.000Z",
"time_start":"2020-08-01T18:22:44.000Z",
"visibility":"all"
}
```

`curl --location --request GET 'localhost:55443/user_events?user_id=user_id1&time_start=2020-08-01T19:22:00&time_end=2020-08-01T21:22:00'`

```json
{"event_ids":["event_id3"]}
```

`curl --location --request POST 'localhost:55443/users/free_slot' --header 'Content-Type: application/json' --data-raw '{
"user_ids": ["user_id1", "user_id2", "user_id3"],
"slot_interval_min": 30,
"from": "2021-10-08T10:15:00.000Z"
}'`

`{
"time_start": "2021-10-08T12:00:00.000Z",
"time_end": "2021-10-08T12:30:00.000Z"
}`
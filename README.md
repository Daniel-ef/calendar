### Example commands
Start the containers: `docker-compose up -d`

Rebuild and start: `docker-compose up -d --build`

Login: `docker exec -it postgres psql -U user`

This command will be determined by the docker-compose.yml file if different env variables are used.

Eg. `docker exec -it postgres psql -U user testdb` as this takes a database name.

Show tables: `\dt`

Show databases: `\l`

Сервис должен иметь HTTP API, позволяющее:
* создать пользователя;
* создать встречу в календаре пользователя со списком приглашенных пользователей;
* получить детали встречи;
* принять или отклонить приглашение другого пользователя;
* найти все встречи пользователя для заданного промежутка времени;
* для заданного списка пользователей и минимальной продолжительности встречи, найти ближайшей интервал времени, в котором все эти пользователи свободны.
  У встреч в календаре должна быть возможна настройка повторов. В повторах нужно поддержать все возможности, доступные в Google-календаре, кроме Сustom.

Будет плюсом, если вы также реализуете одну или несколько из следующих функций:
* аутентификация пользователя;
* поддержка видимости встреч (если встреча приватная, другие пользователи могут получить только информацию о занятости пользователя, но не детали встречи);
* настройки часового пояса пользователя и его рабочего времени, использование этих настроек для поиска интервала времени, в котором участники свободны;
* настройки нотификации пользователя перед встречей (саму нотификацию достаточно реализовать записью в лог);
* поддержка Custom повторов, как в Google-календаре;
* другие функции, которые кажутся вам полезными в календаре.



```
CREATE TABLE COMPANY (ID INT PRIMARY KEY NOT NULL, NAME text);
INSERT INTO  company(id,name) values (1,'test');
SELECT * from company;
```
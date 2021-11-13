DROP SCHEMA IF EXISTS calendar CASCADE;

CREATE SCHEMA calendar;

SET search_path TO calendar;

CREATE TYPE calendar.visibility_t AS ENUM('all', 'participants');

create TYPE calendar.notification_type_t AS ENUM('email', 'sms', 'telegram');

CREATE TYPE calendar.notification_t AS(
    before_start interval,
    notif_type calendar.notification_type_t
);

CREATE TABLE calendar.users(
    user_id bigserial PRIMARY KEY,
    email text,
    phone text,
    first_name text,
    second_name text,
    day_start time NOT NULL,
    day_end time NOT NULL,


    CHECK(email is not null or phone is not null)
);

CREATE TABLE calendar.meeting_rooms(
                                       id bigserial PRIMARY KEY,
                                       name text NOT NULL
);

CREATE TABLE calendar.meetings(
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    description text,
    creator bigserial REFERENCES users,
    time_start timestamptz NOT NULL,
    time_end timestamptz NOT NULL,
    meeting_room bigserial REFERENCES meeting_rooms,
    notifications notification_t[],
    connection_url text
);


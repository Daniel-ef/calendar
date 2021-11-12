DROP SCHEMA IF EXISTS calendar CASCADE;

CREATE SCHEMA calendar;

CREATE TYPE calendar.visibility_t AS ENUM('all', 'participants');

create TYPE calendar.notification_type_t as ENUM('email', 'sms', 'telegram');

CREATE TYPE calendar.notification_t(
    before interval,
    type notification_type_t
)

CREATE TABLE calendar.users(
    user_id bigserial PRIMARY KEY,
    email text,
    phone text,
    first_name text,
    second_name text,
    day_start time IS NOT NULL,
    day_end time IS NOT NULL,


    CHECK(email is not null or phone is not null)
);

CREATE TABLE calendar.meetings(
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    description text,
    creator bigserial REFERENCES users,
    start timestamptz NOT NULL,
    end timestamptz NOT NULL,
    meeting_room bigserial REFERENCES meeting_rooms,
    notifications array(notification_t),
    connection_url text
);

CREATE TABLE calendar.meeting_rooms(
    id bigserial PRIMARY KEY,
    name text NOT NULL,
);

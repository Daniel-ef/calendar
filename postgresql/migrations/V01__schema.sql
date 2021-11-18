DROP SCHEMA IF EXISTS calendar CASCADE;

CREATE SCHEMA calendar;

SET search_path TO calendar;

CREATE TYPE calendar.visibility_t AS ENUM('all', 'participants');

create TYPE calendar.notification_method_t AS ENUM('email', 'sms', 'telegram');

CREATE TYPE calendar.notification_step_t as ENUM('m', 'h', 'd', 'w');

CREATE TYPE calendar.accepted_t AS ENUM('yes', 'no', 'maybe');

CREATE TYPE calendar.repeat_t AS ENUM('day', 'workday', 'week', 'month', 'year');

CREATE TABLE calendar.users(
    user_id text PRIMARY KEY,
    email text NOT NULL UNIQUE,
    phone text NOT NULL UNIQUE,
    first_name text,
    last_name text,
    workday_start time,
    workday_end time
);

CREATE TABLE calendar.event_rooms(
   room_id text PRIMARY KEY,
   name text NOT NULL
);

CREATE TABLE calendar.events(
    event_id text PRIMARY KEY,
    name text NOT NULL,
    description text,
    version bigserial NOT NULL,
    visibility visibility_t NOT NULL DEFAULT 'all',
    creator text REFERENCES users,
    time_start timestamptz NOT NULL,
    time_end timestamptz NOT NULL,
    repeat repeat_t,
    -- whole_day
    event_room text REFERENCES event_rooms,
    event_link text
);
CREATE INDEX events_time_start_idx ON calendar.events(time_start);

CREATE TABLE calendar.notifications (
    event_id text REFERENCES events,
    before_start int NOT NULL,
    step notification_step_t NOT NULL,
    method notification_method_t NOT NULL
    );

CREATE TABLE calendar.invitations (
    event_id text REFERENCES events,
    user_id text REFERENCES users,
    accepted accepted_t
);
CREATE UNIQUE INDEX invitations_event_user_idx ON invitations(event_id, user_id);

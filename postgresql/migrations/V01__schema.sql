DROP SCHEMA IF EXISTS calendar CASCADE;

CREATE SCHEMA calendar;

SET search_path TO calendar;

CREATE TYPE calendar.visibility_t AS ENUM('all', 'participants');

create TYPE calendar.notification_method_t AS ENUM('email', 'sms', 'telegram');

CREATE TYPE calendar.notification_step_t as ENUM('m', 'h', 'd', 'w');

CREATE TYPE calendar.accepted_t AS ENUM('yes', 'no', 'maybe');

CREATE TABLE calendar.users(
    user_id text PRIMARY KEY,
    email text NOT NULL,
    phone text NOT NULL,
    first_name text,
    last_name text,
    day_start time NOT NULL,
    day_end time NOT NULL
);
CREATE UNIQUE INDEX users_email_idx ON calendar.users(email);
CREATE UNIQUE INDEX users_phone_idx ON calendar.users(phone);

CREATE TABLE calendar.meeting_rooms(
   room_id text PRIMARY KEY,
   name text NOT NULL
);

CREATE TABLE calendar.meetings(
    meet_id text PRIMARY KEY,
    name text NOT NULL,
    description text,
    version bigserial NOT NULL,
    visibility visibility_t NOT NULL DEFAULT 'all',
    creator text REFERENCES users,
    time_start timestamptz NOT NULL,
    time_end timestamptz NOT NULL,
    -- whole_day
    meeting_room text REFERENCES meeting_rooms,
    meeting_link text
);
CREATE INDEX meetings_time_start_idx ON calendar.meetings(time_start);

CREATE TABLE calendar.notifications (
    meet_id text REFERENCES meetings,
    before_start int,
    step notification_step_t,
    method notification_method_t
    );

CREATE TABLE calendar.invitations (
    meet_id text REFERENCES meetings,
    user_id text REFERENCES users,
    accepted accepted_t
);

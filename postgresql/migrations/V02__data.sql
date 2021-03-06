-- INSERT INTO calendar.users(
--     user_id, email, phone, first_name, last_name, workday_start, workday_end
-- ) VALUES
--     ('user_id1', 'abc@mail.ru', '+712345678', 'Sobaka', 'Pushok', '12:00', '20:00'),
--     ('user_id2', 'def@gmail.com', '+712345679', null, null, '05:10', '13:42');
--
-- INSERT INTO calendar.event_rooms(room_id, name)
-- VALUES ('123', 'Golang'), ('456', 'C++');
--
-- INSERT INTO calendar.events(
--     event_id, name, description, version, visibility, creator,
--                               time_start, time_end, repeat, event_room, event_link
-- ) VALUES
-- ('event_id1', 'Event1', 'Event description', 0, 'all', 'user_id1',
--     '2020-08-01T18:22:44', '2020-08-01T19:22:44', null, '123', 'zoom.us'),
-- ('event_id2', 'Event2', null, 0, 'participants', 'user_id2',
--     '2020-08-01T13:15:00', '2020-08-01T13:45:00', 'workday', '456', 'zoom.us');
--
-- INSERT INTO calendar.notifications(
--     event_id, before_start, step, method
-- ) VALUES
-- ('event_id1', '15', 'm', 'sms'),
-- ('event_id1', '12', 'h', 'telegram');
--
-- INSERT INTO calendar.invitations(
--     event_id, user_id, accepted
-- ) VALUES
-- ('event_id1', 'user_id2', 'yes'),
-- ('event_id2', 'user_id1', 'maybe');

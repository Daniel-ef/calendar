package queries

const EventRoomInsert = `
INSERT INTO calendar.event_rooms(room_id, name) VALUES($1, $2);`

package queries

const UserEventsSelect = `
SELECT m.event_id, m.time_start, m.time_end, m.repeat
FROM calendar.events as m
INNER JOIN calendar.invitations as i
ON i.event_id = m.event_id
WHERE i.user_id = $1 and ((m.time_end > $2 and m.time_start < $3) or m.repeat is not null);

`

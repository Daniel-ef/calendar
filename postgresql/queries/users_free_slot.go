package queries

const UsersFreeSlotSelect = `
SELECT e.time_start, e.time_end, e.repeat
FROM calendar.events as e
JOIN calendar.invitations as i
ON e.event_id = i.event_id
WHERE i.user_id = ANY($1) and ((e.time_end > $2 and e.time_start < $3) or e.repeat is not null);
`

package queries

const EventCreate = `INSERT INTO calendar.events(
	event_id,
	name,
	description,
	creator,
	time_start,
	time_end,
	repeat,
	visibility,
	event_room,
	event_link
) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING event_id;
`

const NotificationInsert = `
INSERT INTO calendar.notifications(
	SELECT $1 as event_id,
	UNNEST($2::int[]) as before_start,
	UNNEST($3::calendar.notification_step_t[]) as step,
	UNNEST($4::calendar.notification_method_t[]) as method
);`

// Don't do it
const _ = `
SELECT name, description, visibility,creator, 
	time_start, time_end, repeat, event_room, event_link,
	array_agg(JSON_BUILD_OBJECT(
		'user_id', i.user_id, 'accepted', i.accepted::text))::jsonb[] 
	as participants,
	array_agg(JSON_BUILD_OBJECT(
		'before_start', n.before_start, 'step', n.step, 'method', n.method)) as notifications
INNER JOIN calendar.invitations as i 
	ON m.event_id = i.event_id 
FULL JOIN calendar.notifications as n
	ON n.event_id = m.event_id
WHERE m.event_id = $1
GROUP BY m.event_id;
`

const EventSelect = `
SELECT name, description, visibility,creator,
time_start, time_end, repeat, event_room, event_link
FROM calendar.events
WHERE event_id = $1
LIMIT 1;
`

const NotificationsSelect = `
SELECT before_start, step, method 
FROM calendar.notifications
where event_id = $1
`

const InvitationsSelect = `
SELECT user_id, accepted
FROM calendar.invitations
WHERE event_id = $1
`

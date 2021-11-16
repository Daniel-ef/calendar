package queries

const EventCreate = `INSERT INTO calendar.events(
	event_id,
	name,
	description,
	creator,
	time_start,
	time_end,
	event_room,
	event_link
) VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING event_id;
`

const NotificationInsert = `
INSERT INTO calendar.notifications(
	SELECT $1 as event_id,
	UNNEST($2::int[]) as before_start,
	UNNEST($3::calendar.notification_step_t[]) as step,
	UNNEST($4::calendar.notification_method_t[]) as method
);`

const InvitationsInsert = `INSERT INTO calendar.invitations(
	SELECT $1 as event_id, UNNEST($2::text[]) as user_id
);`

const EventInfoSelect = `
SELECT name, description, visibility,creator, 
	time_start, time_end, event_room, event_link,
	array_agg(JSON_BUILD_OBJECT('user_id', i.user_id, 'accepted', i.accepted::text))::jsonb[] as participants,
	array_agg(n.before_start)::int[] as before_starts,
	array_agg(n.step)::text[] as steps, 
	array_agg(n.method)::text[] as methods 
FROM calendar.events as m
INNER JOIN calendar.invitations as i 
ON m.event_id = i.event_id 
JOIN calendar.notifications as n
ON m.event_id = n.event_id
WHERE m.event_id = $1 
GROUP BY m.event_id;
`

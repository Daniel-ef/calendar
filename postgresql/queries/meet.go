package queries

const MeetCreate = `INSERT INTO calendar.meetings(
	meet_id,
	name,
	description,
	creator,
	time_start,
	time_end,
	meeting_room,
	meeting_link
) VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING meet_id;
`

const NotificationCreate = `
INSERT INTO calendar.notifications(
	SELECT $1 as meet_id,
	UNNEST($2::int[]) as before_start,
	UNNEST($3::calendar.notification_step_t[]) as step,
	UNNEST($4::calendar.notification_method_t[]) as method
);`

const InvitationsCreate = `INSERT INTO calendar.invitations(
	SELECT $1 as meet_id, UNNEST($2::text[]) as user_id
);`

const MeetInfo = `
SELECT name, description, visibility,creator, 
	time_start, time_end, meeting_room, meeting_link, 
	array_agg(user_id)::text[] as participants
FROM calendar.meetings as m
INNER JOIN calendar.invitations as i on m.meet_id = i.meet_id
WHERE m.meet_id = $1
GROUP BY m.meet_id;
`

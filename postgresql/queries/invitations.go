package queries

const InvitationsUpdate = `
UPDATE calendar.invitations
SET accepted = $3
WHERE user_id = $1 and event_id = $2
`
const InvitationsInsert = `INSERT INTO calendar.invitations(
	SELECT $1 as event_id, UNNEST($2::text[]) as user_id
);`

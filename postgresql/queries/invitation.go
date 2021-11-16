package queries

const UpdateInvitation = `
INSERT INTO calendar.invitations (
	user_id, event_id, value
) VALUES ($1, $2, $3)
ON CONFLICT (invitations_event_user_idx)
DO UPDATE
SET value = $3
`

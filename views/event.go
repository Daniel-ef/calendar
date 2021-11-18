package views

import (
	"calendar/models"
	"calendar/postgresql/queries"
	"calendar/restapi/operations"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
)

type Notification struct {
	BeforeStart int64  `json:"before_start,omitempty"`
	Step        string `json:"step,omitempty"`
	Method      string `json:"method,omitempty"`
}

func (a *Notification) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

type Participant struct {
	UserId   string `json:"user_id,omitempty"`
	Accepted string `json:"accepted,omitempty"`
}

func (a *Participant) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func GetUserIds(participants []*models.Participant, creator_id string) (userIds []string) {
	userIdsMap := map[string]struct{}{creator_id: {}}
	userIds = append(userIds, creator_id)
	for _, participant := range participants {
		userId := *participant.UserID
		if _, ok := userIdsMap[userId]; !ok {
			userIdsMap[userId] = struct{}{}
			userIds = append(userIds, userId)
		}
	}
	return userIds

}

func MakeParticipants(participants []Participant) []*models.Participant {
	response := []*models.Participant{}
	for i := range participants {
		response = append(response, &models.Participant{UserID: &participants[i].UserId,
			Accepted: models.Accepted(participants[i].Accepted)})
	}
	return response
}

func MakeNotifications(notifications []Notification) []*models.Notification {
	response := []*models.Notification{}
	for i := range notifications {
		rNotification := models.Notification{BeforeStart: &notifications[i].BeforeStart,
			Step: &notifications[i].Step, Method: &notifications[i].Method}
		response = append(response, &rNotification)
	}
	return response
}

func NewEventCreateHandler(dbClient *sqlx.DB) operations.PostEventCreateHandlerFunc {
	return func(params operations.PostEventCreateParams) middleware.Responder {
		event := params.Body
		trx, err := dbClient.Beginx()

		// events
		var repeat, eventRoom *string
		if event.Repeat != "" {
			repeat = &event.Repeat
		}
		if event.EventRoom != "" {
			eventRoom = &event.EventRoom
		}
		rows, err := trx.Query(queries.EventCreate,
			uuid.New().String(),
			event.Name,
			event.Description,
			event.Creator,
			event.TimeStart,
			event.TimeEnd,
			repeat,
			event.Visibility,
			eventRoom,
			event.EventLink,
		)
		if err != nil {
			log.Print("Error while creating event: ", err.Error())
			_ = trx.Rollback()
			return operations.NewPostEventCreateInternalServerError()
		}
		eventId := new(string)
		if !rows.Next() {
			log.Println("Error while fetching event_id: empty rows")
			_ = trx.Rollback()
			return operations.NewPostEventCreateInternalServerError()
		}
		if err := rows.Scan(eventId); err != nil {
			log.Println("Error while scanning event_id: ", err.Error())
			_ = trx.Rollback()
			return operations.NewPostEventCreateInternalServerError()
		}
		_ = rows.Close()

		// Invitations
		userIds := GetUserIds(event.Participants, *event.Creator)

		if _, err = trx.Exec(queries.InvitationsInsert,
			eventId, pq.StringArray(userIds)); err != nil {
			log.Print("Error while creating invitations: ", err.Error())
			return operations.NewPostEventCreateInternalServerError()
		}

		// Notifications
		beforeStarts := []int64{}
		steps := []string{}
		methods := []string{}
		for _, notification := range event.Notifications {
			beforeStarts = append(beforeStarts, *notification.BeforeStart)
			steps = append(steps, *notification.Step)
			methods = append(methods, *notification.Method)
		}
		if _, err = trx.Exec(queries.NotificationInsert,
			eventId, pq.Int64Array(beforeStarts), pq.StringArray(steps),
			pq.StringArray(methods)); err != nil {
			log.Print("Error while creating notifications: ", err.Error())
			return operations.NewPostEventCreateInternalServerError()
		}

		// Commit
		if err = trx.Commit(); err != nil {
			log.Print("Error while committing: ", err.Error())
		}

		return &operations.PostEventCreateOK{Payload: &models.EventCreateResponse{EventID: eventId}}
	}
}

func NewEventInfoHandler(dbClient *sqlx.DB) operations.GetEventInfoHandlerFunc {
	return func(params operations.GetEventInfoParams) middleware.Responder {
		rows, err := dbClient.Query(queries.EventSelect, params.EventID)
		if err != nil {
			log.Print("Error while fetching event: ", err.Error())
			return operations.NewGetEventInfoInternalServerError()
		}
		if !rows.Next() {
			return operations.NewGetEventInfoNotFound()
		}
		response := models.EventInfo{}
		response.Visibility = new(string)
		response.Creator = new(string)
		response.TimeStart = new(strfmt.DateTime)
		response.TimeEnd = new(strfmt.DateTime)

		repeat := sql.NullString{}
		eventRoom := sql.NullString{}
		if err := rows.Scan(&response.Name, &response.Description, &response.Visibility,
			&response.Creator, &response.TimeStart, &response.TimeEnd,
			&repeat, &eventRoom, &response.EventLink,
		); err != nil {
			log.Print("Error while scanning event: ", err.Error())
			return operations.NewGetEventInfoInternalServerError()
		}
		response.Repeat = repeat.String
		response.EventRoom = eventRoom.String
		err = rows.Close()

		rows, err = dbClient.Query(queries.InvitationsSelect, params.EventID)
		if err != nil {
			log.Print("Error while fetching invitations: ", err.Error())
			return operations.NewGetEventInfoInternalServerError()
		}

		for rows.Next() {
			p := models.Participant{}
			p.UserID = new(string)
			accepted := sql.NullString{}
			if err := rows.Scan(p.UserID, &accepted); err != nil {
				log.Print("Error while scanning invitation: ", err.Error())
				break
			}
			p.Accepted = models.Accepted(accepted.String)
			response.Participants = append(response.Participants, &p)
		}

		rows.NextResultSet()
		rows, err = dbClient.Query(queries.NotificationsSelect, params.EventID)
		if err != nil {
			log.Print("Error while fetching notifications: ", err.Error())
			return operations.NewGetEventInfoInternalServerError()
		}

		for rows.Next() {
			n := models.Notification{}
			n.BeforeStart = new(int64)
			n.Step = new(string)
			n.Method = new(string)
			if err := rows.Scan(n.BeforeStart, n.Step, n.Method); err != nil {
				log.Print("Error while scanning notifications: ", err.Error())
				break
			}
			response.Notifications = append(response.Notifications, &n)
		}

		if *response.Visibility == "participants" && params.UserID != nil &&
			!Consists(response.Participants, params.UserID) {
			response.Name = nil
			response.Description = ""
			response.EventRoom = ""
			response.Notifications = []*models.Notification{}
			response.EventLink = ""
		}

		return &operations.GetEventInfoOK{Payload: &response}
	}
}

func Consists(participants []*models.Participant, id *string) bool {
	for _, p := range participants {
		if *id == *p.UserID {
			return true
		}
	}
	return false
}

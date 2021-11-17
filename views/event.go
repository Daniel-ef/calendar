package views

import (
	"calendar/models"
	"calendar/postgresql/queries"
	"calendar/restapi/operations"
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
		rows, err := trx.Query(queries.EventCreate,
			uuid.New().String(),
			event.Name,
			event.Description,
			event.Creator,
			event.TimeStart,
			event.TimeEnd,
			event.EventRoom,
			event.EventLink,
		)
		if err != nil {
			log.Print("Error while creating event: ", err.Error())
			_ = trx.Rollback()
			return operations.NewPostEventCreateBadRequest()
		}
		eventId := new(string)
		if !rows.Next() {
			log.Println("Error while fetching event_id: empty rows")
			_ = trx.Rollback()
			return operations.NewPostEventCreateBadRequest()
		}
		if err := rows.Scan(eventId); err != nil {
			log.Println("Error while scanning event_id: ", err.Error())
			_ = trx.Rollback()
			return operations.NewPostEventCreateBadRequest()
		}
		_ = rows.Close()

		// Invitations
		userIds := GetUserIds(event.Participants, *event.Creator)

		if _, err = trx.Exec(queries.InvitationsInsert,
			eventId, pq.StringArray(userIds)); err != nil {
			log.Print("Error while creating invitations: ", err.Error())
			return operations.NewPostEventCreateBadRequest()
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
			return operations.NewPostEventCreateBadRequest()
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
		eventId := params.EventID
		rows, err := dbClient.Query(queries.EventInfoSelect, eventId)
		if err != nil {
			log.Print("Error while fetching event: ", err.Error())
			return operations.NewGetEventInfoBadRequest()
		}
		if !rows.Next() {
			return operations.NewGetEventInfoNotFound()
		}
		response := models.EventInfo{}
		response.Visibility = new(string)
		response.Creator = new(string)
		response.TimeStart = new(strfmt.DateTime)
		response.TimeEnd = new(strfmt.DateTime)

		participants := []Participant{}
		notifications := []Notification{}
		if err := rows.Scan(&response.Name, &response.Description, &response.Visibility,
			&response.Creator, &response.TimeStart, &response.TimeEnd, &response.EventRoom,
			&response.EventLink, pq.Array(&participants), pq.Array(&notifications),
		); err != nil {
			log.Print("Error while fetching event: ", err.Error())
			return operations.NewGetEventInfoBadRequest()
		}
		response.Participants = MakeParticipants(participants)
		response.Notifications = MakeNotifications(notifications)

		return &operations.GetEventInfoOK{Payload: &response}
	}
}

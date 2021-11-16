package views

import (
	"calendar/models"
	"calendar/postgresql/queries"
	"calendar/restapi/operations"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
)

//type Notification struct {
//	before_start      int64
//	notification_type string
//}
//
//func (n *Notification) Value() (driver.Value, error) {
//	return fmt.Sprintf("('%d','%s')", n.before_start, n.notification_type), nil
//}

type Participant struct {
	UserId   string `json:"user_id,omitempty"`
	Accepted string `json:"accepted,omitempty"`
}

func (a Participant) Value() (driver.Value, error) {
	return json.Marshal(a)
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
	for _, participant := range participants {
		response = append(response, &models.Participant{UserID: &participant.UserId,
			Accepted: models.Accepted(participant.Accepted)})
	}
	return response
}

func MakeNotifications(starts pq.Int64Array, steps pq.StringArray,
	methods pq.StringArray) (ret []*models.Notification) {
	if !(len(starts) == len(steps) && len(steps) == len(methods)) {
		log.Println("Invalid data in notifications")
		return ret
	}
	for i := range starts {
		ret = append(ret, &models.Notification{
			BeforeStart: &starts[i], Method: &steps[i], Step: &methods[i]})
	}
	return
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
		beforeStarts := pq.Int64Array{}
		steps := pq.StringArray{}
		methods := pq.StringArray{}

		participants := []Participant{}
		if err := rows.Scan(&response.Name, &response.Description, &response.Visibility,
			&response.Creator, &response.TimeStart, &response.TimeEnd, &response.EventRoom,
			&response.EventLink, pq.Array(&participants),
			&beforeStarts, &steps, &methods); err != nil {
			log.Print("Error while fetching event: ", err.Error())
			return operations.NewGetEventInfoBadRequest()
		}
		log.Println(participants)
		response.Participants = MakeParticipants(participants)
		response.Notifications = MakeNotifications(beforeStarts, steps, methods)

		return &operations.GetEventInfoOK{Payload: &response}
	}
}

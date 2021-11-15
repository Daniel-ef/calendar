package views

import (
	"calendar/models"
	"calendar/postgresql/queries"
	"calendar/restapi/operations"
	"database/sql/driver"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"time"
)

type Notification struct {
	before_start      int64
	notification_type string
}

type Meeting struct {
	meet_id       string
	name          string
	description   string
	version       uint
	creator       string
	time_start    time.Time
	time_end      time.Time
	meeting_room  string
	notifications []Notification
	meeting_link  string
}

func (n *Notification) Value() (driver.Value, error) {
	return fmt.Sprintf("('%d','%s')", n.before_start, n.notification_type), nil
}

func GetUserIds(participants []string) (userIds []string) {
	userIdsMap := map[string]struct{}{}
	for _, participant := range participants {
		if _, ok := userIdsMap[participant]; !ok {
			userIdsMap[participant] = struct{}{}
			userIds = append(userIds, participant)
		}
	}
	return userIds

}

func RemoveCreator(p []string, c string) []string {
	for i := range p {
		if p[i] == c {
			p[i] = p[len(p)-1]
			return p[:len(p)-1]
		}
	}
	return p
}

func NewMeetCreateHandler(dbClient *sqlx.DB) operations.PostMeetCreateHandlerFunc {
	return func(params operations.PostMeetCreateParams) middleware.Responder {
		meet := params.Body
		trx, err := dbClient.Beginx()

		// Meetings
		rows, err := trx.Query(queries.MeetCreate,
			uuid.New().String(),
			meet.Name,
			meet.Description,
			meet.Creator,
			meet.TimeStart,
			meet.TimeEnd,
			meet.MeetingRoom,
			meet.MeetingLink,
		)
		if err != nil {
			log.Print("Error while creating meeting: ", err.Error())
			_ = trx.Rollback()
			return operations.NewPostMeetCreateBadRequest()
		}
		meetId := new(string)
		if !rows.Next() {
			log.Println("Error while fetching meet_id: empty rows")
			_ = trx.Rollback()
			return operations.NewPostMeetCreateBadRequest()
		}
		if err := rows.Scan(meetId); err != nil {
			log.Println("Error while scanning meet_id: ", err.Error())
			_ = trx.Rollback()
			return operations.NewPostMeetCreateBadRequest()
		}
		_ = rows.Close()

		// Invitations
		userIds := GetUserIds(append(meet.Participants, *meet.Creator))

		if _, err = trx.Exec(queries.InvitationsCreate,
			meetId, pq.StringArray(userIds)); err != nil {
			log.Print("Error while creating invitations: ", err.Error())
			return operations.NewPostMeetCreateBadRequest()
		}

		// Notifications
		beforeStarts := []int64{}
		steps := []string{}
		methods := []string{}
		for _, notification := range meet.Notifications {
			beforeStarts = append(beforeStarts, *notification.BeforeStart)
			steps = append(steps, *notification.Step)
			methods = append(methods, *notification.Method)
		}
		if _, err = trx.Exec(queries.NotificationCreate,
			meetId, pq.Int64Array(beforeStarts), pq.StringArray(steps),
			pq.StringArray(methods)); err != nil {
			log.Print("Error while creating notifications: ", err.Error())
			return operations.NewPostMeetCreateBadRequest()
		}

		// Commit
		if err = trx.Commit(); err != nil {
			log.Print("Error while committing: ", err.Error())
		}

		return &operations.PostMeetCreateOK{Payload: &models.MeetCreateResponse{MeetID: meetId}}
	}
}

func NewMeetInfoHandler(dbClient *sqlx.DB) operations.GetMeetInfoHandlerFunc {
	return func(params operations.GetMeetInfoParams) middleware.Responder {
		meetId := params.MeetingID
		rows, err := dbClient.Query(queries.MeetInfo, meetId)
		if err != nil {
			log.Print("Error while fetching meeting: ", err.Error())
			return operations.NewGetMeetInfoBadRequest()
		}
		if !rows.Next() {
			return operations.NewGetMeetInfoNotFound()
		}
		response := models.MeetInfo{}
		response.Visibility = new(string)
		response.Creator = new(string)
		response.TimeStart = new(strfmt.DateTime)
		response.TimeEnd = new(strfmt.DateTime)
		if err := rows.Scan(&response.Name, &response.Description, response.Visibility,
			response.Creator, response.TimeStart, response.TimeEnd, &response.MeetingRoom,
			&response.MeetingLink, pq.Array(&response.Participants)); err != nil {
			log.Print("Error while fetching meeting: ", err.Error())
			return operations.NewGetMeetInfoBadRequest()
		}
		response.Participants = RemoveCreator(response.Participants, *response.Creator)

		return &operations.GetMeetInfoOK{Payload: &response}
	}
}

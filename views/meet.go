package views

import (
	"calendar/models"
	"calendar/postgresql/queries"
	"calendar/restapi/operations"
	"database/sql/driver"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
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

func makeNotifications(notifications []*models.Notification) (ret []*Notification) {
	for _, notification := range notifications {
		ret = append(ret, &Notification{*notification.BeforeStart, *notification.NotificationType})
	}
	return
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

func NewMeetCreateHandler(dbClient *sqlx.DB) operations.PostMeetCreateHandlerFunc {
	return func(params operations.PostMeetCreateParams) middleware.Responder {
		meet := params.Body
		trx, err := dbClient.Beginx()
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
		if success := rows.Next(); !success {
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

		userIds := GetUserIds(append(meet.Participants, *meet.Creator))

		if _, err = trx.Exec(queries.InvitationsCreate,
			meetId, pq.StringArray(userIds)); err != nil {
			log.Print("Error while creating invitations: ", err.Error())
			return operations.NewPostMeetCreateBadRequest()
		}
		if err = trx.Commit(); err != nil {
			log.Print("Error while committing: ", err.Error())
		}

		// notifications

		return &operations.PostMeetCreateOK{Payload: &models.MeetCreateResponse{MeetID: meetId}}
	}
}

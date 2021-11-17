package views

import (
	"calendar/models"
	"calendar/postgresql/queries"
	"calendar/restapi/operations"
	"calendar/utils"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type EventMeta struct {
	eventId    string
	timeStart  strfmt.DateTime
	timeEnd    strfmt.DateTime
	repeatType string
}

func NewUserEventsHandler(dbClient *sqlx.DB) operations.GetUserEventsHandlerFunc {
	return func(params operations.GetUserEventsParams) middleware.Responder {
		rows, err := dbClient.Query(queries.UserEventsSelect,
			params.UserID, params.TimeStart, params.TimeEnd)
		if err != nil {
			log.Print("Error while fetching user events: ", err.Error())
			return operations.NewGetUserEventsInternalServerError()
		}
		response := models.UserEventsResponse{}
		for rows.Next() {
			var eventMeta EventMeta
			if err := rows.Scan(
				&eventMeta.eventId, &eventMeta.timeStart,
				&eventMeta.timeEnd, &eventMeta.repeatType,
			); err != nil {
				log.Print("Error while scanning user_id: ", err.Error())
				return operations.NewGetUserEventsInternalServerError()
			}
			if utils.CheckEvent(time.Time(eventMeta.timeStart), time.Time(eventMeta.timeEnd),
				time.Time(params.TimeStart), time.Time(params.TimeEnd), eventMeta.repeatType) {
				response.EventIds = append(response.EventIds, eventMeta.eventId)
			}
		}
		return &operations.GetUserEventsOK{Payload: &response}
	}
}

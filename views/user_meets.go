package views

import (
	"calendar/models"
	"calendar/postgresql/queries"
	"calendar/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jmoiron/sqlx"
	"log"
)

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
			var eventId string
			if err := rows.Scan(&eventId); err != nil {
				log.Print("Error while scanning user_id: ", err.Error())
				return operations.NewGetUserEventsInternalServerError()
			}
			response.EventIds = append(response.EventIds, eventId)
		}
		return &operations.GetUserEventsOK{Payload: &response}
	}
}

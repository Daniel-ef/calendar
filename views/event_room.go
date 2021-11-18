package views

import (
	"calendar/postgresql/queries"
	"calendar/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
)

func NewEventRoomCreateHandler(dbClient *sqlx.DB) operations.PostEventRoomCreateHandlerFunc {
	return func(params operations.PostEventRoomCreateParams) middleware.Responder {

		roomId := params.Body.RoomID
		if roomId == "" {
			roomId = uuid.New().String()
		}
		_, err := dbClient.Exec(queries.EventRoomInsert,
			roomId,
			params.Body.Name,
		)
		if err != nil {
			log.Print("Error while creating event room: ", err.Error())
			return operations.NewPostUsersCreateInternalServerError()
		}
		return &operations.PostEventRoomCreateOK{
			Payload: &operations.PostEventRoomCreateOKBody{RoomID: &roomId}}
	}
}

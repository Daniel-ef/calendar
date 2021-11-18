package views

import (
	"calendar/postgresql/queries"
	"calendar/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jmoiron/sqlx"
	"log"
)

func NewInvitationUpdateHandler(dbClient *sqlx.DB) operations.PostInvitationUpdateHandlerFunc {
	return func(params operations.PostInvitationUpdateParams) middleware.Responder {
		rows, err := dbClient.Exec(queries.InvitationsUpdate,
			params.Body.UserID, params.Body.EventID, params.Body.Accepted)
		if err != nil {
			log.Print("Error while updating invitation: ", err.Error())
			return operations.NewPostInvitationUpdateInternalServerError()
		}
		if effected, _ := rows.RowsAffected(); effected == 0 {
			log.Print("Can't find invitation with such params: ", err.Error())
			return operations.NewPostInvitationUpdateInternalServerError()
		}
		return operations.NewPostInvitationUpdateOK()
	}
}

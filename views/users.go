package views

import (
	"calendar/models"
	"calendar/postgresql/queries"
	"calendar/restapi/operations"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
)

type User struct {
	UserId    string `db:"user_id"`
	Email     string
	Phone     string
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	DayStart  string `db:"day_start"`
	DayEnd    string `db:"day_end"`
}

func NewUsersCreateHandler(dbClient *sqlx.DB) operations.PostUsersCreateHandlerFunc {
	return func(params operations.PostUsersCreateParams) middleware.Responder {
		userInfo := params.Body
		rows, err := dbClient.NamedQuery(queries.UserCreate,
			User{
				UserId:    uuid.New().String(),
				Email:     *userInfo.Email,
				Phone:     *userInfo.Phone,
				FirstName: userInfo.FirstName,
				LastName:  userInfo.LastName,
				DayStart:  userInfo.DayStart,
				DayEnd:    userInfo.DayEnd,
			},
		)
		if err != nil {
			log.Print("Error while creating user: ", err.Error())
			return operations.NewPostUsersCreateBadRequest()
		}

		userId := new(string)
		success := rows.Next()
		if !success {
			log.Println("Error while fetching user_id: ", err.Error())
			return operations.NewPostUsersCreateBadRequest()
		}
		if err := rows.Scan(userId); err != nil {
			log.Println("Error while scanning user_id: ", err.Error())
			return operations.NewPostUsersCreateBadRequest()
		}

		return &operations.PostUsersCreateOK{Payload: &models.UsersCreateResponse{userId}}
	}
}

func NewUsersInfoHandler(dbClient *sqlx.DB) operations.GetUsersInfoHandlerFunc {
	return func(params operations.GetUsersInfoParams) middleware.Responder {
		query := queries.UserFind
		if params.Email == nil && params.Phone == nil ||
			params.Email != nil && params.Phone != nil {
			return operations.NewGetUsersInfoBadRequest()
		}
		if params.Email != nil {
			query += fmt.Sprintf(" email='%s'", *params.Email)
		} else {
			query += fmt.Sprintf(" phone='%s'", *params.Phone)
		}

		user := User{}
		if err := dbClient.Get(&user, query); err != nil {
			log.Print("Error while creating user: ", err.Error())
			return operations.NewGetUsersInfoBadRequest()
		}
		retUserInfo := models.UserInfo{
			UserID:    user.UserId,
			Email:     &user.Email,
			Phone:     &user.Phone,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			DayStart:  user.DayStart,
			DayEnd:    user.DayEnd,
		}
		response := operations.NewGetUsersInfoOK()
		response.SetPayload(&retUserInfo)
		return response
	}
}

package views

import (
	"calendar/models"
	"calendar/postgresql/queries"
	"calendar/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
)

type User struct {
	UserId       string `db:"user_id"`
	Email        string
	Phone        string
	FirstName    string  `db:"first_name"`
	LastName     string  `db:"last_name"`
	WorkDayStart *string `db:"workday_start"`
	WorkDayEnd   *string `db:"workday_end"`
}

func NewUsersCreateHandler(dbClient *sqlx.DB) operations.PostUsersCreateHandlerFunc {
	return func(params operations.PostUsersCreateParams) middleware.Responder {
		userInfo := params.Body
		user := User{
			UserId:    uuid.New().String(),
			Email:     *userInfo.Email,
			Phone:     *userInfo.Phone,
			FirstName: userInfo.FirstName,
			LastName:  userInfo.LastName,
		}
		if len(userInfo.WorkdayStart) != 0 {
			user.WorkDayStart = &userInfo.WorkdayStart
		}
		if len(userInfo.WorkdayEnd) != 0 {
			user.WorkDayEnd = &userInfo.WorkdayEnd
		}
		rows, err := dbClient.NamedQuery(queries.UserInsert, user)
		if err != nil {
			log.Print("Error while creating user: ", err.Error())
			return operations.NewPostUsersCreateInternalServerError()
		}

		userId := new(string)
		success := rows.Next()
		if !success {
			log.Println("Error while fetching user_id: ", err.Error())
			return operations.NewPostUsersCreateInternalServerError()
		}
		if err := rows.Scan(userId); err != nil {
			log.Println("Error while scanning user_id: ", err.Error())
			return operations.NewPostUsersCreateInternalServerError()
		}

		return &operations.PostUsersCreateOK{Payload: &models.UsersCreateResponse{userId}}
	}
}

func NewUsersInfoHandler(dbClient *sqlx.DB) operations.GetUsersInfoHandlerFunc {
	return func(params operations.GetUsersInfoParams) middleware.Responder {
		user := User{}
		if err := dbClient.Get(&user, queries.UserSelect, params.UserID); err != nil {
			log.Print("Error while creating user: ", err.Error())
			return operations.NewGetUsersInfoInternalServerError()
		}
		retUserInfo := models.UserInfo{
			UserID:    user.UserId,
			Email:     &user.Email,
			Phone:     &user.Phone,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
		if user.WorkDayStart != nil {
			retUserInfo.WorkdayStart = *user.WorkDayStart
		}
		if user.WorkDayEnd != nil {
			retUserInfo.WorkdayEnd = *user.WorkDayEnd
		}
		response := operations.NewGetUsersInfoOK()
		response.SetPayload(&retUserInfo)
		return response
	}
}

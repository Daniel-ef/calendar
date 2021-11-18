package restapi

import (
	"calendar/postgresql"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"os"
	"testing"
)

func CreateUser(t *testing.T, ts *httptest.Server, email string, phone string) string {
	request := map[string]interface{}{
		"email": email,
		"phone": phone,
	}
	requestB, _ := json.Marshal(request)
	resp, statusCode, _ := MakeResponse("POST", ts.URL+"/users/create", nil,
		requestB)
	assert.Equal(t, statusCode, 200, "")
	return resp["user_id"].(string)
}

func createEvent(t *testing.T, ts *httptest.Server,
	user_id string, start string, end string, repeat string) string {
	request := map[string]interface{}{
		"name":       "Interview1",
		"creator":    user_id,
		"time_start": start,
		"time_end":   end,
		"repeat":     repeat,
		"visibility": "all",
	}
	requestB, _ := json.Marshal(request)
	resp, statusCode, _ := MakeResponse("POST", ts.URL+"/event/create", nil,
		requestB)
	assert.Equal(t, statusCode, 200, "")
	return resp["event_id"].(string)
}

func TestUsersFreeSlot(t *testing.T) {
	handler, err := GetAPIHandler()
	if err != nil {
		t.Fatal("get api handler", err)
	}
	ts := httptest.NewServer(handler)
	defer ts.Close()

	psql := postgresql.GetPostgresClient()
	// Don't do it at home
	dat, err := os.ReadFile("../postgresql/migrations/V01__schema.sql")
	if err != nil {
		panic(err)
	}
	_, err = psql.Exec(string(dat))
	if err != nil {
		panic(err)
	}

	// users
	userId1 := CreateUser(t, ts, "abc@mail.ru", "+71234567")
	userId2 := CreateUser(t, ts, "bcd@mail.ru", "+71234561")
	userId3 := CreateUser(t, ts, "cde@mail.ru", "+71234562")

	// simple
	_ = createEvent(t, ts, userId1,
		"2021-10-08T10:30:00.000Z", "2021-10-08T12:00:00.000Z", "")
	_ = createEvent(t, ts, userId2,
		"2021-10-08T12:30:00.000Z", "2021-10-08T14:00:00.000Z", "")
	_ = createEvent(t, ts, userId3,
		"2021-10-08T13:00:00.000Z", "2021-10-09T12:00:00.000Z", "")

	request := map[string]interface{}{
		"user_ids":          []string{userId1, userId2, userId3},
		"slot_interval_min": 30,
		"from":              "2021-10-08T10:15:00.000Z",
	}
	requestB, _ := json.Marshal(request)
	resp, statusCode, _ := MakeResponse("POST", ts.URL+"/users/free_slot", nil,
		requestB)
	assert.Equal(t, 200, statusCode, "")
	assert.Equal(t, []string{resp["time_start"].(string), resp["time_end"].(string)},
		[]string{"2021-10-08T12:00:00.000Z", "2021-10-08T12:30:00.000Z"})

	// with repeat
	//_ = createEvent(t, ts, userId1,
	//	"2021-10-08T10:30:00.000Z", "2021-10-08T12:00:00.000Z", "day")
	//_ = createEvent(t, ts, userId2,
	//	"2021-10-08T12:30:00.000Z", "2021-10-08T14:00:00.000Z", "week")
	//_ = createEvent(t, ts, userId3,
	//	"2021-10-08T13:00:00.000Z", "2021-10-09T12:00:00.000Z", "workday")
	//
	//request := map[string]interface{}{
	//	"user_ids":          []string{userId1, userId2, userId3},
	//	"slot_interval_min": 30,
	//	"from":              "2021-10-15T10:15:00.000Z",
	//}
	//requestB, _ := json.Marshal(request)
	//resp, statusCode, _ := MakeResponse("POST", ts.URL+"/users/free_slot", nil,
	//	requestB)
	//assert.Equal(t, 200, statusCode, "")
	//assert.Equal(t, []string{resp["time_start"].(string), resp["time_end"].(string)},
	//	[]string{"2021-10-15T12:00:00.000Z", "2021-10-15T12:30:00.000Z"})

}

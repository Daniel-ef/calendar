package restapi

import (
	"calendar/postgresql"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestEvents(t *testing.T) {
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
	body, statusCode, _ := MakeResponse("POST", ts.URL+"/users/create", nil,
		[]byte(`{
			"email": "abc@mail.ru",
			"phone": "+744222134"
		}`))
	assert.Equal(t, statusCode, 200, "")
	userId1 := body["user_id"].(string)
	body, statusCode, _ = MakeResponse("POST", ts.URL+"/users/create", nil,
		[]byte(`{
			"email": "abc1@mail.ru",
			"phone": "+744222132"
		}`))
	assert.Equal(t, statusCode, 200, "")
	userId2 := body["user_id"].(string)

	// test base creation
	request := map[string]interface{}{
		"name":       "Interview1",
		"creator":    userId1,
		"time_start": "2020-08-01T18:22:44.000Z",
		"time_end":   "2020-08-01T19:22:44.000Z",
		"visibility": "all",
	}
	requestB, _ := json.Marshal(request)
	body, statusCode, _ = MakeResponse("POST", ts.URL+"/event/create", nil,
		requestB)
	assert.Equal(t, statusCode, 200, "")

	params := url.Values{
		"event_id": {body["event_id"].(string)},
	}
	body, statusCode, _ = MakeResponse("GET", ts.URL+"/event/info?", &params, nil)
	assert.Equal(t, statusCode, 200, "")
	participants := body["participants"].([]interface{})
	participant := participants[0].(map[string]interface{})
	assert.Equal(t, participant["user_id"].(string), userId1)
	delete(body, "participants")
	delete(body, "notifications")
	assert.Equal(t, body, request)

	// test full creation
	request = map[string]interface{}{
		"name": "Smeshariki",
	}
	requestB, _ = json.Marshal(request)
	body, statusCode, _ = MakeResponse("POST", ts.URL+"/event/room/create", nil,
		requestB)
	assert.Equal(t, statusCode, 200, "")
	roomId := body["room_id"].(string)

	bigRequest := map[string]interface{}{
		"name":        "Interview1",
		"description": "some description",
		"creator":     userId1,
		"participants": []interface{}{
			map[string]interface{}{
				"user_id": userId2,
			},
		},
		"time_start": "2020-08-01T18:22:44.000Z",
		"time_end":   "2020-08-01T19:22:44.000Z",
		"repeat":     "workday",
		"event_room": roomId,
		"event_link": "zoom.us",
		"visibility": "participants",
		"notifications": []interface{}{
			map[string]interface{}{
				"before_start": 60,
				"step":         "m",
				"method":       "telegram",
			},
			map[string]interface{}{
				"before_start": 12,
				"step":         "h",
				"method":       "sms",
			},
		},
	}
	requestB, _ = json.Marshal(bigRequest)
	body, statusCode, _ = MakeResponse("POST", ts.URL+"/event/create", nil,
		requestB)
	assert.Equal(t, statusCode, 200, "")
	eventId := body["event_id"].(string)

	request = map[string]interface{}{
		"user_id":  userId1,
		"event_id": eventId,
		"accepted": "maybe",
	}
	requestB, _ = json.Marshal(request)
	body, statusCode, _ = MakeResponse("POST", ts.URL+"/invitation/update", nil,
		requestB)
	assert.Equal(t, statusCode, 200, "")

	params = url.Values{
		"event_id": {eventId},
	}
	body, statusCode, _ = MakeResponse("GET", ts.URL+"/event/info?", &params, nil)
	assert.Equal(t, statusCode, 200, "")
	participants = body["participants"].([]interface{})
	participant = participants[0].(map[string]interface{})
	assert.Equal(t, participant["user_id"].(string), userId1)
	assert.Equal(t, participant["accepted"].(string), "maybe")
	delete(body, "participants")
	delete(body, "notifications")
	delete(bigRequest, "participants")
	delete(bigRequest, "notifications")

	assert.Equal(t, body, bigRequest)

	params = url.Values{
		"event_id": {eventId},
		"user_id":  {"anyUser"},
	}
	body, statusCode, _ = MakeResponse("GET", ts.URL+"/event/info?", &params, nil)
	assert.Equal(t, statusCode, 200, "")
	participants = body["participants"].([]interface{})
	participant = participants[0].(map[string]interface{})
	assert.Equal(t, participant["user_id"].(string), userId1)
	assert.Equal(t, participant["accepted"].(string), "maybe")
	assert.Equal(t, body["name"], nil)
}

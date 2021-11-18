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

func TestUserEvents(t *testing.T) {
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
	resp, statusCode, _ := MakeResponse("POST", ts.URL+"/users/create", nil,
		[]byte(`{
			"email": "abc@mail.ru",
			"phone": "+744222134"
		}`))
	assert.Equal(t, statusCode, 200, "")
	userId1 := resp["user_id"].(string)

	// events
	request := map[string]interface{}{
		"name":       "Interview1",
		"creator":    userId1,
		"time_start": "2021-10-07T10:30:00.000Z",
		"time_end":   "2021-10-07T13:00:00.000Z",
		"repeat":     "week",
		"visibility": "all",
	}
	requestB, _ := json.Marshal(request)
	resp, statusCode, _ = MakeResponse("POST", ts.URL+"/event/create", nil,
		requestB)
	assert.Equal(t, statusCode, 200, "")
	eventId1 := resp["event_id"].(string)

	request = map[string]interface{}{
		"name":       "Interview1",
		"creator":    userId1,
		"time_start": "2021-10-08T18:30:00.000Z",
		"time_end":   "2021-10-09T12:00:00.000Z",
		"repeat":     "workday",
		"visibility": "all",
	}
	requestB, _ = json.Marshal(request)
	resp, statusCode, _ = MakeResponse("POST", ts.URL+"/event/create", nil,
		requestB)
	assert.Equal(t, statusCode, 200, "")
	eventId2 := resp["event_id"].(string)

	// user_events
	params := url.Values{
		"user_id":    {userId1},
		"time_start": {"2021-10-07T11:00:00Z"},
		"time_end":   {"2021-10-07T12:00:00Z"},
	}
	resp, statusCode, _ = MakeResponse("GET", ts.URL+"/user_events?", &params, nil)
	assert.Equal(t, statusCode, 200, "")
	checkEvents(t, map[string]struct{}{eventId1: {}}, resp["event_ids"].([]interface{}))

	params = url.Values{
		"user_id":    {userId1},
		"time_start": {"2021-10-09T11:00:00Z"},
		"time_end":   {"2021-10-09T13:00:00Z"},
	}
	resp, statusCode, _ = MakeResponse("GET", ts.URL+"/user_events?", &params, nil)
	assert.Equal(t, statusCode, 200, "")
	checkEvents(t, map[string]struct{}{eventId2: {}}, resp["event_ids"].([]interface{}))

	params = url.Values{
		"user_id":    {userId1},
		"time_start": {"2021-10-14T11:00:00Z"},
		"time_end":   {"2021-10-14T19:00:00Z"},
	}
	resp, statusCode, _ = MakeResponse("GET", ts.URL+"/user_events?", &params, nil)
	assert.Equal(t, statusCode, 200, "")
	checkEvents(t, map[string]struct{}{eventId1: {}, eventId2: {}}, resp["event_ids"].([]interface{}))
}

func checkEvents(t *testing.T, m map[string]struct{}, resp []interface{}) {
	for _, el := range resp {
		_, ok := m[el.(string)]
		assert.True(t, ok)
	}
}

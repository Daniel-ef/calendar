package restapi

import (
	"calendar/postgresql"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestUsers(t *testing.T) {
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

	_, statusCode, _ := MakeResponse("POST", ts.URL+"/users/create", nil,
		[]byte(`{
			"email": "abc@mail.ru",
		}`))
	assert.Equal(t, statusCode, 400, "")

	// test creation
	resp, statusCode, _ := MakeResponse("POST", ts.URL+"/users/create", nil,
		[]byte(`{
			"email": "e24@mail.ru",
			"phone": "+744222134"
		}`))
	assert.Equal(t, statusCode, 200, "")
	assert.True(t, len(resp["user_id"].(string)) > 0, "")

	params := url.Values{
		"user_id": {resp["user_id"].(string)},
	}
	resp, statusCode, _ = MakeResponse("GET", ts.URL+"/users/info?", &params, nil)
	assert.Equal(t, statusCode, 200, "")
	assert.Equal(t, resp["email"].(string), "e24@mail.ru")
	assert.Equal(t, resp["phone"].(string), "+744222134")

	// test big creation
	resp, statusCode, _ = MakeResponse("POST", ts.URL+"/users/create", nil,
		[]byte(`{
			"email": "abc@mail.ru",
			"first_name": "Sobaka",
			"last_name": "Pushok",
			"phone": "+712345678",
			"workday_end": "20:00",
			"workday_start": "12:00"
		}`))
	assert.Equal(t, statusCode, 200, "")
	assert.True(t, len(resp["user_id"].(string)) > 0, "")

	params = url.Values{
		"user_id": {resp["user_id"].(string)},
	}
	resp, statusCode, _ = MakeResponse("GET", ts.URL+"/users/info?", &params, nil)
	assert.Equal(t, statusCode, 200, "")
	delete(resp, "user_id")
	assert.Equal(t, resp, map[string]interface{}{
		"email":         "abc@mail.ru",
		"phone":         "+712345678",
		"first_name":    "Sobaka",
		"last_name":     "Pushok",
		"workday_start": "0000-01-01T12:00:00Z",
		"workday_end":   "0000-01-01T20:00:00Z",
	})
}

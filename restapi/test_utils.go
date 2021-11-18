package restapi

import (
	"bytes"
	"calendar/restapi/operations"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/loads"
	"io"
	"net/http"
	"net/url"
)

func getAPI() (*operations.CalendarAPIAPI, error) {
	swaggerSpec, err := loads.Analyzed(SwaggerJSON, "")
	if err != nil {
		return nil, err
	}
	api := operations.NewCalendarAPIAPI(swaggerSpec)
	return api, nil
}

func GetAPIHandler() (http.Handler, error) {
	api, err := getAPI()
	if err != nil {
		return nil, err
	}
	h := configureAPI(api)
	err = api.Validate()
	if err != nil {
		return nil, err
	}
	return h, nil
}

func MakeResponse(method string, url string, params *url.Values, body []byte) (
	response map[string]interface{}, statusCode int, err error) {
	var res *http.Response
	if params != nil {
		url += params.Encode()
	}
	if method == "GET" {
		res, err = http.Get(url)
	} else {
		res, err = http.Post(url,
			"application/json", bytes.NewBuffer(body))
	}
	if err != nil {
		return
	}
	statusCode = res.StatusCode
	if statusCode != 200 {
		return
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &response)
	fmt.Println(response)
	return
}

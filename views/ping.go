package views

import (
	"calendar/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

var PingHandler = operations.GetPingHandlerFunc(func(params operations.GetPingParams) middleware.Responder {
	return operations.NewGetPingOK()
})

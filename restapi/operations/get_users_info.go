// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetUsersInfoHandlerFunc turns a function with the right signature into a get users info handler
type GetUsersInfoHandlerFunc func(GetUsersInfoParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetUsersInfoHandlerFunc) Handle(params GetUsersInfoParams) middleware.Responder {
	return fn(params)
}

// GetUsersInfoHandler interface for that can handle valid get users info params
type GetUsersInfoHandler interface {
	Handle(GetUsersInfoParams) middleware.Responder
}

// NewGetUsersInfo creates a new http.Handler for the get users info operation
func NewGetUsersInfo(ctx *middleware.Context, handler GetUsersInfoHandler) *GetUsersInfo {
	return &GetUsersInfo{Context: ctx, Handler: handler}
}

/* GetUsersInfo swagger:route GET /users/info getUsersInfo

GetUsersInfo get users info API

*/
type GetUsersInfo struct {
	Context *middleware.Context
	Handler GetUsersInfoHandler
}

func (o *GetUsersInfo) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetUsersInfoParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetPingOKCode is the HTTP code returned for type GetPingOK
const GetPingOKCode int = 200

/*GetPingOK Ok

swagger:response getPingOK
*/
type GetPingOK struct {
}

// NewGetPingOK creates GetPingOK with default headers values
func NewGetPingOK() *GetPingOK {

	return &GetPingOK{}
}

// WriteResponse to the client
func (o *GetPingOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

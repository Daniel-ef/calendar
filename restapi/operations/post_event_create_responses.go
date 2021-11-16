// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"calendar/models"
)

// PostEventCreateOKCode is the HTTP code returned for type PostEventCreateOK
const PostEventCreateOKCode int = 200

/*PostEventCreateOK Ok

swagger:response postEventCreateOK
*/
type PostEventCreateOK struct {

	/*
	  In: Body
	*/
	Payload *models.EventCreateResponse `json:"body,omitempty"`
}

// NewPostEventCreateOK creates PostEventCreateOK with default headers values
func NewPostEventCreateOK() *PostEventCreateOK {

	return &PostEventCreateOK{}
}

// WithPayload adds the payload to the post event create o k response
func (o *PostEventCreateOK) WithPayload(payload *models.EventCreateResponse) *PostEventCreateOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post event create o k response
func (o *PostEventCreateOK) SetPayload(payload *models.EventCreateResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostEventCreateOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostEventCreateBadRequestCode is the HTTP code returned for type PostEventCreateBadRequest
const PostEventCreateBadRequestCode int = 400

/*PostEventCreateBadRequest Creation failed

swagger:response postEventCreateBadRequest
*/
type PostEventCreateBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

// NewPostEventCreateBadRequest creates PostEventCreateBadRequest with default headers values
func NewPostEventCreateBadRequest() *PostEventCreateBadRequest {

	return &PostEventCreateBadRequest{}
}

// WithPayload adds the payload to the post event create bad request response
func (o *PostEventCreateBadRequest) WithPayload(payload *models.ErrorResponse) *PostEventCreateBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post event create bad request response
func (o *PostEventCreateBadRequest) SetPayload(payload *models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostEventCreateBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// UsersFreeSlotResponse users free slot response
//
// swagger:model UsersFreeSlotResponse
type UsersFreeSlotResponse struct {

	// time end
	// Required: true
	// Format: date-time
	TimeEnd *strfmt.DateTime `json:"time_end"`

	// time start
	// Required: true
	// Format: date-time
	TimeStart *strfmt.DateTime `json:"time_start"`
}

// Validate validates this users free slot response
func (m *UsersFreeSlotResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTimeEnd(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTimeStart(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *UsersFreeSlotResponse) validateTimeEnd(formats strfmt.Registry) error {

	if err := validate.Required("time_end", "body", m.TimeEnd); err != nil {
		return err
	}

	if err := validate.FormatOf("time_end", "body", "date-time", m.TimeEnd.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *UsersFreeSlotResponse) validateTimeStart(formats strfmt.Registry) error {

	if err := validate.Required("time_start", "body", m.TimeStart); err != nil {
		return err
	}

	if err := validate.FormatOf("time_start", "body", "date-time", m.TimeStart.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this users free slot response based on context it is used
func (m *UsersFreeSlotResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *UsersFreeSlotResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UsersFreeSlotResponse) UnmarshalBinary(b []byte) error {
	var res UsersFreeSlotResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

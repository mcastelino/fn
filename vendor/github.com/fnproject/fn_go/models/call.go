// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Call call
// swagger:model Call
type Call struct {

	// App name that is assigned to a route that is being executed.
	// Read Only: true
	AppName string `json:"app_name,omitempty"`

	// Time when call completed, whether it was successul or failed. Always in UTC.
	// Read Only: true
	CompletedAt strfmt.DateTime `json:"completed_at,omitempty"`

	// Time when call was submitted. Always in UTC.
	// Read Only: true
	CreatedAt strfmt.DateTime `json:"created_at,omitempty"`

	// Call execution error, if status is 'error'.
	// Read Only: true
	Error string `json:"error,omitempty"`

	// Call UUID ID.
	// Read Only: true
	ID string `json:"id,omitempty"`

	// App route that is being executed.
	// Read Only: true
	Path string `json:"path,omitempty"`

	// Time when call started execution. Always in UTC.
	// Read Only: true
	StartedAt strfmt.DateTime `json:"started_at,omitempty"`

	// A histogram of stats for a call, each is a snapshot of a calls state at the timestamp.
	// Read Only: true
	Stats []*Stat `json:"stats"`

	// Call execution status.
	// Read Only: true
	Status string `json:"status,omitempty"`
}

// Validate validates this call
func (m *Call) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCompletedAt(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateStartedAt(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateStats(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Call) validateCompletedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.CompletedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("completed_at", "body", "date-time", m.CompletedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Call) validateCreatedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("created_at", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Call) validateStartedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.StartedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("started_at", "body", "date-time", m.StartedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Call) validateStats(formats strfmt.Registry) error {

	if swag.IsZero(m.Stats) { // not required
		return nil
	}

	for i := 0; i < len(m.Stats); i++ {

		if swag.IsZero(m.Stats[i]) { // not required
			continue
		}

		if m.Stats[i] != nil {

			if err := m.Stats[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("stats" + "." + strconv.Itoa(i))
				}
				return err
			}

		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Call) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Call) UnmarshalBinary(b []byte) error {
	var res Call
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

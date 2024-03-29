// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	custom "github.com/kihamo/boggart/providers/pass24online/static/model"
)

// Feed feed
//
// swagger:model Feed
type Feed struct {

	// event data
	EventData map[string]interface{} `json:"eventData,omitempty"`

	// happened at
	// Format: date-time
	HappenedAt *custom.DateTime `json:"happenedAt,omitempty"`

	// initiated by
	InitiatedBy *FeedInitiatedBy `json:"initiatedBy,omitempty"`

	// message
	Message string `json:"message,omitempty"`

	// subject
	Subject *Pass `json:"subject,omitempty"`

	// subject type
	SubjectType string `json:"subjectType,omitempty"`

	// title
	Title string `json:"title,omitempty"`

	// type
	Type uint64 `json:"type,omitempty"`
}

// Validate validates this feed
func (m *Feed) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHappenedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInitiatedBy(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSubject(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Feed) validateHappenedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.HappenedAt) { // not required
		return nil
	}

	if m.HappenedAt != nil {
		if err := m.HappenedAt.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("happenedAt")
			}
			return err
		}
	}

	return nil
}

func (m *Feed) validateInitiatedBy(formats strfmt.Registry) error {

	if swag.IsZero(m.InitiatedBy) { // not required
		return nil
	}

	if m.InitiatedBy != nil {
		if err := m.InitiatedBy.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("initiatedBy")
			}
			return err
		}
	}

	return nil
}

func (m *Feed) validateSubject(formats strfmt.Registry) error {

	if swag.IsZero(m.Subject) { // not required
		return nil
	}

	if m.Subject != nil {
		if err := m.Subject.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("subject")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Feed) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Feed) UnmarshalBinary(b []byte) error {
	var res Feed
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// FeedInitiatedBy feed initiated by
//
// swagger:model FeedInitiatedBy
type FeedInitiatedBy struct {

	// id
	ID uint64 `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// type
	Type string `json:"type,omitempty"`
}

// Validate validates this feed initiated by
func (m *FeedInitiatedBy) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *FeedInitiatedBy) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FeedInitiatedBy) UnmarshalBinary(b []byte) error {
	var res FeedInitiatedBy
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

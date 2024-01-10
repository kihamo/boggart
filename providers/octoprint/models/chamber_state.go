// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ChamberState chamber state
//
// swagger:model ChamberState
type ChamberState struct {

	// chamber
	Chamber *TemperatureData `json:"chamber,omitempty"`

	// history
	History []*ChamberStateHistoryItems0 `json:"history"`
}

// Validate validates this chamber state
func (m *ChamberState) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateChamber(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHistory(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ChamberState) validateChamber(formats strfmt.Registry) error {
	if swag.IsZero(m.Chamber) { // not required
		return nil
	}

	if m.Chamber != nil {
		if err := m.Chamber.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("chamber")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("chamber")
			}
			return err
		}
	}

	return nil
}

func (m *ChamberState) validateHistory(formats strfmt.Registry) error {
	if swag.IsZero(m.History) { // not required
		return nil
	}

	for i := 0; i < len(m.History); i++ {
		if swag.IsZero(m.History[i]) { // not required
			continue
		}

		if m.History[i] != nil {
			if err := m.History[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("history" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("history" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this chamber state based on the context it is used
func (m *ChamberState) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateChamber(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateHistory(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ChamberState) contextValidateChamber(ctx context.Context, formats strfmt.Registry) error {

	if m.Chamber != nil {

		if swag.IsZero(m.Chamber) { // not required
			return nil
		}

		if err := m.Chamber.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("chamber")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("chamber")
			}
			return err
		}
	}

	return nil
}

func (m *ChamberState) contextValidateHistory(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.History); i++ {

		if m.History[i] != nil {

			if swag.IsZero(m.History[i]) { // not required
				return nil
			}

			if err := m.History[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("history" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("history" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ChamberState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ChamberState) UnmarshalBinary(b []byte) error {
	var res ChamberState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ChamberStateHistoryItems0 chamber state history items0
//
// swagger:model ChamberStateHistoryItems0
type ChamberStateHistoryItems0 struct {

	// chamber
	Chamber *TemperatureData `json:"chamber,omitempty"`

	// time
	Time int64 `json:"time,omitempty"`
}

// Validate validates this chamber state history items0
func (m *ChamberStateHistoryItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateChamber(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ChamberStateHistoryItems0) validateChamber(formats strfmt.Registry) error {
	if swag.IsZero(m.Chamber) { // not required
		return nil
	}

	if m.Chamber != nil {
		if err := m.Chamber.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("chamber")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("chamber")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this chamber state history items0 based on the context it is used
func (m *ChamberStateHistoryItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateChamber(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ChamberStateHistoryItems0) contextValidateChamber(ctx context.Context, formats strfmt.Registry) error {

	if m.Chamber != nil {

		if swag.IsZero(m.Chamber) { // not required
			return nil
		}

		if err := m.Chamber.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("chamber")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("chamber")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ChamberStateHistoryItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ChamberStateHistoryItems0) UnmarshalBinary(b []byte) error {
	var res ChamberStateHistoryItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

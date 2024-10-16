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
	static "github.com/kihamo/boggart/providers/openweathermap/static/models"
)

// OneCallTimeMachine one call time machine
//
// swagger:model OneCallTimeMachine
type OneCallTimeMachine struct {

	// current
	Current *OneCallTimeMachineCurrent `json:"current,omitempty"`

	// hourly
	Hourly []*OneCallTimeMachineHourlyItems0 `json:"hourly"`

	// lat
	Lat float64 `json:"lat,omitempty"`

	// lon
	Lon float64 `json:"lon,omitempty"`

	// timezone
	Timezone string `json:"timezone,omitempty"`

	// timezone offset
	TimezoneOffset uint64 `json:"timezone_offset,omitempty"`
}

// Validate validates this one call time machine
func (m *OneCallTimeMachine) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCurrent(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHourly(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OneCallTimeMachine) validateCurrent(formats strfmt.Registry) error {
	if swag.IsZero(m.Current) { // not required
		return nil
	}

	if m.Current != nil {
		if err := m.Current.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("current")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("current")
			}
			return err
		}
	}

	return nil
}

func (m *OneCallTimeMachine) validateHourly(formats strfmt.Registry) error {
	if swag.IsZero(m.Hourly) { // not required
		return nil
	}

	for i := 0; i < len(m.Hourly); i++ {
		if swag.IsZero(m.Hourly[i]) { // not required
			continue
		}

		if m.Hourly[i] != nil {
			if err := m.Hourly[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("hourly" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("hourly" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this one call time machine based on the context it is used
func (m *OneCallTimeMachine) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCurrent(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateHourly(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OneCallTimeMachine) contextValidateCurrent(ctx context.Context, formats strfmt.Registry) error {

	if m.Current != nil {

		if swag.IsZero(m.Current) { // not required
			return nil
		}

		if err := m.Current.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("current")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("current")
			}
			return err
		}
	}

	return nil
}

func (m *OneCallTimeMachine) contextValidateHourly(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Hourly); i++ {

		if m.Hourly[i] != nil {

			if swag.IsZero(m.Hourly[i]) { // not required
				return nil
			}

			if err := m.Hourly[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("hourly" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("hourly" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *OneCallTimeMachine) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OneCallTimeMachine) UnmarshalBinary(b []byte) error {
	var res OneCallTimeMachine
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// OneCallTimeMachineCurrent one call time machine current
//
// swagger:model OneCallTimeMachineCurrent
type OneCallTimeMachineCurrent struct {

	// clouds
	Clouds uint64 `json:"clouds,omitempty"`

	// dew point
	DewPoint float64 `json:"dew_point,omitempty"`

	// dt
	// Format: date-time
	Dt static.DateTime `json:"dt,omitempty"`

	// feels like
	FeelsLike float64 `json:"feels_like,omitempty"`

	// humidity
	Humidity uint64 `json:"humidity,omitempty"`

	// pressure
	Pressure float64 `json:"pressure,omitempty"`

	// rain
	Rain float64 `json:"rain,omitempty"`

	// snow
	Snow float64 `json:"snow,omitempty"`

	// sunrise
	// Format: date-time
	Sunrise static.DateTime `json:"sunrise,omitempty"`

	// sunset
	// Format: date-time
	Sunset static.DateTime `json:"sunset,omitempty"`

	// temp
	Temp float64 `json:"temp,omitempty"`

	// uvi
	Uvi float64 `json:"uvi,omitempty"`

	// visibility
	Visibility uint64 `json:"visibility,omitempty"`

	// weather
	Weather []*Weather `json:"weather"`

	// wind deg
	WindDeg uint64 `json:"wind_deg,omitempty"`

	// wind gust
	WindGust float64 `json:"wind_gust,omitempty"`

	// wind speed
	WindSpeed float64 `json:"wind_speed,omitempty"`
}

// Validate validates this one call time machine current
func (m *OneCallTimeMachineCurrent) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSunrise(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSunset(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateWeather(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OneCallTimeMachineCurrent) validateDt(formats strfmt.Registry) error {
	if swag.IsZero(m.Dt) { // not required
		return nil
	}

	if err := m.Dt.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("current" + "." + "dt")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("current" + "." + "dt")
		}
		return err
	}

	return nil
}

func (m *OneCallTimeMachineCurrent) validateSunrise(formats strfmt.Registry) error {
	if swag.IsZero(m.Sunrise) { // not required
		return nil
	}

	if err := m.Sunrise.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("current" + "." + "sunrise")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("current" + "." + "sunrise")
		}
		return err
	}

	return nil
}

func (m *OneCallTimeMachineCurrent) validateSunset(formats strfmt.Registry) error {
	if swag.IsZero(m.Sunset) { // not required
		return nil
	}

	if err := m.Sunset.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("current" + "." + "sunset")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("current" + "." + "sunset")
		}
		return err
	}

	return nil
}

func (m *OneCallTimeMachineCurrent) validateWeather(formats strfmt.Registry) error {
	if swag.IsZero(m.Weather) { // not required
		return nil
	}

	for i := 0; i < len(m.Weather); i++ {
		if swag.IsZero(m.Weather[i]) { // not required
			continue
		}

		if m.Weather[i] != nil {
			if err := m.Weather[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("current" + "." + "weather" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("current" + "." + "weather" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this one call time machine current based on the context it is used
func (m *OneCallTimeMachineCurrent) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSunrise(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSunset(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateWeather(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OneCallTimeMachineCurrent) contextValidateDt(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Dt) { // not required
		return nil
	}

	if err := m.Dt.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("current" + "." + "dt")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("current" + "." + "dt")
		}
		return err
	}

	return nil
}

func (m *OneCallTimeMachineCurrent) contextValidateSunrise(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Sunrise) { // not required
		return nil
	}

	if err := m.Sunrise.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("current" + "." + "sunrise")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("current" + "." + "sunrise")
		}
		return err
	}

	return nil
}

func (m *OneCallTimeMachineCurrent) contextValidateSunset(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Sunset) { // not required
		return nil
	}

	if err := m.Sunset.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("current" + "." + "sunset")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("current" + "." + "sunset")
		}
		return err
	}

	return nil
}

func (m *OneCallTimeMachineCurrent) contextValidateWeather(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Weather); i++ {

		if m.Weather[i] != nil {

			if swag.IsZero(m.Weather[i]) { // not required
				return nil
			}

			if err := m.Weather[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("current" + "." + "weather" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("current" + "." + "weather" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *OneCallTimeMachineCurrent) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OneCallTimeMachineCurrent) UnmarshalBinary(b []byte) error {
	var res OneCallTimeMachineCurrent
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// OneCallTimeMachineHourlyItems0 one call time machine hourly items0
//
// swagger:model OneCallTimeMachineHourlyItems0
type OneCallTimeMachineHourlyItems0 struct {

	// clouds
	Clouds uint64 `json:"clouds,omitempty"`

	// dew point
	DewPoint float64 `json:"dew_point,omitempty"`

	// dt
	// Format: date-time
	Dt static.DateTime `json:"dt,omitempty"`

	// feels like
	FeelsLike float64 `json:"feels_like,omitempty"`

	// humidity
	Humidity uint64 `json:"humidity,omitempty"`

	// pressure
	Pressure float64 `json:"pressure,omitempty"`

	// rain
	Rain *Rain `json:"rain,omitempty"`

	// snow
	Snow *Snow `json:"snow,omitempty"`

	// temp
	Temp float64 `json:"temp,omitempty"`

	// visibility
	Visibility uint64 `json:"visibility,omitempty"`

	// weather
	Weather []*Weather `json:"weather"`

	// wind deg
	WindDeg uint64 `json:"wind_deg,omitempty"`

	// wind gust
	WindGust float64 `json:"wind_gust,omitempty"`

	// wind speed
	WindSpeed float64 `json:"wind_speed,omitempty"`
}

// Validate validates this one call time machine hourly items0
func (m *OneCallTimeMachineHourlyItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRain(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSnow(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateWeather(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OneCallTimeMachineHourlyItems0) validateDt(formats strfmt.Registry) error {
	if swag.IsZero(m.Dt) { // not required
		return nil
	}

	if err := m.Dt.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("dt")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("dt")
		}
		return err
	}

	return nil
}

func (m *OneCallTimeMachineHourlyItems0) validateRain(formats strfmt.Registry) error {
	if swag.IsZero(m.Rain) { // not required
		return nil
	}

	if m.Rain != nil {
		if err := m.Rain.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("rain")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("rain")
			}
			return err
		}
	}

	return nil
}

func (m *OneCallTimeMachineHourlyItems0) validateSnow(formats strfmt.Registry) error {
	if swag.IsZero(m.Snow) { // not required
		return nil
	}

	if m.Snow != nil {
		if err := m.Snow.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("snow")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("snow")
			}
			return err
		}
	}

	return nil
}

func (m *OneCallTimeMachineHourlyItems0) validateWeather(formats strfmt.Registry) error {
	if swag.IsZero(m.Weather) { // not required
		return nil
	}

	for i := 0; i < len(m.Weather); i++ {
		if swag.IsZero(m.Weather[i]) { // not required
			continue
		}

		if m.Weather[i] != nil {
			if err := m.Weather[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("weather" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("weather" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this one call time machine hourly items0 based on the context it is used
func (m *OneCallTimeMachineHourlyItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRain(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSnow(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateWeather(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OneCallTimeMachineHourlyItems0) contextValidateDt(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Dt) { // not required
		return nil
	}

	if err := m.Dt.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("dt")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("dt")
		}
		return err
	}

	return nil
}

func (m *OneCallTimeMachineHourlyItems0) contextValidateRain(ctx context.Context, formats strfmt.Registry) error {

	if m.Rain != nil {

		if swag.IsZero(m.Rain) { // not required
			return nil
		}

		if err := m.Rain.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("rain")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("rain")
			}
			return err
		}
	}

	return nil
}

func (m *OneCallTimeMachineHourlyItems0) contextValidateSnow(ctx context.Context, formats strfmt.Registry) error {

	if m.Snow != nil {

		if swag.IsZero(m.Snow) { // not required
			return nil
		}

		if err := m.Snow.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("snow")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("snow")
			}
			return err
		}
	}

	return nil
}

func (m *OneCallTimeMachineHourlyItems0) contextValidateWeather(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Weather); i++ {

		if m.Weather[i] != nil {

			if swag.IsZero(m.Weather[i]) { // not required
				return nil
			}

			if err := m.Weather[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("weather" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("weather" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *OneCallTimeMachineHourlyItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OneCallTimeMachineHourlyItems0) UnmarshalBinary(b []byte) error {
	var res OneCallTimeMachineHourlyItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

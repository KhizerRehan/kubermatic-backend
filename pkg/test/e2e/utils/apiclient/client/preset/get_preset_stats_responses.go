// Code generated by go-swagger; DO NOT EDIT.

package preset

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// GetPresetStatsReader is a Reader for the GetPresetStats structure.
type GetPresetStatsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPresetStatsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPresetStatsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetPresetStatsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewGetPresetStatsForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetPresetStatsNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetPresetStatsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetPresetStatsOK creates a GetPresetStatsOK with default headers values
func NewGetPresetStatsOK() *GetPresetStatsOK {
	return &GetPresetStatsOK{}
}

/*
GetPresetStatsOK describes a response with status code 200, with default header values.

PresetStats
*/
type GetPresetStatsOK struct {
	Payload *models.PresetStats
}

// IsSuccess returns true when this get preset stats o k response has a 2xx status code
func (o *GetPresetStatsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get preset stats o k response has a 3xx status code
func (o *GetPresetStatsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get preset stats o k response has a 4xx status code
func (o *GetPresetStatsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get preset stats o k response has a 5xx status code
func (o *GetPresetStatsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get preset stats o k response a status code equal to that given
func (o *GetPresetStatsOK) IsCode(code int) bool {
	return code == 200
}

func (o *GetPresetStatsOK) Error() string {
	return fmt.Sprintf("[GET /api/v2/presets/{preset_name}/stats][%d] getPresetStatsOK  %+v", 200, o.Payload)
}

func (o *GetPresetStatsOK) String() string {
	return fmt.Sprintf("[GET /api/v2/presets/{preset_name}/stats][%d] getPresetStatsOK  %+v", 200, o.Payload)
}

func (o *GetPresetStatsOK) GetPayload() *models.PresetStats {
	return o.Payload
}

func (o *GetPresetStatsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PresetStats)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPresetStatsUnauthorized creates a GetPresetStatsUnauthorized with default headers values
func NewGetPresetStatsUnauthorized() *GetPresetStatsUnauthorized {
	return &GetPresetStatsUnauthorized{}
}

/*
GetPresetStatsUnauthorized describes a response with status code 401, with default header values.

EmptyResponse is a empty response
*/
type GetPresetStatsUnauthorized struct {
}

// IsSuccess returns true when this get preset stats unauthorized response has a 2xx status code
func (o *GetPresetStatsUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get preset stats unauthorized response has a 3xx status code
func (o *GetPresetStatsUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get preset stats unauthorized response has a 4xx status code
func (o *GetPresetStatsUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this get preset stats unauthorized response has a 5xx status code
func (o *GetPresetStatsUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this get preset stats unauthorized response a status code equal to that given
func (o *GetPresetStatsUnauthorized) IsCode(code int) bool {
	return code == 401
}

func (o *GetPresetStatsUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/v2/presets/{preset_name}/stats][%d] getPresetStatsUnauthorized ", 401)
}

func (o *GetPresetStatsUnauthorized) String() string {
	return fmt.Sprintf("[GET /api/v2/presets/{preset_name}/stats][%d] getPresetStatsUnauthorized ", 401)
}

func (o *GetPresetStatsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetPresetStatsForbidden creates a GetPresetStatsForbidden with default headers values
func NewGetPresetStatsForbidden() *GetPresetStatsForbidden {
	return &GetPresetStatsForbidden{}
}

/*
GetPresetStatsForbidden describes a response with status code 403, with default header values.

EmptyResponse is a empty response
*/
type GetPresetStatsForbidden struct {
}

// IsSuccess returns true when this get preset stats forbidden response has a 2xx status code
func (o *GetPresetStatsForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get preset stats forbidden response has a 3xx status code
func (o *GetPresetStatsForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get preset stats forbidden response has a 4xx status code
func (o *GetPresetStatsForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this get preset stats forbidden response has a 5xx status code
func (o *GetPresetStatsForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this get preset stats forbidden response a status code equal to that given
func (o *GetPresetStatsForbidden) IsCode(code int) bool {
	return code == 403
}

func (o *GetPresetStatsForbidden) Error() string {
	return fmt.Sprintf("[GET /api/v2/presets/{preset_name}/stats][%d] getPresetStatsForbidden ", 403)
}

func (o *GetPresetStatsForbidden) String() string {
	return fmt.Sprintf("[GET /api/v2/presets/{preset_name}/stats][%d] getPresetStatsForbidden ", 403)
}

func (o *GetPresetStatsForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetPresetStatsNotFound creates a GetPresetStatsNotFound with default headers values
func NewGetPresetStatsNotFound() *GetPresetStatsNotFound {
	return &GetPresetStatsNotFound{}
}

/*
GetPresetStatsNotFound describes a response with status code 404, with default header values.

EmptyResponse is a empty response
*/
type GetPresetStatsNotFound struct {
}

// IsSuccess returns true when this get preset stats not found response has a 2xx status code
func (o *GetPresetStatsNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get preset stats not found response has a 3xx status code
func (o *GetPresetStatsNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get preset stats not found response has a 4xx status code
func (o *GetPresetStatsNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this get preset stats not found response has a 5xx status code
func (o *GetPresetStatsNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this get preset stats not found response a status code equal to that given
func (o *GetPresetStatsNotFound) IsCode(code int) bool {
	return code == 404
}

func (o *GetPresetStatsNotFound) Error() string {
	return fmt.Sprintf("[GET /api/v2/presets/{preset_name}/stats][%d] getPresetStatsNotFound ", 404)
}

func (o *GetPresetStatsNotFound) String() string {
	return fmt.Sprintf("[GET /api/v2/presets/{preset_name}/stats][%d] getPresetStatsNotFound ", 404)
}

func (o *GetPresetStatsNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetPresetStatsDefault creates a GetPresetStatsDefault with default headers values
func NewGetPresetStatsDefault(code int) *GetPresetStatsDefault {
	return &GetPresetStatsDefault{
		_statusCode: code,
	}
}

/*
GetPresetStatsDefault describes a response with status code -1, with default header values.

errorResponse
*/
type GetPresetStatsDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the get preset stats default response
func (o *GetPresetStatsDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this get preset stats default response has a 2xx status code
func (o *GetPresetStatsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get preset stats default response has a 3xx status code
func (o *GetPresetStatsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get preset stats default response has a 4xx status code
func (o *GetPresetStatsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get preset stats default response has a 5xx status code
func (o *GetPresetStatsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get preset stats default response a status code equal to that given
func (o *GetPresetStatsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *GetPresetStatsDefault) Error() string {
	return fmt.Sprintf("[GET /api/v2/presets/{preset_name}/stats][%d] getPresetStats default  %+v", o._statusCode, o.Payload)
}

func (o *GetPresetStatsDefault) String() string {
	return fmt.Sprintf("[GET /api/v2/presets/{preset_name}/stats][%d] getPresetStats default  %+v", o._statusCode, o.Payload)
}

func (o *GetPresetStatsDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetPresetStatsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

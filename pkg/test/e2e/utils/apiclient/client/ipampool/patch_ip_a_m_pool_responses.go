// Code generated by go-swagger; DO NOT EDIT.

package ipampool

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// PatchIPAMPoolReader is a Reader for the PatchIPAMPool structure.
type PatchIPAMPoolReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PatchIPAMPoolReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPatchIPAMPoolOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewPatchIPAMPoolUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewPatchIPAMPoolForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewPatchIPAMPoolDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPatchIPAMPoolOK creates a PatchIPAMPoolOK with default headers values
func NewPatchIPAMPoolOK() *PatchIPAMPoolOK {
	return &PatchIPAMPoolOK{}
}

/*
PatchIPAMPoolOK describes a response with status code 200, with default header values.

EmptyResponse is a empty response
*/
type PatchIPAMPoolOK struct {
}

// IsSuccess returns true when this patch Ip a m pool o k response has a 2xx status code
func (o *PatchIPAMPoolOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this patch Ip a m pool o k response has a 3xx status code
func (o *PatchIPAMPoolOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this patch Ip a m pool o k response has a 4xx status code
func (o *PatchIPAMPoolOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this patch Ip a m pool o k response has a 5xx status code
func (o *PatchIPAMPoolOK) IsServerError() bool {
	return false
}

// IsCode returns true when this patch Ip a m pool o k response a status code equal to that given
func (o *PatchIPAMPoolOK) IsCode(code int) bool {
	return code == 200
}

func (o *PatchIPAMPoolOK) Error() string {
	return fmt.Sprintf("[PATCH /api/v2/seeds/{seed_name}/ipampools/{ipampool_name}][%d] patchIpAMPoolOK ", 200)
}

func (o *PatchIPAMPoolOK) String() string {
	return fmt.Sprintf("[PATCH /api/v2/seeds/{seed_name}/ipampools/{ipampool_name}][%d] patchIpAMPoolOK ", 200)
}

func (o *PatchIPAMPoolOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewPatchIPAMPoolUnauthorized creates a PatchIPAMPoolUnauthorized with default headers values
func NewPatchIPAMPoolUnauthorized() *PatchIPAMPoolUnauthorized {
	return &PatchIPAMPoolUnauthorized{}
}

/*
PatchIPAMPoolUnauthorized describes a response with status code 401, with default header values.

EmptyResponse is a empty response
*/
type PatchIPAMPoolUnauthorized struct {
}

// IsSuccess returns true when this patch Ip a m pool unauthorized response has a 2xx status code
func (o *PatchIPAMPoolUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this patch Ip a m pool unauthorized response has a 3xx status code
func (o *PatchIPAMPoolUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this patch Ip a m pool unauthorized response has a 4xx status code
func (o *PatchIPAMPoolUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this patch Ip a m pool unauthorized response has a 5xx status code
func (o *PatchIPAMPoolUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this patch Ip a m pool unauthorized response a status code equal to that given
func (o *PatchIPAMPoolUnauthorized) IsCode(code int) bool {
	return code == 401
}

func (o *PatchIPAMPoolUnauthorized) Error() string {
	return fmt.Sprintf("[PATCH /api/v2/seeds/{seed_name}/ipampools/{ipampool_name}][%d] patchIpAMPoolUnauthorized ", 401)
}

func (o *PatchIPAMPoolUnauthorized) String() string {
	return fmt.Sprintf("[PATCH /api/v2/seeds/{seed_name}/ipampools/{ipampool_name}][%d] patchIpAMPoolUnauthorized ", 401)
}

func (o *PatchIPAMPoolUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewPatchIPAMPoolForbidden creates a PatchIPAMPoolForbidden with default headers values
func NewPatchIPAMPoolForbidden() *PatchIPAMPoolForbidden {
	return &PatchIPAMPoolForbidden{}
}

/*
PatchIPAMPoolForbidden describes a response with status code 403, with default header values.

EmptyResponse is a empty response
*/
type PatchIPAMPoolForbidden struct {
}

// IsSuccess returns true when this patch Ip a m pool forbidden response has a 2xx status code
func (o *PatchIPAMPoolForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this patch Ip a m pool forbidden response has a 3xx status code
func (o *PatchIPAMPoolForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this patch Ip a m pool forbidden response has a 4xx status code
func (o *PatchIPAMPoolForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this patch Ip a m pool forbidden response has a 5xx status code
func (o *PatchIPAMPoolForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this patch Ip a m pool forbidden response a status code equal to that given
func (o *PatchIPAMPoolForbidden) IsCode(code int) bool {
	return code == 403
}

func (o *PatchIPAMPoolForbidden) Error() string {
	return fmt.Sprintf("[PATCH /api/v2/seeds/{seed_name}/ipampools/{ipampool_name}][%d] patchIpAMPoolForbidden ", 403)
}

func (o *PatchIPAMPoolForbidden) String() string {
	return fmt.Sprintf("[PATCH /api/v2/seeds/{seed_name}/ipampools/{ipampool_name}][%d] patchIpAMPoolForbidden ", 403)
}

func (o *PatchIPAMPoolForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewPatchIPAMPoolDefault creates a PatchIPAMPoolDefault with default headers values
func NewPatchIPAMPoolDefault(code int) *PatchIPAMPoolDefault {
	return &PatchIPAMPoolDefault{
		_statusCode: code,
	}
}

/*
PatchIPAMPoolDefault describes a response with status code -1, with default header values.

errorResponse
*/
type PatchIPAMPoolDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the patch IP a m pool default response
func (o *PatchIPAMPoolDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this patch IP a m pool default response has a 2xx status code
func (o *PatchIPAMPoolDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this patch IP a m pool default response has a 3xx status code
func (o *PatchIPAMPoolDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this patch IP a m pool default response has a 4xx status code
func (o *PatchIPAMPoolDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this patch IP a m pool default response has a 5xx status code
func (o *PatchIPAMPoolDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this patch IP a m pool default response a status code equal to that given
func (o *PatchIPAMPoolDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *PatchIPAMPoolDefault) Error() string {
	return fmt.Sprintf("[PATCH /api/v2/seeds/{seed_name}/ipampools/{ipampool_name}][%d] patchIPAMPool default  %+v", o._statusCode, o.Payload)
}

func (o *PatchIPAMPoolDefault) String() string {
	return fmt.Sprintf("[PATCH /api/v2/seeds/{seed_name}/ipampools/{ipampool_name}][%d] patchIPAMPool default  %+v", o._statusCode, o.Payload)
}

func (o *PatchIPAMPoolDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *PatchIPAMPoolDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

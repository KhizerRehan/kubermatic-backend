// Code generated by go-swagger; DO NOT EDIT.

package aks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// ValidateAKSCredentialsReader is a Reader for the ValidateAKSCredentials structure.
type ValidateAKSCredentialsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ValidateAKSCredentialsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewValidateAKSCredentialsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewValidateAKSCredentialsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewValidateAKSCredentialsOK creates a ValidateAKSCredentialsOK with default headers values
func NewValidateAKSCredentialsOK() *ValidateAKSCredentialsOK {
	return &ValidateAKSCredentialsOK{}
}

/*
ValidateAKSCredentialsOK describes a response with status code 200, with default header values.

EmptyResponse is a empty response
*/
type ValidateAKSCredentialsOK struct {
}

// IsSuccess returns true when this validate a k s credentials o k response has a 2xx status code
func (o *ValidateAKSCredentialsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this validate a k s credentials o k response has a 3xx status code
func (o *ValidateAKSCredentialsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this validate a k s credentials o k response has a 4xx status code
func (o *ValidateAKSCredentialsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this validate a k s credentials o k response has a 5xx status code
func (o *ValidateAKSCredentialsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this validate a k s credentials o k response a status code equal to that given
func (o *ValidateAKSCredentialsOK) IsCode(code int) bool {
	return code == 200
}

func (o *ValidateAKSCredentialsOK) Error() string {
	return fmt.Sprintf("[GET /api/v2/providers/aks/validatecredentials][%d] validateAKSCredentialsOK ", 200)
}

func (o *ValidateAKSCredentialsOK) String() string {
	return fmt.Sprintf("[GET /api/v2/providers/aks/validatecredentials][%d] validateAKSCredentialsOK ", 200)
}

func (o *ValidateAKSCredentialsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewValidateAKSCredentialsDefault creates a ValidateAKSCredentialsDefault with default headers values
func NewValidateAKSCredentialsDefault(code int) *ValidateAKSCredentialsDefault {
	return &ValidateAKSCredentialsDefault{
		_statusCode: code,
	}
}

/*
ValidateAKSCredentialsDefault describes a response with status code -1, with default header values.

errorResponse
*/
type ValidateAKSCredentialsDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the validate a k s credentials default response
func (o *ValidateAKSCredentialsDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this validate a k s credentials default response has a 2xx status code
func (o *ValidateAKSCredentialsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this validate a k s credentials default response has a 3xx status code
func (o *ValidateAKSCredentialsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this validate a k s credentials default response has a 4xx status code
func (o *ValidateAKSCredentialsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this validate a k s credentials default response has a 5xx status code
func (o *ValidateAKSCredentialsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this validate a k s credentials default response a status code equal to that given
func (o *ValidateAKSCredentialsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *ValidateAKSCredentialsDefault) Error() string {
	return fmt.Sprintf("[GET /api/v2/providers/aks/validatecredentials][%d] validateAKSCredentials default  %+v", o._statusCode, o.Payload)
}

func (o *ValidateAKSCredentialsDefault) String() string {
	return fmt.Sprintf("[GET /api/v2/providers/aks/validatecredentials][%d] validateAKSCredentials default  %+v", o._statusCode, o.Payload)
}

func (o *ValidateAKSCredentialsDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ValidateAKSCredentialsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

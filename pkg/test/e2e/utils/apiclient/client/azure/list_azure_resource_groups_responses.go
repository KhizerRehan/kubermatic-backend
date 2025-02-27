// Code generated by go-swagger; DO NOT EDIT.

package azure

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// ListAzureResourceGroupsReader is a Reader for the ListAzureResourceGroups structure.
type ListAzureResourceGroupsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListAzureResourceGroupsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListAzureResourceGroupsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListAzureResourceGroupsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListAzureResourceGroupsOK creates a ListAzureResourceGroupsOK with default headers values
func NewListAzureResourceGroupsOK() *ListAzureResourceGroupsOK {
	return &ListAzureResourceGroupsOK{}
}

/*
ListAzureResourceGroupsOK describes a response with status code 200, with default header values.

AzureResourceGroupsList
*/
type ListAzureResourceGroupsOK struct {
	Payload *models.AzureResourceGroupsList
}

// IsSuccess returns true when this list azure resource groups o k response has a 2xx status code
func (o *ListAzureResourceGroupsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this list azure resource groups o k response has a 3xx status code
func (o *ListAzureResourceGroupsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list azure resource groups o k response has a 4xx status code
func (o *ListAzureResourceGroupsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this list azure resource groups o k response has a 5xx status code
func (o *ListAzureResourceGroupsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this list azure resource groups o k response a status code equal to that given
func (o *ListAzureResourceGroupsOK) IsCode(code int) bool {
	return code == 200
}

func (o *ListAzureResourceGroupsOK) Error() string {
	return fmt.Sprintf("[GET /api/v2/providers/azure/resourcegroups][%d] listAzureResourceGroupsOK  %+v", 200, o.Payload)
}

func (o *ListAzureResourceGroupsOK) String() string {
	return fmt.Sprintf("[GET /api/v2/providers/azure/resourcegroups][%d] listAzureResourceGroupsOK  %+v", 200, o.Payload)
}

func (o *ListAzureResourceGroupsOK) GetPayload() *models.AzureResourceGroupsList {
	return o.Payload
}

func (o *ListAzureResourceGroupsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.AzureResourceGroupsList)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListAzureResourceGroupsDefault creates a ListAzureResourceGroupsDefault with default headers values
func NewListAzureResourceGroupsDefault(code int) *ListAzureResourceGroupsDefault {
	return &ListAzureResourceGroupsDefault{
		_statusCode: code,
	}
}

/*
ListAzureResourceGroupsDefault describes a response with status code -1, with default header values.

errorResponse
*/
type ListAzureResourceGroupsDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the list azure resource groups default response
func (o *ListAzureResourceGroupsDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this list azure resource groups default response has a 2xx status code
func (o *ListAzureResourceGroupsDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this list azure resource groups default response has a 3xx status code
func (o *ListAzureResourceGroupsDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this list azure resource groups default response has a 4xx status code
func (o *ListAzureResourceGroupsDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this list azure resource groups default response has a 5xx status code
func (o *ListAzureResourceGroupsDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this list azure resource groups default response a status code equal to that given
func (o *ListAzureResourceGroupsDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *ListAzureResourceGroupsDefault) Error() string {
	return fmt.Sprintf("[GET /api/v2/providers/azure/resourcegroups][%d] listAzureResourceGroups default  %+v", o._statusCode, o.Payload)
}

func (o *ListAzureResourceGroupsDefault) String() string {
	return fmt.Sprintf("[GET /api/v2/providers/azure/resourcegroups][%d] listAzureResourceGroups default  %+v", o._statusCode, o.Payload)
}

func (o *ListAzureResourceGroupsDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ListAzureResourceGroupsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

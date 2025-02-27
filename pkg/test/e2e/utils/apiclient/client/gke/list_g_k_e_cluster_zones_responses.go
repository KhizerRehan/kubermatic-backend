// Code generated by go-swagger; DO NOT EDIT.

package gke

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// ListGKEClusterZonesReader is a Reader for the ListGKEClusterZones structure.
type ListGKEClusterZonesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListGKEClusterZonesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListGKEClusterZonesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewListGKEClusterZonesUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewListGKEClusterZonesForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewListGKEClusterZonesDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListGKEClusterZonesOK creates a ListGKEClusterZonesOK with default headers values
func NewListGKEClusterZonesOK() *ListGKEClusterZonesOK {
	return &ListGKEClusterZonesOK{}
}

/*
ListGKEClusterZonesOK describes a response with status code 200, with default header values.

GKEZoneList
*/
type ListGKEClusterZonesOK struct {
	Payload models.GKEZoneList
}

// IsSuccess returns true when this list g k e cluster zones o k response has a 2xx status code
func (o *ListGKEClusterZonesOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this list g k e cluster zones o k response has a 3xx status code
func (o *ListGKEClusterZonesOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list g k e cluster zones o k response has a 4xx status code
func (o *ListGKEClusterZonesOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this list g k e cluster zones o k response has a 5xx status code
func (o *ListGKEClusterZonesOK) IsServerError() bool {
	return false
}

// IsCode returns true when this list g k e cluster zones o k response a status code equal to that given
func (o *ListGKEClusterZonesOK) IsCode(code int) bool {
	return code == 200
}

func (o *ListGKEClusterZonesOK) Error() string {
	return fmt.Sprintf("[GET /api/v2/projects/{project_id}/kubernetes/clusters/{cluster_id}/providers/gke/zones][%d] listGKEClusterZonesOK  %+v", 200, o.Payload)
}

func (o *ListGKEClusterZonesOK) String() string {
	return fmt.Sprintf("[GET /api/v2/projects/{project_id}/kubernetes/clusters/{cluster_id}/providers/gke/zones][%d] listGKEClusterZonesOK  %+v", 200, o.Payload)
}

func (o *ListGKEClusterZonesOK) GetPayload() models.GKEZoneList {
	return o.Payload
}

func (o *ListGKEClusterZonesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListGKEClusterZonesUnauthorized creates a ListGKEClusterZonesUnauthorized with default headers values
func NewListGKEClusterZonesUnauthorized() *ListGKEClusterZonesUnauthorized {
	return &ListGKEClusterZonesUnauthorized{}
}

/*
ListGKEClusterZonesUnauthorized describes a response with status code 401, with default header values.

EmptyResponse is a empty response
*/
type ListGKEClusterZonesUnauthorized struct {
}

// IsSuccess returns true when this list g k e cluster zones unauthorized response has a 2xx status code
func (o *ListGKEClusterZonesUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this list g k e cluster zones unauthorized response has a 3xx status code
func (o *ListGKEClusterZonesUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list g k e cluster zones unauthorized response has a 4xx status code
func (o *ListGKEClusterZonesUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this list g k e cluster zones unauthorized response has a 5xx status code
func (o *ListGKEClusterZonesUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this list g k e cluster zones unauthorized response a status code equal to that given
func (o *ListGKEClusterZonesUnauthorized) IsCode(code int) bool {
	return code == 401
}

func (o *ListGKEClusterZonesUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/v2/projects/{project_id}/kubernetes/clusters/{cluster_id}/providers/gke/zones][%d] listGKEClusterZonesUnauthorized ", 401)
}

func (o *ListGKEClusterZonesUnauthorized) String() string {
	return fmt.Sprintf("[GET /api/v2/projects/{project_id}/kubernetes/clusters/{cluster_id}/providers/gke/zones][%d] listGKEClusterZonesUnauthorized ", 401)
}

func (o *ListGKEClusterZonesUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewListGKEClusterZonesForbidden creates a ListGKEClusterZonesForbidden with default headers values
func NewListGKEClusterZonesForbidden() *ListGKEClusterZonesForbidden {
	return &ListGKEClusterZonesForbidden{}
}

/*
ListGKEClusterZonesForbidden describes a response with status code 403, with default header values.

EmptyResponse is a empty response
*/
type ListGKEClusterZonesForbidden struct {
}

// IsSuccess returns true when this list g k e cluster zones forbidden response has a 2xx status code
func (o *ListGKEClusterZonesForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this list g k e cluster zones forbidden response has a 3xx status code
func (o *ListGKEClusterZonesForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list g k e cluster zones forbidden response has a 4xx status code
func (o *ListGKEClusterZonesForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this list g k e cluster zones forbidden response has a 5xx status code
func (o *ListGKEClusterZonesForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this list g k e cluster zones forbidden response a status code equal to that given
func (o *ListGKEClusterZonesForbidden) IsCode(code int) bool {
	return code == 403
}

func (o *ListGKEClusterZonesForbidden) Error() string {
	return fmt.Sprintf("[GET /api/v2/projects/{project_id}/kubernetes/clusters/{cluster_id}/providers/gke/zones][%d] listGKEClusterZonesForbidden ", 403)
}

func (o *ListGKEClusterZonesForbidden) String() string {
	return fmt.Sprintf("[GET /api/v2/projects/{project_id}/kubernetes/clusters/{cluster_id}/providers/gke/zones][%d] listGKEClusterZonesForbidden ", 403)
}

func (o *ListGKEClusterZonesForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewListGKEClusterZonesDefault creates a ListGKEClusterZonesDefault with default headers values
func NewListGKEClusterZonesDefault(code int) *ListGKEClusterZonesDefault {
	return &ListGKEClusterZonesDefault{
		_statusCode: code,
	}
}

/*
ListGKEClusterZonesDefault describes a response with status code -1, with default header values.

errorResponse
*/
type ListGKEClusterZonesDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the list g k e cluster zones default response
func (o *ListGKEClusterZonesDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this list g k e cluster zones default response has a 2xx status code
func (o *ListGKEClusterZonesDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this list g k e cluster zones default response has a 3xx status code
func (o *ListGKEClusterZonesDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this list g k e cluster zones default response has a 4xx status code
func (o *ListGKEClusterZonesDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this list g k e cluster zones default response has a 5xx status code
func (o *ListGKEClusterZonesDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this list g k e cluster zones default response a status code equal to that given
func (o *ListGKEClusterZonesDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *ListGKEClusterZonesDefault) Error() string {
	return fmt.Sprintf("[GET /api/v2/projects/{project_id}/kubernetes/clusters/{cluster_id}/providers/gke/zones][%d] listGKEClusterZones default  %+v", o._statusCode, o.Payload)
}

func (o *ListGKEClusterZonesDefault) String() string {
	return fmt.Sprintf("[GET /api/v2/projects/{project_id}/kubernetes/clusters/{cluster_id}/providers/gke/zones][%d] listGKEClusterZones default  %+v", o._statusCode, o.Payload)
}

func (o *ListGKEClusterZonesDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ListGKEClusterZonesDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

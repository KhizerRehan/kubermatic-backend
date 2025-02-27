// Code generated by go-swagger; DO NOT EDIT.

package project

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// NewCreateGroupProjectBindingParams creates a new CreateGroupProjectBindingParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCreateGroupProjectBindingParams() *CreateGroupProjectBindingParams {
	return &CreateGroupProjectBindingParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCreateGroupProjectBindingParamsWithTimeout creates a new CreateGroupProjectBindingParams object
// with the ability to set a timeout on a request.
func NewCreateGroupProjectBindingParamsWithTimeout(timeout time.Duration) *CreateGroupProjectBindingParams {
	return &CreateGroupProjectBindingParams{
		timeout: timeout,
	}
}

// NewCreateGroupProjectBindingParamsWithContext creates a new CreateGroupProjectBindingParams object
// with the ability to set a context for a request.
func NewCreateGroupProjectBindingParamsWithContext(ctx context.Context) *CreateGroupProjectBindingParams {
	return &CreateGroupProjectBindingParams{
		Context: ctx,
	}
}

// NewCreateGroupProjectBindingParamsWithHTTPClient creates a new CreateGroupProjectBindingParams object
// with the ability to set a custom HTTPClient for a request.
func NewCreateGroupProjectBindingParamsWithHTTPClient(client *http.Client) *CreateGroupProjectBindingParams {
	return &CreateGroupProjectBindingParams{
		HTTPClient: client,
	}
}

/*
CreateGroupProjectBindingParams contains all the parameters to send to the API endpoint

	for the create group project binding operation.

	Typically these are written to a http.Request.
*/
type CreateGroupProjectBindingParams struct {

	// Body.
	Body *models.GroupProjectBindingBody

	// ProjectID.
	ProjectID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the create group project binding params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateGroupProjectBindingParams) WithDefaults() *CreateGroupProjectBindingParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the create group project binding params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateGroupProjectBindingParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the create group project binding params
func (o *CreateGroupProjectBindingParams) WithTimeout(timeout time.Duration) *CreateGroupProjectBindingParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create group project binding params
func (o *CreateGroupProjectBindingParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create group project binding params
func (o *CreateGroupProjectBindingParams) WithContext(ctx context.Context) *CreateGroupProjectBindingParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create group project binding params
func (o *CreateGroupProjectBindingParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create group project binding params
func (o *CreateGroupProjectBindingParams) WithHTTPClient(client *http.Client) *CreateGroupProjectBindingParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create group project binding params
func (o *CreateGroupProjectBindingParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create group project binding params
func (o *CreateGroupProjectBindingParams) WithBody(body *models.GroupProjectBindingBody) *CreateGroupProjectBindingParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create group project binding params
func (o *CreateGroupProjectBindingParams) SetBody(body *models.GroupProjectBindingBody) {
	o.Body = body
}

// WithProjectID adds the projectID to the create group project binding params
func (o *CreateGroupProjectBindingParams) WithProjectID(projectID string) *CreateGroupProjectBindingParams {
	o.SetProjectID(projectID)
	return o
}

// SetProjectID adds the projectId to the create group project binding params
func (o *CreateGroupProjectBindingParams) SetProjectID(projectID string) {
	o.ProjectID = projectID
}

// WriteToRequest writes these params to a swagger request
func (o *CreateGroupProjectBindingParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// path param project_id
	if err := r.SetPathParam("project_id", o.ProjectID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

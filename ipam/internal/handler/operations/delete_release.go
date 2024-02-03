// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeleteReleaseHandlerFunc turns a function with the right signature into a delete release handler
type DeleteReleaseHandlerFunc func(DeleteReleaseParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteReleaseHandlerFunc) Handle(params DeleteReleaseParams) middleware.Responder {
	return fn(params)
}

// DeleteReleaseHandler interface for that can handle valid delete release params
type DeleteReleaseHandler interface {
	Handle(DeleteReleaseParams) middleware.Responder
}

// NewDeleteRelease creates a new http.Handler for the delete release operation
func NewDeleteRelease(ctx *middleware.Context, handler DeleteReleaseHandler) *DeleteRelease {
	return &DeleteRelease{Context: ctx, Handler: handler}
}

/*
	DeleteRelease swagger:route DELETE /release deleteRelease

Release IPv6 range
*/
type DeleteRelease struct {
	Context *middleware.Context
	Handler DeleteReleaseHandler
}

func (o *DeleteRelease) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDeleteReleaseParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

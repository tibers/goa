/*
Package goa provides the runtime support for goa microservices.

Code Generation

goa service development begins with writing the *design* of a service. The design is described using
the goa language implemented by the github.com/goadesign/goa/design/apidsl package. The `goagen` tool
consumes the metadata produced from executing the design language to generate service specific code
that glues the underlying HTTP server with action specific code and data structures.

The goa package contains supporting functionality for the generated code including basic request
and response state management through the RequestData and ResponseData structs, error handling via the
service and controller ErrorHandler field, middleware support via the Middleware data structure as
well as decoding and encoding algorithms.

Request Context

The RequestData and ResponseData structs provides access to the request and response state. goa request
handlers also accept a golang.org/x/net/Context interface as first parameter so that deadlines and
cancelation signals may easily be implemented.

The request state exposes the underlying http.Request object as well as the deserialized payload (request
body) and parameters (both path and querystring parameters). Generated action specific contexts wrap
the context.Context, ResponseData and RequestData data structures. They expose properly typed fields
that correspond to the request parameters and body data structure descriptions appearing in the design.

The response state exposes  the response status and body length as well as the underlying ResponseWriter.
Action contexts provide action specific helper methods that write the responses as described in the
design optionally taking an instance of the media type for responses that contain a body.

Here is an example showing an "update" action corresponding to following design (extract):

	Resource("bottle", func() {
		DefaultMedia(Bottle)
		Action("update", func() {
			Params(func() {
				Param("bottleID", Integer)
			})
			Payload(UpdateBottlePayload)
			Response(OK)
			Response(NotFound)
		})
	})

The action signature generated by goagen is:

	type BottleController interface {
		goa.Controller
		Update(*UpdateBottleContext) error
	}

where UpdateBottleContext is:

	type UpdateBottleContext struct {
        	context.Context   // Timeout and deadline support
		*goa.ResponseData // Response state access
		*goa.RequestData  // Request state access
        	BottleID  int     // Properly typed parameter fields
        	Payload   *UpdateBottlePayload // Properly typed payload
	}

and implements:

	func (ctx *UpdateBottleContext) OK(resp *Bottle) error
	func (ctx *UpdateBottleContext) NotFound() error

The definitions of the Bottle and UpdateBottlePayload data structures are ommitted for brievity.

Controllers

There is one controller interface generated per resource defined via the design language. The
interface exposes the controller actions. User code must provide data structures that implement these
interfaces when mounting a controller onto a service. The controller data structure should include
an anonymous field of type *goa.Controller which takes care of implementing the middleware and
error handler handling.

Error Handling

The controller action methods generated by goagen such as the Update method of the BottleController
interface shown above all return an error value. The controller or service-wide error handler (if no
controller specific error handler) function is invoked whenever the value returned by a controller
action is not nil. The handler gets both the request context and the error as argument.

The default handler implementation returns a response with status code 500 containing the error
message in the body. A different error handler can be specificied using the SetErrorHandler
function on either a controller or service wide. goa comes with an alternative error handler - the
TerseErrorHandler - which also returns a response with status 500 but does not write the error
message to the body of the response.

Middleware

A goa middleware is a function that takes and returns a Handler. A Handler is a the low level
function which handles incoming HTTP requests. goagen generates the handlers code so each handler
creates the action specific context and calls the controller action with it.

Middleware can be added to a goa service or a specific controller using the Service type Use method.
goa comes with a few stock middleware that handle common needs such as logging, panic recovery or
using the RequestID header to trace requests across multiple services.

Validation

The goa design language documented in the dsl package makes it possible to attach validations to
data structure definitions. One specific type of validation consists of defining the format that a
data structure string field must follow. Example of formats include email, data time, hostnames etc.
The ValidateFormat function provides the implementation for the format validation invoked from the
code generated by goagen.

Encoding

The goa design language makes it possible to specify the encodings supported by the API both as
input (Consumes) and output (Produces). goagen uses that information to registed the corresponding
packages with the service encoders and decoders via the SetEncoder and SetDecoder methods. The
service exposes the Decode, DecodeRequest, Encode and EncodeResponse that implement a simple content
type negotiation algorithm for picking the right encoder for the "Accept" request header.

Versioning

The VersionMux interface implemented by the RootMux struct exposes methods used by the generated
code to setup the routing to versioned endpoints. The DSL defines how the API handles versioning:
via request path, header, querystring or a combination. The generated code uses the VersionMux
interface to setup the root mux accordingly.
*/
package goa

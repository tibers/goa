// Code generated by goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service HTTP server types
//
// Command:
// $ goa gen goa.design/goa/examples/security/design -o
// $(GOPATH)/src/goa.design/goa/examples/security

package server

import (
	securedservice "goa.design/goa/examples/security/gen/secured_service"
)

// SigninUnauthorizedResponseBody is the type of the "secured_service" service
// "signin" endpoint HTTP response body for the "unauthorized" error.
type SigninUnauthorizedResponseBody string

// SecureUnauthorizedResponseBody is the type of the "secured_service" service
// "secure" endpoint HTTP response body for the "unauthorized" error.
type SecureUnauthorizedResponseBody string

// DoublySecureUnauthorizedResponseBody is the type of the "secured_service"
// service "doubly_secure" endpoint HTTP response body for the "unauthorized"
// error.
type DoublySecureUnauthorizedResponseBody string

// AlsoDoublySecureUnauthorizedResponseBody is the type of the
// "secured_service" service "also_doubly_secure" endpoint HTTP response body
// for the "unauthorized" error.
type AlsoDoublySecureUnauthorizedResponseBody string

// NewSigninUnauthorizedResponseBody builds the HTTP response body from the
// result of the "signin" endpoint of the "secured_service" service.
func NewSigninUnauthorizedResponseBody(res securedservice.Unauthorized) SigninUnauthorizedResponseBody {
	body := SigninUnauthorizedResponseBody(res)
	return body
}

// NewSecureUnauthorizedResponseBody builds the HTTP response body from the
// result of the "secure" endpoint of the "secured_service" service.
func NewSecureUnauthorizedResponseBody(res securedservice.Unauthorized) SecureUnauthorizedResponseBody {
	body := SecureUnauthorizedResponseBody(res)
	return body
}

// NewDoublySecureUnauthorizedResponseBody builds the HTTP response body from
// the result of the "doubly_secure" endpoint of the "secured_service" service.
func NewDoublySecureUnauthorizedResponseBody(res securedservice.Unauthorized) DoublySecureUnauthorizedResponseBody {
	body := DoublySecureUnauthorizedResponseBody(res)
	return body
}

// NewAlsoDoublySecureUnauthorizedResponseBody builds the HTTP response body
// from the result of the "also_doubly_secure" endpoint of the
// "secured_service" service.
func NewAlsoDoublySecureUnauthorizedResponseBody(res securedservice.Unauthorized) AlsoDoublySecureUnauthorizedResponseBody {
	body := AlsoDoublySecureUnauthorizedResponseBody(res)
	return body
}

// NewSigninPayload builds a secured_service service signin endpoint payload.
func NewSigninPayload() *securedservice.SigninPayload {
	return &securedservice.SigninPayload{}
}

// NewSecurePayload builds a secured_service service secure endpoint payload.
func NewSecurePayload(fail *bool, token *string) *securedservice.SecurePayload {
	return &securedservice.SecurePayload{
		Fail:  fail,
		Token: token,
	}
}

// NewDoublySecurePayload builds a secured_service service doubly_secure
// endpoint payload.
func NewDoublySecurePayload(key *string, token *string) *securedservice.DoublySecurePayload {
	return &securedservice.DoublySecurePayload{
		Key:   key,
		Token: token,
	}
}

// NewAlsoDoublySecurePayload builds a secured_service service
// also_doubly_secure endpoint payload.
func NewAlsoDoublySecurePayload(key *string, oauthToken *string, token *string) *securedservice.AlsoDoublySecurePayload {
	return &securedservice.AlsoDoublySecurePayload{
		Key:        key,
		OauthToken: oauthToken,
		Token:      token,
	}
}

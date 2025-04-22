package id1_client

import "errors"

var ErrNotFound = errors.New("not found")
var ErrNotAuthenticated = errors.New("not authenticated")
var ErrNotAuthorized = errors.New("not authorized")
var ErrInvalidInput = errors.New("invalid input")
var ErrUnexpected = errors.New("unexpected")
var ErrTimeout = errors.New("timeout")

func httpStatusErr(statusCode int) error {
	switch statusCode {
	case 404:
		return ErrNotFound
	case 400:
		return ErrInvalidInput
	case 401:
		return ErrNotAuthenticated
	case 403:
		return ErrNotAuthorized
	case 500:
		return ErrUnexpected
	default:
		return nil
	}
}

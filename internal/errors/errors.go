package errors

import "errors"

var (
	ErrUnmarsh   = errors.New("problem obtain information from the request")
	ErrMarsh     = errors.New("problem creating the response")
	ErrAPIClient = errors.New("Problem creating Dynamo Client")
)

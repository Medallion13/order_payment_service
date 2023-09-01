package errors

import "errors"

var (
	ErrUnmarsh            = errors.New("problem obtain information from the request")
	ErrMarsh              = errors.New("problem creating the response")
	ErrAPIClient          = errors.New("problem creating dynamo client")
	ErrBuildingExpression = errors.New("occurs an error building the expresion for update")
)

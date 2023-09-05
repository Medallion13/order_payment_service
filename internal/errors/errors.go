package errors

import "errors"

var (
	ErrUnmarsh            = errors.New("problem obtain information from the request")
	ErrMarsh              = errors.New("problem creating the response")
	ErrAPIClient          = errors.New("problem creating dynamo client")
	ErrBuildingExpression = errors.New("occurs an error building the expresion for update")
	ErrUpdateDynamo       = errors.New("an error occurred updating the values")
	ErrEventBridgeClient  = errors.New("problems creating event bridge client")
	ErrEventBrigePutEvent = errors.New("an error ocurred putting the event")
)

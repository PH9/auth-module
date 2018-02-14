package main

//
type LogError struct {
	HTTPStatus int
	Step       LogStep
	Code       ErrorCode
	Type       LogErrorType
	Error      error
}

// LogStep
type LogStep int

const (
	ReceiveFromClient LogStep = 1 + iota
	CallThirdParty
	ResponseFromThirdParty
	SuccessResponseToClient
	ErrorResponseToClient
)

//
type ErrorCode int

const (
	StatusUnauthorized ErrorCode = 401
	TokenIsUnder32Char ErrorCode = 401001
	TokenNotFound      ErrorCode = 401002
)

//
type LogErrorType string

const (
	BussinessError LogErrorType = "B"
	TechnicalError LogErrorType = "T"
)

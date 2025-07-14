package utils

import (
	"fmt"
)

const (
	Success     string = "OK"
	SuccessDesc string = "SUCCESS"

	IncompleteRequest string = "ERR_INCOMPLETE"
	IncompleteDesc    string = "ERR_INCOMPLETE"
)

type BadRequestError struct {
	Code    int
	Message string
}

type UnauthorizedError struct {
	Code    int
	Message string
}

type ForbiddenError struct {
	Code    int
	Message string
}

type UnprocessableContentError struct {
	Code    int
	Message string
}

type NotFoundError struct {
	Code    int
	Message string
}

type InternalServerError struct {
	Code    int
	Message string
}

type UnprocessableError struct {
	Code    int
	Message string
}

type ConflictError struct {
	Code    int
	Message string
}

func (err *BadRequestError) Error() string {
	return fmt.Sprintf("BadRequest %v", err.Message)
}

func (err *NotFoundError) Error() string {
	return fmt.Sprintf("NotFound %v", err.Message)
}

func (err *UnprocessableContentError) Error() string {
	return fmt.Sprintf("UnprocessableContent %v", err.Message)
}

func (err *InternalServerError) Error() string {
	return fmt.Sprintf("InternalServerError %v", err.Message)
}

func (err *ConflictError) Error() string {
	return fmt.Sprintf("ConflictError %v", err.Message)
}

func (err *UnauthorizedError) Error() string {
	return fmt.Sprintf("UnauthorizedError %v", err.Message)
}

func (err *ForbiddenError) Error() string {
	return fmt.Sprintf("Forbidden %v", err.Message)
}

package services

import "fmt"

type ErrorAuthentication struct{ Err error }
type ErrorAuthorization struct{ Err error }
type ErrorInternal struct{ Err error }
type ErrorEmpty struct{ What []string }
type ErrorNotFound struct{ What []string }
type ErrorDataAccess struct{ Err error }
type ErrorNoAccess struct{ What []string }
type ErrorUnknown struct{ What []string }
type ErrorNoActive struct {
	What string
	For  string
}
type ErrorAlreadyExists struct{ What []string }

func (e ErrorAuthentication) Error() string {
	return fmt.Sprintf("Authentication error: %v", e.Err)
}

func (e ErrorAuthentication) Unwrap() error {
	return e.Err
}

func (e ErrorAuthorization) Error() string {
	return fmt.Sprintf("Authorization error: %v", e.Err)
}

func (e ErrorAuthorization) Unwrap() error {
	return e.Err
}

func (e ErrorInternal) Error() string {
	return fmt.Sprintf("Internal error occured: %v", e.Err)
}

func (e ErrorInternal) Unwrap() error {
	return e.Err
}

func (e ErrorEmpty) Error() string {
	return fmt.Sprintf("Following information can't be empty: %v", e.What)
}

func (e ErrorNotFound) Error() string {
	return fmt.Sprintf("Not found: %v", e.What)
}

func (e ErrorDataAccess) Error() string {
	return fmt.Sprintf("Error during data access: '%v'", e.Err)
}

func (e ErrorDataAccess) Unwrap() error {
	return e.Err
}

func (e ErrorNoAccess) Error() string {
	return fmt.Sprintf("No access: %v", e.What)
}

func (e ErrorUnknown) Error() string {
	return fmt.Sprintf("Unknown: %v", e.What)
}

func (e ErrorNoActive) Error() string {
	return fmt.Sprintf("No active '%v' for '%v'", e.What, e.For)
}

func (e ErrorAlreadyExists) Error() string {
	return fmt.Sprintf("'%v' already exists", e.What)
}


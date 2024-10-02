package cmnerrors

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
type ErrorConflict struct{ What string }

// Creators
func Authentication(err error) ErrorAuthentication {
	return ErrorAuthentication{err}
}

func Authorization(err error) ErrorAuthorization {
	return ErrorAuthorization{err}
}

func Internal(err error) ErrorInternal {
	return ErrorInternal{err}
}

func Empty(what []string) ErrorEmpty {
	return ErrorEmpty{what}
}

func NotFound(what []string) ErrorNotFound {
	return ErrorNotFound{what}
}

func DataAccess(err error) ErrorDataAccess {
	return ErrorDataAccess{err}
}

func NoAccess(what []string) ErrorNoAccess {
	return ErrorNoAccess{what}
}

func Unknown(what []string) ErrorUnknown {
	return ErrorUnknown{what}
}

func NoActive(what string, for_target string) ErrorNoActive {
	return ErrorNoActive{what, for_target}
}

func AlreadyExists(what []string) ErrorAlreadyExists {
	return ErrorAlreadyExists{what}
}

func Conflict(what string) ErrorConflict {
	return ErrorConflict{what}
}

// Error implementation
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

func (e ErrorConflict) Error() string {
	return fmt.Sprintf("Request reached conflicting state: %v", e.What)
}


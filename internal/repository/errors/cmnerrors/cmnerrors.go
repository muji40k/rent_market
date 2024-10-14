package cmnerrors

import "fmt"

// Any error other error being returned by repository should be treated as
// internal error
type ErrorDuplicate struct{ What []string }
type ErrorNotFound struct{ What []string }

// Creators
func Duplicate(what ...string) ErrorDuplicate {
	return ErrorDuplicate{what}
}

func NotFound(what ...string) ErrorNotFound {
	return ErrorNotFound{what}
}

// Error implementation
func (e ErrorNotFound) Error() string {
	return fmt.Sprintf("Unable to find: %v", e.What)
}

func (e ErrorDuplicate) Error() string {
	return fmt.Sprintf("Object with %v already exists", e.What)
}


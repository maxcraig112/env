package env

import "fmt"

// NotSetError is returned (or panicked with) when a required environment
// variable has not been set.
type NotSetError struct {
	Key string
}

func (e *NotSetError) Error() string {
	return fmt.Sprintf("env: required environment variable %q is not set", e.Key)
}

// ParseError is returned (or panicked with) when an environment variable is
// set but its value cannot be parsed as the requested type.
type ParseError struct {
	Key   string
	Value string
	Type  string
	Err   error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("env: environment variable %q has value %q which cannot be parsed as %s: %v", e.Key, e.Value, e.Type, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

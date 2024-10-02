package guard

import (
	"errors"
	"fmt"
)

var ErrGuard = errors.New("guard error")

// NotNil checks if the provided value is nil and returns a guard error if it is.
func NotNil(value interface{}, param string) error {
	if value == nil {
		return fmt.Errorf("%w: %s cannot be nil", ErrGuard, param)
	}
	return nil
}

// NotEmpty checks if the provided string is empty and returns a guard error if it is.
func NotEmpty(value, param string) error {
	if value == "" {
		return fmt.Errorf("%w: %s cannot be empty", ErrGuard, param)
	}
	return nil
}

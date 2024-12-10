package common

import (
	"errors"
)

func MergeErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	var combinedErr string
	for _, err := range errs {
		combinedErr += err.Error() + "\n"
	}
	return errors.New(combinedErr)
}

func MergeErrorsMessages(errs []string) error {
	if len(errs) == 0 {
		return nil
	}
	var combinedErr string
	for _, err := range errs {
		combinedErr += err + "\n"
	}
	return errors.New(combinedErr)
}

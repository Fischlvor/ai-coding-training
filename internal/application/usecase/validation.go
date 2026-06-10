package usecase

import (
	"fmt"
	"strings"
)

func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	return nil
}

func validateID(id string, field string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: %s is required", ErrValidation, field)
	}
	return nil
}

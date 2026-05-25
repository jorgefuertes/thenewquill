package process

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/pkg/validator"
)

func (p Process) Validate() error {
	if err := validator.Validate(p); err != nil {
		return err
	}

	return nil
}

func (s *Service) ValidateAll() []error {
	validationErrors := []error{}

	for _, k := range AllowedTables() {
		if !s.Exists(k) {
			validationErrors = append(validationErrors, fmt.Errorf("missing process table: %q", k.String()))
		}
	}

	return validationErrors
}

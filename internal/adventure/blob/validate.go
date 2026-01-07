package blob

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/pkg/validator"
)

func (b Blob) Validate() error {
	if err := validator.Validate(b); err != nil {
		return err
	}

	return nil
}

func (s *Service) ValidateAll() []error {
	validationErrors := []error{}

	blobs := s.db.Query(database.FilterByKind(kind.Blob))
	defer blobs.Close()

	var b Blob
	for blobs.Next(&b) {
		if err := b.Validate(); err != nil {
			validationErrors = append(
				validationErrors,
				fmt.Errorf("%w: blob %d %q", err, b.ID, s.db.GetLabelOrBlank(b.LabelID)),
			)
		}
	}

	return validationErrors
}

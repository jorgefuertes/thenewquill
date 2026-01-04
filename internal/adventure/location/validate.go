package location

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/pkg/validator"
)

func (l Location) Validate() error {
	if err := validator.Validate(l); err != nil {
		return err
	}

	for _, conn := range l.Conns {
		if err := validator.Validate(conn); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) ValidateAll() []error {
	validationErrors := []error{}

	locations := s.db.Query(database.FilterByKind(kind.Location))
	defer locations.Close()

	var loc Location
	for locations.Next(&loc) {
		if err := loc.Validate(); err != nil {
			validationErrors = append(
				validationErrors,
				fmt.Errorf("%w: location %d %q", err, loc.ID, s.db.GetLabelOrBlank(loc.LabelID)),
			)
		}

		for i, conn := range loc.Conns {
			connErr := fmt.Errorf(
				"[LOC:%d:%s], [CONN#%d]: [WORD:%d:%s]->[LOC:%d:%s]",
				loc.ID,
				s.db.GetLabelOrBlank(loc.LabelID),
				i,
				conn.WordID,
				s.db.GetLabelFromRecordOrBlank(conn.WordID),
				conn.LocationID,
				s.db.GetLabelFromRecordOrBlank(conn.LocationID),
			)

			wordExists := s.db.Query(database.FilterByID(conn.WordID), database.FilterByKind(kind.Word)).Exists()
			locationExists := s.db.Query(database.FilterByID(conn.LocationID), database.FilterByKind(kind.Location)).
				Exists()

			if !wordExists {
				validationErrors = append(validationErrors, fmt.Errorf("%w: %w", ErrConnWordNotFound, connErr))
			}

			if !locationExists {
				validationErrors = append(validationErrors, fmt.Errorf("%w: %w", ErrConnLocationNotFound, connErr))
			}
		}
	}

	return validationErrors
}

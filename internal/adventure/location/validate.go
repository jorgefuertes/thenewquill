package location

import (
	"errors"
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

func (s *Service) ValidateAll() error {
	locations := s.db.Query(database.FilterByKind(kind.Location))
	defer locations.Close()

	var loc Location
	for locations.Next(&loc) {
		if err := loc.Validate(); err != nil {
			return errors.Join(err, fmt.Errorf("location %d: %s", loc.ID, s.db.GetLabelOrBlank(loc.LabelID)))
		}

		for i, conn := range loc.Conns {
			connErr := fmt.Errorf(
				"location %q, connection #%d: %q-->%q",
				s.db.GetLabelOrBlank(loc.LabelID),
				i,
				s.db.GetLabelFromRecordOrBlank(conn.WordID),
				s.db.GetLabelFromRecordOrBlank(conn.LocationID),
			)

			wordExists := s.db.Query(database.FilterByID(conn.WordID), database.FilterByKind(kind.Word)).Exists()
			locationExists := s.db.Query(database.FilterByID(conn.LocationID), database.FilterByKind(kind.Location)).
				Exists()

			if !wordExists {
				return fmt.Errorf("%w: %w", ErrConnWordNotFound, connErr)
			}

			if !locationExists {
				return fmt.Errorf("%w: %w", ErrConnLocationNotFound, connErr)
			}
		}
	}

	return nil
}

package compiler

import (
	"fmt"
	"os"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
	"github.com/jorgefuertes/thenewquill/internal/database"
)

func validateSection(a *adventure.Adventure, s *status.Status, k kind.Kind) {
	validators := map[kind.Kind]func() error{
		kind.Param:     a.Config.ValidateAll,
		kind.Word:      a.Words.ValidateAll,
		kind.Message:   a.Messages.ValidateAll,
		kind.Variable:  a.Variables.ValidateAll,
		kind.Item:      a.Items.ValidateAll,
		kind.Character: a.Characters.ValidateAll,
		kind.Location:  a.Locations.ValidateAll,
	}

	if err := validators[k](); err != nil {
		fmt.Printf("❗ Validation error in section %q\n", s.Section.String())
		fmt.Println(err)

		os.Exit(1)
	}
}

func replaceLocationConnectionsIDs(a *adventure.Adventure) {
	c := a.DB.Query(database.FilterByKind(kind.Location))
	defer c.Close()

	var loc location.Location
	for c.Next(&loc) {
		for i, conn := range loc.Conns {
			d, err := a.Locations.Get().WithLabelID(conn.LocationID).First()
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"❌ cannot find location %q with label %q\n",
					conn.LocationID,
					a.DB.GetLabelOrBlank(conn.LocationID),
				)
				os.Exit(1)
			}

			conn.LocationID = d.ID
			loc.Conns[i] = conn
		}

		if err := a.DB.Update(&loc); err != nil {
			fmt.Fprintf(os.Stderr, "❌ cannot update location %q: %s\n", a.DB.GetLabelOrBlank(loc.LabelID), err)
			os.Exit(1)
		}
	}
}

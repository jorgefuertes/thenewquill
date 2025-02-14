package adventure

import (
	"fmt"

	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/msg"
	"thenewquill/internal/adventure/obj"
	"thenewquill/internal/adventure/vars"
	"thenewquill/internal/adventure/voc"
	"thenewquill/internal/util"
)

type Adventure struct {
	Vars       vars.Store
	Vocabulary voc.Vocabulary
	Messages   msg.Messages
	Locations  loc.Locations
	Objects    obj.Items
}

func New() *Adventure {
	return &Adventure{
		Vars:       vars.New(),
		Vocabulary: voc.New(),
		Messages:   msg.New(),
		Locations:  loc.New(),
		Objects:    obj.NewItems(),
	}
}

func (a *Adventure) Dump() {
	fmt.Println("--- VARS BEGIN ---")
	for k, v := range a.Vars.All() {
		fmt.Printf("[var] %s = %v\n", k, v)
	}
	fmt.Println("--- VARS END ---")

	fmt.Println("--- MESSAGES BEGIN ---")
	for _, m := range a.Messages {
		fmt.Printf("[%s msg] %s: %s\n", m.Type, m.Label, m.Text)
	}
	fmt.Println("--- MESSAGES END ---")

	fmt.Println("--- LOCATIONS BEGIN ---")
	for _, l := range a.Locations.All() {
		fmt.Printf("[loc:%-5s] %-10s \"%s\"\n",
			util.LimitStr(l.Label, 5),
			util.LimitStr(l.Title, 10),
			util.LimitStr(l.Description, 50))

		if len(l.Conns) > 0 {
			fmt.Printf("            exits: ")
			for _, c := range l.Conns {
				fmt.Printf("%s->%s ", c.Word, c.To.Label)
			}
			fmt.Println()
		}
	}
	fmt.Println("--- LOCATIONS END ---")
}

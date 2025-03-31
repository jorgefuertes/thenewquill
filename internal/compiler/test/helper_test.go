package compiler_test

import (
	"reflect"
	"testing"

	"thenewquill/internal/adventure"
	"thenewquill/internal/adventure/character"
	"thenewquill/internal/adventure/config"
	"thenewquill/internal/adventure/item"
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/msg"
	"thenewquill/internal/adventure/vars"
	"thenewquill/internal/adventure/words"

	"github.com/stretchr/testify/assert"
)

func assertEqualAventures(t *testing.T, expected, actual *adventure.Adventure) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	if expected == nil {
		t.Error("expected adventure is nil")

		return
	}

	if actual == nil {
		t.Error("actual adventure is nil")

		return
	}

	assertEqualConfig(t, expected.Config, actual.Config)
	assertEqualVars(t, expected.Vars, actual.Vars)
	assertEqualWords(t, expected.Words, actual.Words)
	assertEqualMessages(t, expected.Messages, actual.Messages)
	assertEqualLocations(t, expected.Locations, actual.Locations)
	assertEqualItems(t, expected.Items, actual.Items)
	assertEqualChars(t, expected.Chars, actual.Chars)
}

func assertEqualConfig(t *testing.T, expected, actual config.Config) {
	t.Helper()

	if expected.Title != actual.Title {
		t.Errorf("expected title %s, got %s", expected.Title, actual.Title)
	}

	if expected.Author != actual.Author {
		t.Errorf("expected author %s, got %s", expected.Author, actual.Author)
	}

	if expected.Description != actual.Description {
		t.Errorf("expected description %s, got %s", expected.Description, actual.Description)
	}

	if expected.Version != actual.Version {
		t.Errorf("expected version %s, got %s", expected.Version, actual.Version)
	}

	if expected.Date != actual.Date {
		t.Errorf("expected date %s, got %s", expected.Date, actual.Date)
	}

	if expected.Lang != actual.Lang {
		t.Errorf("expected lang %s, got %s", expected.Lang, actual.Lang)
	}
}

func assertEqualVars(t *testing.T, expected, actual vars.Store) {
	t.Helper()

	if expected.Len() != actual.Len() {
		t.Errorf("expected vars %d, got %d", expected.Len(), actual.Len())
	}

	for k, v := range expected.Regs {
		if actual.Regs[k] != v {
			t.Errorf("expected var %s = %v, got %v", k, v, actual.Regs[k])
		}

		expectedType := reflect.TypeOf(v)
		actualType := reflect.TypeOf(actual.Regs[k])

		if expectedType != actualType {
			t.Errorf("expected var %s type %s, got %s", k, expectedType, actualType)
		}
	}
}

func assertEqualWords(t *testing.T, expected, actual words.Store) {
	t.Helper()

	if expected.Len() != actual.Len() {
		t.Errorf("expected words %d, got %d", expected.Len(), actual.Len())
	}

	for _, w := range expected {
		if !actual.Exists(w.Type, w.Label) {
			t.Errorf("expected %s with label %s not found in actual", w.Type.String(), w.Label)

			continue
		}

		w2 := actual.Get(w.Type, w.Label)
		assert.ElementsMatch(t, w.Synonyms, w2.Synonyms, "synonyms for %s %s doesn't match", w.Type.String(), w.Label)
	}
}

func assertEqualMessages(t *testing.T, expected, actual msg.Store) {
	t.Helper()

	if expected.Len() != actual.Len() {
		t.Errorf("expected messages %d, got %d", expected.Len(), actual.Len())
	}

	for _, m := range expected {
		if !actual.Exists(m.Label) {
			t.Errorf("expected message %s not found in actual", m.Label)

			continue
		}

		m2 := actual.Get(m.Label)
		assert.Equal(t, m.Text, m2.Text, "text for %s doesn't match", m.Label)
		assert.Equal(t, m.Plurals, m2.Plurals, "plurals for %s doesn't match", m.Label)
	}
}

func assertEqualLocations(t *testing.T, expected, actual loc.Store) {
	t.Helper()

	if expected.Len() != actual.Len() {
		t.Errorf("expected locations %d, got %d", expected.Len(), actual.Len())
	}

	for _, l := range expected.GetAll() {
		if !actual.Exists(l.Label) {
			t.Errorf("expected location %s not found in actual", l.Label)

			continue
		}

		l2 := actual.Get(l.Label)
		assert.Equal(t, l.Title, l2.Title, "location %s.Title doesn't match", l.Label)
		assert.Equal(t, l.Description, l2.Description, "location %s.Description doesn't match", l.Label)
		assert.Equal(t, len(l.Conns), len(l2.Conns), "location %s.Conns length differs", l.Label)
		for _, c := range l.Conns {
			to2 := l2.GetConn(c.Word)
			assert.NotNil(t, to2, "connection %s for %s doesn't exist", c.Word.Label, l.Label)
			assert.Equal(t, c.To.Label, to2.Label, "connection %s for %s doesn't match", c.Word.Label, l.Label)
		}

		assert.Equal(t, l.Vars.Len(), l2.Vars.Len(), "vars for %s differ in length", l.Label)
		for k, v := range l.Vars.Regs {
			assert.Equal(t, v, l2.Vars.Regs[k], "var %s for %s doesn't match", k, l.Label)
		}
	}
}

func assertEqualItems(t *testing.T, expected, actual item.Store) {
	t.Helper()

	if expected.Len() != actual.Len() {
		t.Errorf("expected items %d, got %d", expected.Len(), actual.Len())
	}

	for _, i := range expected.GetAll() {
		if !actual.Exists(i.Label) {
			t.Errorf("expected item %s not found in actual", i.Label)

			continue
		}

		i2 := actual.Get(i.Label)
		assert.Equal(t, i.Noun.Label, i2.Noun.Label, "item %s.Noun doesn't match", i.Label)
		assert.Equal(t, i.Adjective.Label, i2.Adjective.Label, "item %s.Adjective doesn't match", i.Label)
		assert.Equal(t, i.Description, i2.Description, "item %s.Description doesn't match", i.Label)
		assert.Equal(t, i.Weight, i2.Weight, "item %s.Weight doesn't match", i.Label)
		assert.Equal(t, i.MaxWeight, i2.MaxWeight, "item %s.MaxWeight doesn't match", i.Label)
		assert.Equal(t, i.IsContainer, i2.IsContainer, "item %s.IsContainer doesn't match", i.Label)
		assert.Equal(t, i.IsWearable, i2.IsWearable, "item %s.IsWearable doesn't match", i.Label)
		assert.Equal(t, i.IsCreated, i2.IsCreated, "item %s.IsCreated doesn't match", i.Label)

		if i.Location != nil {
			assert.Equal(t, i.Location.Label, i2.Location.Label, "item %s.Location doesn't match", i.Label)
		} else {
			assert.Nil(t, i2.Location, "item %s.Location should be nil", i.Label)
		}

		if i.Inside != nil {
			assert.Equal(t, i.Inside.Label, i2.Inside.Label, "item %s.Inside doesn't match", i.Label)
		} else {
			assert.Nil(t, i2.Inside, "item %s.Inside should be nil", i.Label)
		}

		if i.CarriedBy != nil {
			assert.Equal(t, i.CarriedBy.Label, i2.CarriedBy.Label, "item %s.CarriedBy doesn't match", i.Label)
		} else {
			assert.Nil(t, i2.CarriedBy, "item %s.CarriedBy should be nil", i.Label)
		}

		if i.WornBy != nil {
			assert.Equal(t, i.WornBy.Label, i2.WornBy.Label, "item %s.WornBy doesn't match", i.Label)
		} else {
			assert.Nil(t, i2.WornBy, "item %s.WornBy should be nil", i.Label)
		}

		assert.Equal(t, len(i.Vars.Regs), len(i2.Vars.Regs), "item %s.Vars length differs", i.Label)
		for k, v := range i.Vars.Regs {
			assert.Equal(t, v, i2.Vars.Regs[k], "item %s.Vars[%s] doesn't match", k, i.Label)
		}
	}
}

func assertEqualChars(t *testing.T, expected, actual character.Store) {
	t.Helper()

	if expected.Len() != actual.Len() {
		t.Errorf("expected chars %d, got %d", expected.Len(), actual.Len())
	}

	for _, c := range expected.GetAll() {
		if !actual.Exists(c.Label) {
			t.Errorf("expected char %s not found in actual", c.Label)

			continue
		}

		c2 := actual.Get(c.Label)
		assert.Equal(t, c.Name.Label, c2.Name.Label, "name for %s doesn't match", c.Label)
		assert.Equal(t, c.Adjective.Label, c2.Adjective.Label, "adjective for %s doesn't match", c.Label)
		assert.Equal(t, c.Description, c2.Description, "description for %s doesn't match", c.Label)

		if c.Location != nil {
			assert.Equal(t, c.Location.Label, c2.Location.Label, "location for %s doesn't match", c.Label)
		} else {
			assert.Nil(t, c2.Location, "location for %s should be nil", c.Label)
		}

		assert.Equal(t, c.Created, c2.Created, "created for %s doesn't match", c.Label)
		assert.Equal(t, c.Human, c2.Human, "human for %s doesn't match", c.Label)

		assert.Equal(t, len(c.Vars.Regs), len(c2.Vars.Regs), "vars for %s differ in length", c.Label)
		for k, v := range c.Vars.Regs {
			assert.Equal(t, v, c2.Vars.Regs[k], "var %s for %s doesn't match", k, c.Label)
		}
	}
}

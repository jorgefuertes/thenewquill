package primitive

import (
	"errors"
	"fmt"
	"regexp"
)

type Label string

var (
	Underscore Label = "_"
	Wildcard   Label = "*"
)

func New(name string) Label {
	return Label(name)
}

func (l Label) String() string {
	return string(l)
}

func LabelFromLabelOrString(labelOrString any) (Label, error) {
	switch v := any(labelOrString).(type) {
	case Label:
		return v, nil
	case string:
		return Label(v), nil
	default:
		return Label(""), errors.Join(ErrInvalidLabelType, fmt.Errorf("%T:%v", v, v))
	}
}

func (l Label) Validate(allowComposite bool) error {
	validLabelRg := regexp.MustCompile(`^[\d\p{L}\-_]{1,25}$`)
	validComposedLabelRg := regexp.MustCompile(`^[\d\p{L}\-_\.]{1,50}$`)

	if allowComposite && validComposedLabelRg.MatchString(l.String()) {
		return nil
	}

	if validLabelRg.MatchString(l.String()) {
		return nil
	}

	return errors.Join(ErrInvalidLabelName, errors.New(l.String()))
}

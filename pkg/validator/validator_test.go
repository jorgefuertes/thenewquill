package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	type TestStruct struct {
		Label  string   `valid:"matches(^[\\d\\p{L}\\-_]+$)"`
		Name   string   `valid:"required"`
		Age    int      `valid:"required,min=18,max=99"`
		City   string   `valid:"in=zaragoza|madrid|barcelona"`
		Phone  string   `valid:"numeric"`
		Emails []string `valid:"len(6|50),count(1|4),matches(^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$)"`
		Syns   []string `valid:"len(1|25),count(1|25)"`
	}

	testCases := []struct {
		name     string
		input    TestStruct
		expected error
	}{
		{
			name: "valid",
			input: TestStruct{
				Label: "label",
				Name:  "John Doe",
				City:  "zaragoza",
				Age:   25,
				Emails: []string{
					"john@doe.com",
					"john.doe@gmail.com",
				},
				Syns: []string{"john", "doe", "y", "z"},
			},
			expected: nil,
		},
		{
			name: "no syns",
			input: TestStruct{
				Label: "label",
				Name:  "John Doe",
				City:  "zaragoza",
				Age:   25,
				Emails: []string{
					"john@doe.com",
					"john.doe@gmail.com",
				},
				Syns: []string{},
			},
			expected: errors.New("validating TestStruct: Syns must have at least 1 elements"),
		},
		{
			name: "long syns",
			input: TestStruct{
				Label: "label",
				Name:  "John Doe",
				City:  "zaragoza",
				Age:   25,
				Emails: []string{
					"john@doe.com",
					"john.doe@gmail.com",
				},
				Syns: []string{
					"john", "doe", "y", "z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l",
					"m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
				},
			},
			expected: errors.New("validating TestStruct: Syns must have less or equal than 25 elements"),
		},
		{
			name: "a long syn",
			input: TestStruct{
				Label: "label",
				Name:  "John Doe",
				City:  "zaragoza",
				Age:   25,
				Emails: []string{
					"john@doe.com",
					"john.doe@gmail.com",
				},
				Syns: []string{"john", "doe", "thisisaverylongsynonymmorethan25characters"},
			},
			expected: errors.New(
				"validating TestStruct: Syns[2]: thisisaverylongsynonymmorethan25characters must have less or equal than 25 characters",
			),
		},
		{
			name: "missing name",
			input: TestStruct{
				Label: "label",
				Name:  "",
				Age:   25,
				City:  "zaragoza",
				Emails: []string{
					"john@doe.com",
					"john.doe@gmail.com",
				},
				Syns: []string{"john", "doe", "y", "z"},
			},
			expected: errors.New("validating TestStruct: Name is required"),
		},
		{
			name: "missing age",
			input: TestStruct{
				Label: "label",
				Name:  "Noage Jhonson",
				Age:   0,
				City:  "zaragoza",
				Emails: []string{
					"john@doe.com",
					"john.doe@gmail.com",
				},
				Syns: []string{"john", "doe", "y", "z"},
			},
			expected: errors.New("validating TestStruct: Age is required"),
		},
		{
			name: "underage",
			input: TestStruct{
				Label: "label",
				Name:  "Underage Jhonson",
				Age:   5,
				City:  "zaragoza",
				Emails: []string{
					"john@doe.com",
					"john.doe@gmail.com",
				},
				Syns: []string{"john", "doe", "y", "z"},
			},
			expected: errors.New("validating TestStruct: Age must be at least 18"),
		},
		{
			name: "on age",
			input: TestStruct{
				Label: "label",
				Name:  "Onage Jhonson",
				Age:   20,
				City:  "zaragoza",
				Emails: []string{
					"john@doe.com",
					"john.doe@gmail.com",
				},
				Syns: []string{"john", "doe", "y", "z"},
			},
			expected: nil,
		},
		{
			name: "overage",
			input: TestStruct{
				Label: "label",
				Name:  "Underage Jhonson",
				Age:   100,
				City:  "zaragoza",
				Emails: []string{
					"john@doe.com",
					"john.doe@gmail.com",
				},
				Syns: []string{"john", "doe", "y", "z"},
			},
			expected: errors.New("validating TestStruct: Age must be 99 as maximum"),
		},
		{
			name: "invalid label",
			input: TestStruct{
				Label: "label not valid",
				Name:  "Jhon Jhonson",
				Age:   19,
				City:  "zaragoza",
				Emails: []string{
					"john@doe.com",
					"john.doe@gmail.com",
				},
				Syns: []string{"john", "doe", "y", "z"},
			},
			expected: errors.New("validating TestStruct: Label does not match regexp `^[\\d\\p{L}\\-_]+$`"),
		},
		{
			name: "valid label",
			input: TestStruct{
				Label: "labelvalid",
				Name:  "Jhon Jhonson",
				Age:   19,
				City:  "zaragoza",
				Emails: []string{
					"john@doe.com",
					"john.doe@gmail.com",
				},
				Syns: []string{"john", "doe", "y", "z"},
			},
			expected: nil,
		},
		{
			name: "invalid city",
			input: TestStruct{
				Label: "label",
				Name:  "Jhon Jhonson",
				Age:   19,
				City:  "gotham",
				Emails: []string{
					"john@doe.com",
					"john.doe@gmail.com",
				},
				Syns: []string{"john", "doe", "y", "z"},
			},
			expected: errors.New("validating TestStruct: City must be one of zaragoza, madrid, barcelona"),
		},
		{
			name: "no emails",
			input: TestStruct{
				Label: "label",
				Name:  "John Doe",
				City:  "zaragoza",
				Age:   25,
				Syns:  []string{"john", "doe", "y", "z"},
			},
			expected: errors.New("validating TestStruct: Emails must have at least 1 elements"),
		},
		{
			name: "too many emails",
			input: TestStruct{
				Label: "label",
				Name:  "John Doe",
				City:  "zaragoza",
				Age:   25,
				Emails: []string{
					"john@doe.com",
					"john.doe@gmail.com",
					"john@doe.com",
					"john.doe@gmail.com",
					"john@doe.com",
					"john.doe@gmail.com",
				},
				Syns: []string{"john", "doe", "y", "z"},
			},
			expected: errors.New("validating TestStruct: Emails must have less or equal than 4 elements"),
		},
		{
			name: "short email",
			input: TestStruct{
				Label: "label",
				Name:  "John Doe",
				City:  "zaragoza",
				Age:   25,
				Emails: []string{
					"john@doe.com",
					"john.doe@gmail.com",
					"j@d.c",
				},
				Syns: []string{"john", "doe", "y", "z"},
			},
			expected: errors.New("validating TestStruct: Emails[2]: j@d.c must have at least 6 characters"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.input)

			if tc.expected != nil {
				require.Error(t, err)
				require.Equal(t, tc.expected, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

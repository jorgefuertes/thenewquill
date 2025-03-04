package bin

import (
	"fmt"
	"strings"

	"thenewquill/internal/adventure/words"
	"thenewquill/internal/compiler/section"
)

type label struct {
	Name string
	Sec  section.Section
}

type labels []label

func (s *labels) Add(name string, sec section.Section) int {
	if !s.Exists(name, sec) {
		*s = append(*s, label{Name: name, Sec: sec})
	}

	id, _ := s.GetID(name, sec)

	return id
}

func (s *labels) GetID(name string, sec section.Section) (int, bool) {
	for i, l := range *s {
		if l.Name == name && l.Sec == sec {
			return i, true
		}
	}

	return -1, false
}

func (s *labels) Exists(name string, sec section.Section) bool {
	for _, l := range *s {
		if l.Name == name && l.Sec == sec {
			return true
		}
	}

	return false
}

func (s labels) Get(id int) (label, bool) {
	if id < 0 || id >= len(s) {
		return label{}, false
	}

	return s[id], true
}

func (s labels) GetSection(sec section.Section) map[int]string {
	labels := make(map[int]string, 0)

	for i, l := range s {
		if l.Sec == sec {
			labels[i] = l.Name
		}
	}

	return labels
}

func composeWordName(name string, t words.WordType) string {
	return fmt.Sprintf("%s#%s", name, t.String())
}

func splitWordName(name string) (string, words.WordType) {
	parts := strings.Split(name, "#")
	if len(parts) != 2 {
		return "", words.Unknown
	}

	t := words.WordTypeFromString(parts[1])
	if t == words.Unknown {
		return "", words.Unknown
	}

	return parts[0], t
}

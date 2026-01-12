package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"github.com/jorgefuertes/thenewquill/pkg/log"
)

type Lang string

const (
	EN Lang = "en"
	ES Lang = "es"
)

type SL struct {
	Adverb      []*word.Word
	Verb        []*word.Word
	Adjective1  []*word.Word
	Noun1       []*word.Word
	Preposition []*word.Word
	Adjective2  []*word.Word
	Noun2       []*word.Word
	Executed    bool
}

type Parser struct {
	wordStore *word.Service
	lang      Lang
	lastInput string
	splitRg   *regexp.Regexp
	history   []SL
}

func New(wordStore *word.Service, lang Lang) (*Parser, error) {
	words := wordStore.Get().WithType(word.Conjunction).All()

	conjunctions := ""
	for _, w := range words {
		for _, syn := range w.Synonyms {
			if strings.Contains("_*", syn) {
				continue
			}

			conjunctions += `|\b` + regexp.QuoteMeta(syn) + `\b`
		}
	}

	splitStr := `(?i)(?:\.|,|;|¡|!|¿|\?|\n` + conjunctions + `)+`
	splitRg := regexp.MustCompile(splitStr)
	fmt.Printf("Split regex: %s\n", splitRg.String())

	return &Parser{
		wordStore: wordStore,
		lang:      lang,
		lastInput: "",
		splitRg:   splitRg,
		history:   []SL{},
	}, nil
}

func (p *Parser) Parse(input string) {
	var phrases []string

	parts := p.splitRg.Split(strings.TrimSpace(input), -1)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			phrases = append(phrases, part)
		}
	}

	log.Debug("Parsed phrases: %+v", phrases)
	_ = phrases
}

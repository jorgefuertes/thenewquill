package parser

import (
	"regexp"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"github.com/jorgefuertes/thenewquill/internal/lang"
)

type Parser struct {
	wordStore *word.Service
	lang      lang.Lang
	splitRg   *regexp.Regexp
	Sentences []LS
	cursor    int
}

func New(wordStore *word.Service) (*Parser, error) {
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

	l := wordStore.GetLang()

	return &Parser{
		wordStore: wordStore,
		lang:      l,
		splitRg:   splitRg,
		Sentences: []LS{},
		cursor:    -1,
	}, nil
}

func (p *Parser) Parse(input string) {
	var phrases []string

	parts := p.splitRg.Split(strings.TrimSpace(input), -1)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.ToLower(part)
		if part != "" {
			phrases = append(phrases, part)
		}
	}

	p.transformToLogicalSentences(phrases)
}

func (p *Parser) Reset() {
	p.Sentences = []LS{}
	p.cursor = -1
}

func (p *Parser) Len() int {
	return len(p.Sentences)
}

func (p *Parser) transformToLogicalSentences(phrases []string) {
	// TODO: implement talking with other characters when text is between quotes
	// 1st pass - tokenize phrases into words
	for _, phrase := range phrases {
		ls := NewLS()
		tokens := strings.Split(phrase, " ")
		for _, token := range tokens {
			w, err := p.wordStore.GetAnyWith(token, word.Verb, word.Adverb, word.Noun, word.Adjective)
			if err == nil {
				ls.addWord(w)
			}
		}

		if ls.isEmpty() {
			continue
		}

		p.Sentences = append(p.Sentences, ls)
	}

	// 2nd pass - add implied verbs if missing
	for _, ls := range p.Sentences {
		if !ls.Has(word.Verb) && ls.Has(word.Noun) {
			noun := ls.Get(word.Noun, First)
			if noun == nil {
				continue
			}

			if noun.IsConnection {
				verb := p.wordStore.GetDefaultVerbSyns(p.lang, lang.Go)
				if verb != nil {
					ls.addWordAt(verb, 0)
				}

				return
			}

			if noun.IsItem {
				verb := p.wordStore.GetDefaultVerbSyns(p.lang, lang.Examine)
				if verb != nil {
					ls.setVerb(verb)
				}

				return
			}

			if noun.IsCharacter {
				verb := p.wordStore.GetDefaultVerbSyns(p.lang, lang.Talk)
				if verb != nil {
					ls.setVerb(verb)
				}

				return
			}
		}
	}
}

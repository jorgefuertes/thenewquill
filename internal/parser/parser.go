package parser

import (
	"regexp"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"github.com/jorgefuertes/thenewquill/internal/lang"
	"github.com/jorgefuertes/thenewquill/internal/util"
)

var (
	phraseSeparatorRg        = regexp.MustCompile(`\s+|,|\.|;|:`)
	quoteRg                  = regexp.MustCompile(`["""''']`)
	exclamationAndQuestionRg = regexp.MustCompile(`[¡!¿?]`)
)

type Parser struct {
	wordStore *word.Service
	lang      lang.Lang
	Sentences []LS
	cursor    int
}

func New(wordStore *word.Service) (*Parser, error) {
	return &Parser{
		wordStore: wordStore,
		lang:      wordStore.GetLang(),
		Sentences: []LS{NewLS()},
		cursor:    -1,
	}, nil
}

func (p *Parser) getCurrentLS() LS {
	return p.Sentences[len(p.Sentences)-1]
}

func (p *Parser) setCurrentLS(ls LS) {
	p.Sentences[len(p.Sentences)-1] = ls
}

func (p *Parser) closeLS() {
	ls := p.getCurrentLS()
	if ls.isEmpty() {
		ls.original = strings.TrimSpace(ls.original)
		p.setCurrentLS(ls)

		return
	}

	ls.original = strings.TrimSpace(ls.original)
	p.setCurrentLS(ls)

	ls2 := NewLS()
	ls2.sub = ls.sub

	p.Sentences = append(p.Sentences, ls2)
}

func (p *Parser) removeEmptyLS() {
	if p.getCurrentLS().isEmpty() {
		p.Sentences = p.Sentences[:len(p.Sentences)-1]
	}
}

func (p *Parser) addWord(current string) {
	ls := p.getCurrentLS()

	needle := util.NormalizeString(current)
	if needle == "" {
		return
	}

	w, err := p.wordStore.GetAnyWith(needle, word.Verb, word.Noun, word.Pronoun, word.Adjective, word.Adverb,
		word.Conjunction, word.Preposition)
	if err != nil {
		return
	}

	if w.Type == word.Conjunction {
		if !ls.isEmpty() {
			p.remove(current)
			p.closeLS()

			return
		}
	}

	ls.addWord(w)
	p.setCurrentLS(ls)
}

func (p *Parser) remove(s string) {
	rg := regexp.MustCompile(`\s*` + regexp.QuoteMeta(s) + `\s*`)
	ls := p.getCurrentLS()
	ls.original = rg.ReplaceAllString(ls.original, " ")
	p.setCurrentLS(ls)
}

func (p *Parser) Parse(input string) {
	input = strings.TrimSpace(input)
	current := ""

	for _, c := range input {
		ls := p.getCurrentLS()
		ls.original += string(c)
		p.setCurrentLS(ls)

		if exclamationAndQuestionRg.MatchString(string(c)) {
			p.addWord(current)
			p.remove(string(c))
			p.closeLS()
			current = ""
			continue
		}

		if phraseSeparatorRg.MatchString(string(c)) {
			p.addWord(current)
			p.remove(string(c))
			current = ""
			continue
		}

		if quoteRg.MatchString(string(c)) {
			p.addWord(current)
			isSub := p.getCurrentLS().IsSub()
			p.remove(string(c))
			p.closeLS()
			ls := p.getCurrentLS()
			ls.sub = !isSub
			p.setCurrentLS(ls)
			current = ""
			continue
		}

		current += string(c)
	}

	p.addWord(current)
	p.closeLS()
	p.removeEmptyLS()
	p.completeSentences()
}

func (p *Parser) Reset() {
	p.Sentences = []LS{}
	p.cursor = -1
}

func (p *Parser) Len() int {
	return len(p.Sentences)
}

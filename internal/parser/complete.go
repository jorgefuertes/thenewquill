package parser

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"github.com/jorgefuertes/thenewquill/internal/lang"
)

func inferAction(w *word.Word) lang.Action {
	if w.IsConnection {
		return lang.Go
	}

	if w.IsCharacter {
		return lang.Talk
	}

	return lang.Examine
}

func (p *Parser) completeSentences() {
	for i, ls := range p.Sentences {
		if len(ls.words) == 0 {
			continue
		}

		if ls.Has(word.Verb) {
			continue
		}

		for _, w := range ls.words {
			if w.Type == word.Verb {
				break
			}

			if w.Type == word.Noun {
				action := inferAction(w)
				verb, err := p.wordStore.GetDefaultVerbForAction(action)
				if err != nil {
					continue
				}

				ls.setVerb(verb)

				break
			}
		}

		p.Sentences[i] = ls
	}
}

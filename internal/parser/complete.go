package parser

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/item"
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

func (p *Parser) setBindings() {
	var lastItem *item.Item
	var lastNPC *character.Character

	for i, ls := range p.Sentences {
		if len(ls.words) == 0 {
			continue
		}

		if !ls.Has(word.Noun) {
			ls.Item = lastItem
			ls.NPC = lastNPC
			p.Sentences[i] = ls

			continue
		}

		noun := ls.Get(word.Noun, 1)
		adj := ls.Get(word.Adjective, 1)

		q := p.itemStore.Get().WithNameID(noun.ID)
		if adj != nil {
			q = q.WithAdjectiveID(adj.ID)
		}

		if q.Exists() {
			item, _ := q.First()
			lastItem = item
			ls.Item = item
			p.Sentences[i] = ls

			continue
		}

		q2 := p.charStore.Get().WithNameID(noun.ID)
		if adj != nil {
			q2 = q2.WithAdjectiveID(adj.ID)
		}

		if q2.Exists() {
			npc, _ := q2.First()
			lastNPC = npc
			ls.NPC = npc
			p.Sentences[i] = ls

			continue
		}

		ls.Item = lastItem
		ls.NPC = lastNPC
		p.Sentences[i] = ls
	}
}

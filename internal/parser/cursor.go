package parser

func (p *Parser) NextLS() LS {
	p.cursor++
	if p.cursor >= len(p.Sentences) {
		return LS{}
	}

	ls := p.Sentences[p.cursor]

	return ls
}

func (p *Parser) HasRemaining() bool {
	return p.cursor+1 < len(p.Sentences)
}

func (p *Parser) Current() *LS {
	if p.cursor < 0 || p.cursor >= len(p.Sentences) {
		return nil
	}

	return &p.Sentences[p.cursor]
}

package line

import "github.com/jorgefuertes/thenewquill/internal/compiler/rg"

func (l Line) IsCommentBegin() bool {
	return rg.CommentBegin.MatchString(l.text)
}

func (l Line) IsCommentEnd() bool {
	return rg.CommentEnd.MatchString(l.OptimizedText())
}

func (l Line) IsBlank() bool {
	return rg.Blank.MatchString(l.text)
}

func (l Line) IsOneLineComment() bool {
	return rg.OneLinecomment.MatchString(l.text)
}

func (l Line) IsMultilineBegin() bool {
	return rg.MultilineBegin.MatchString(l.text) || rg.Continue.MatchString(l.text)
}

func (l Line) IsMultilineEnd(isHeredoc bool) bool {
	if isHeredoc {
		return rg.MultilineEnd.MatchString(l.text)
	}

	return !rg.Continue.MatchString(l.text)
}

package line

import "github.com/jorgefuertes/thenewquill/internal/compiler/rg"

func (l Line) IsCommentBegin() bool {
	return rg.CommentBegin.MatchString(l.Text)
}

func (l Line) IsCommentEnd() bool {
	return rg.CommentEnd.MatchString(l.OptimizedText())
}

func (l Line) IsBlank() bool {
	return rg.Blank.MatchString(l.Text)
}

func (l Line) IsOneLineComment() bool {
	return rg.OneLinecomment.MatchString(l.Text)
}

func (l Line) IsMultilineBegin() bool {
	return rg.MultilineBegin.MatchString(l.Text) || rg.Continue.MatchString(l.Text)
}

func (l Line) IsMultilineEnd(isHeredoc bool) bool {
	if isHeredoc {
		return rg.MultilineEnd.MatchString(l.Text)
	}

	return !rg.Continue.MatchString(l.Text)
}

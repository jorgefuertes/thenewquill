package config

type Label string

func (l Label) String() string {
	return string(l)
}

const (
	UnknownLabel     Label = "unknown"
	TitleLabel       Label = "title"
	AuthorLabel      Label = "author"
	DescriptionLabel Label = "description"
	DescLabel        Label = "desc"
	VersionLabel     Label = "version"
	DateLabel        Label = "date"
	LangLabel        Label = "lang"
)

func Labels() []Label {
	return []Label{
		TitleLabel,
		AuthorLabel,
		DescriptionLabel,
		DescLabel,
		VersionLabel,
		DateLabel,
		LangLabel,
	}
}

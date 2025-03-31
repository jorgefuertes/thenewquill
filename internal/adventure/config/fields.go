package config

type Field string

func (l Field) String() string {
	return string(l)
}

const (
	UnknownField     Field = "unknown"
	TitleField       Field = "title"
	AuthorField      Field = "author"
	DescriptionField Field = "description"
	DescField        Field = "desc"
	VersionField     Field = "version"
	DateField        Field = "date"
	LangField        Field = "lang"
)

func Fields() []Field {
	return []Field{
		TitleField,
		AuthorField,
		DescriptionField,
		DescField,
		VersionField,
		DateField,
		LangField,
	}
}

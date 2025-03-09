package db

import "thenewquill/internal/compiler/section"

type Exportable interface {
	Export() map[section.Section][][]string
	ExportHeaders() []string
}

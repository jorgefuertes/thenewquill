package word

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/internal/lang"
	"github.com/jorgefuertes/thenewquill/pkg/log"
)

type Service struct {
	db  *database.DB
	cfg *config.Service
}

func NewService(d *database.DB, cfg *config.Service) *Service {
	s := &Service{db: d, cfg: cfg}

	labels := map[uint32]string{
		1: database.LabelAsterisk,
		2: database.LabelUnderscore,
	}

	for _, t := range WordTypes {
		for id, label := range labels {
			_, err := s.Create(&Word{
				LabelID:  id,
				Type:     t,
				Synonyms: []string{label},
			})
			if err != nil {
				panic(err)
			}
		}
	}

	return s
}

func (s *Service) DB() *database.DB {
	return s.db
}

func (s *Service) GetLang() lang.Lang {
	return lang.Lang(s.cfg.GetValueOrBlank(config.LanguageParamLabel))
}

func (s *Service) Create(w *Word) (uint32, error) {
	return s.db.Create(w)
}

func (s *Service) Update(w *Word) error {
	return s.db.Update(w)
}

func (s *Service) Count() int {
	return s.db.CountRecordsByKind(kind.Word)
}

func (s *Service) GetDefaultVerbForAction(a lang.Action) (*Word, error) {
	defaults := map[lang.Lang]map[lang.Action]string{
		lang.EN: {
			lang.Go:      "go",
			lang.Examine: "examine",
			lang.Talk:    "say",
		},
		lang.ES: {
			lang.Go:      "ir",
			lang.Examine: "examinar",
			lang.Talk:    "decir",
		},
	}

	name, ok := defaults[s.GetLang()][a]
	if !ok {
		log.Fatal("no default verb for action %d in language %q", a, s.GetLang().String())
	}

	w, err := s.Get().WithType(Verb).WithLabel(name).First()
	if err != nil {
		return nil, fmt.Errorf("you need to define an action verb with label %q", name)
	}

	return w, nil
}

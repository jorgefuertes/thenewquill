package word

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/internal/lang"
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

func (s *Service) GetDefaultVerbSyns(l lang.Lang, a lang.Action) *Word {
	for _, syn := range lang.GetDefaultSynonymForAction(l, a) {
		w, _ := s.Get().WithType(Verb).WithSynonym(syn).First()

		if w != nil {
			return w
		}
	}

	return nil
}

package process

type Process struct {
	Header   string
	VerbID   uint32
	NounID   uint32
	Condacts []Condact
}

func NewProcess(header string, verbID, nounID uint32) *Process {
	p := &Process{
		Header:   header,
		VerbID:   verbID,
		NounID:   nounID,
		Condacts: []Condact{},
	}

	return p
}

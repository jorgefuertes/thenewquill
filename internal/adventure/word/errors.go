package word

import "errors"

var (
	ErrDuplicatedLabel    = errors.New("duplicated label for type")
	ErrDuplicatedSyn      = errors.New("duplicated synonym for type")
	ErrWordNotFound       = errors.New("word not found")
	ErrMissingGoVerb      = errors.New("missing 'go' verb word")
	ErrMissingTalkVerb    = errors.New("missing 'talk' verb word")
	ErrMissingExamineVerb = errors.New("missing 'examine' verb word")
)

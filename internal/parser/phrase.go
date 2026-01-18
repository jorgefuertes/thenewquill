package parser

import (
	"regexp"
	"strings"
)

type phrase struct {
	str       string
	isTalking bool
}

func partToPhrases(part string) []phrase {
	var phrases []phrase
	// Remove empty quoted sections ("", '' or with only whitespace)
	emptyQuotedRg := regexp.MustCompile(`"\s*"|'\s*'|“\s*”|‘\s*’`)
	part = strings.TrimSpace(emptyQuotedRg.ReplaceAllString(part, " "))
	quotedRg := regexp.MustCompile(`"([^"]+)"|'([^']+)'|“([^”]+)”|‘([^’]+)’`)

	match := quotedRg.FindStringIndex(part)
	lastIndex := 0

	if match != nil {
		if match[0] > lastIndex {
			nonQuoted := strings.TrimSpace(part[lastIndex:match[0]])
			if nonQuoted != "" {
				phrases = append(phrases, phrase{str: strings.ToLower(nonQuoted), isTalking: false})
			}
		}

		quoted := strings.TrimSpace(part[match[0]:match[1]])
		if quoted != "" {
			quoted = strings.Trim(quoted, `"'“”‘’`)
			quoted = strings.TrimSpace(quoted)
			// Only add quoted phrase if there's at least one non-quoted phrase before it
			if quoted != "" && len(phrases) > 0 {
				phrases = append(phrases, phrase{str: strings.ToLower(quoted), isTalking: true})
			}
		}

		lastIndex = match[1]
	}

	if lastIndex < len(part) {
		nonQuoted := strings.TrimSpace(part[lastIndex:])
		if nonQuoted != "" {
			phrases = append(phrases, phrase{str: strings.ToLower(nonQuoted), isTalking: false})
		}
	}

	return phrases
}

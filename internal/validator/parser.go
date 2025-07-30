package validator

import (
	"fmt"
	"strings"
)

func parseValidTags(valid string) []string {
	preSplit := strings.Split(valid, ",")
	tags := []string{}
	unresolved := ""

	for i, current := range preSplit {
		for _, validator := range validators {
			if validator.tag == current {
				tags = append(tags, preSplit[i])

				continue
			}

			if validator.rg != nil {
				if validator.rg.MatchString(current) {
					tags = append(tags, preSplit[i])

					continue
				}
			}

			unresolved = unresolved + current
			if validator.rg != nil && validator.rg.MatchString(unresolved) {
				tags = append(tags, unresolved)
				unresolved = ""
			}
		}
	}

	return tags
}

func getValidatorAndParams(tag string) (validator, string, error) {
	v, err := getValidator(tag)
	if err != nil {
		return v, "", err
	}

	if v.rg == nil {
		return v, "", nil
	}

	matches := v.rg.FindStringSubmatch(tag)
	if len(matches) == 0 {
		return v, "", fmt.Errorf("invalid validator params %q", tag)
	}

	if v.tag == "matches" {
		params := strings.ReplaceAll(matches[1], `\\`, `\`)
		return v, params, nil
	}

	return v, matches[1], nil
}

func getValidator(tag string) (validator, error) {
	for _, validator := range validators {
		if strings.HasPrefix(tag, validator.tag) {
			return validator, nil
		}
	}

	return validator{}, fmt.Errorf("unknown validator %q", tag)
}

package regex

import (
	"regexp"
)

func ShouldIgnoreFile(filename string, regexPatterns *[]string) (bool, error) {
	shouldIgnore := false

	for _, pattern := range *regexPatterns {
		matched, err := regexp.MatchString(pattern, filename)
		if err != nil {
			return false, err
		}
		if matched {
			shouldIgnore = true
			break
		}
	}

	return shouldIgnore, nil
}

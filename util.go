package license

import (
	"regexp"
	"strings"
)

// scan will scan through a block of text, delmited by newline characters,
// and check for lines matching the provided glob text, case-insensitively.
func scan(text, match string) bool {
	// Case-insensitive matching
	text = strings.ToLower(text)
	match = strings.ToLower(match)

	re := regexp.MustCompile(match)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}

package naming

import (
	"regexp"
	"strings"
)

var (
	spaceRE = regexp.MustCompile(`\s`)
)

func Var(s string) string {
	s = strings.ReplaceAll(s, "'", "")
	return spaceRE.ReplaceAllString(s, "")
}

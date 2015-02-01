package targets

import "strings"

// GetToken extracts the next token from the current url part
func GetToken(urlPart *string) string {
	if index := strings.Index(*urlPart, "/"); index != -1 {
		s := *urlPart
		returnee := s[:strings.Index(s, "/")]
		*urlPart = s[strings.Index(s, "/")+1:]
		return returnee
	}

	return ""
}

package github

import "strings"

func getNextPageURL(linkHeader string) string {
	links := strings.Split(linkHeader, ",")
	for _, link := range links {
		parts := strings.Split(link, ";")
		if len(parts) < 2 {
			continue
		}
		// Look for the 'rel="next"' directive
		if strings.TrimSpace(parts[1]) == `rel="next"` {
			// Extract the URL between the angle brackets
			url := strings.TrimSpace(parts[0])
			url = strings.Trim(url, "<>")
			return url
		}
	}
	return ""
}

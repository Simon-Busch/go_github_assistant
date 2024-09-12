package github

import (
	"net/url"
	"strings"
)

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

func extractRepoAndOrg(repoURL string) (string, string) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return "", ""
	}
	parts := strings.Split(u.Path, "/")
	if len(parts) >= 3 {
		organization := parts[2] // Extract the organization
		repository := parts[3]   // Extract the repository name
		return organization, repository
	}
	return "", ""
}

package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Comment struct {
	Body      		string `json:"body"`
	User      		User `json:"user"`
	CreatedAt 		string `json:"created_at"`
}

type IssueOrPR struct {
	Title						string `json:"title"`
	Body						string `json:"body"`
	URL       			string `json:"html_url"`
	State     			string `json:"state"`
	CreatedAt 			string `json:"created_at"`
	UpdatedAt 			string `json:"updated_at"`
	CommentsURL 		string `json:"comments_url"`
	Repository 			string `json:"repository"`
	Organization 		string `json:"organization"`
	RepositoryURL 	string `json:"repository_url"`
}


type User struct {
	Login string `json:"login"`
}

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

func FetchComments(commentsURL, ghUsername, ghToken string) ([]Comment, error) {
	req, err := http.NewRequest("GET", commentsURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return nil, err
	}

	req.SetBasicAuth(ghUsername, ghToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: received non-200 response code: %d", resp.StatusCode)
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	var comments []Comment
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		log.Fatalf("Error decoding comments JSON: %v", err)
		return nil, err
	}

	return comments, nil
}

package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type IssuesResponse struct {
	TotalCount 				int     `json:"total_count"`
	Items      				[]Issue `json:"items"`
}

type Issue struct {
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

func FetchIssues(ghUsername, ghToken string) (*IssuesResponse, error) {
	query := fmt.Sprintf("assignee:%s", url.QueryEscape(ghUsername))
	apiURL := fmt.Sprintf("https://api.github.com/search/issues?q=%s&per_page=100", query)

	var allIssues IssuesResponse
	pageURL := apiURL

	for {
		req, err := http.NewRequest("GET", pageURL, nil)
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

		var issuesPage IssuesResponse
		if err := json.NewDecoder(resp.Body).Decode(&issuesPage); err != nil {
			log.Fatalf("Error decoding JSON: %v", err)
			return nil, err
		}

		for i, issue := range issuesPage.Items {
			org, repo := extractRepoAndOrg(issue.RepositoryURL)
			issuesPage.Items[i].Organization = org
			issuesPage.Items[i].Repository = repo
		}

		allIssues.Items = append(allIssues.Items, issuesPage.Items...)

		linkHeader := resp.Header.Get("Link")
		nextURL := getNextPageURL(linkHeader)
		if nextURL == "" {
			break // No more pages
		}
		pageURL = nextURL
	}

	return &allIssues, nil
}

func (i *IssuesResponse) GetAllOpenedIssues() []Issue {
	var openedIssues []Issue
	for _, issue := range i.Items {
		if issue.State == "open" {
			openedIssues = append(openedIssues, issue)
		}
	}
	return openedIssues
}

func (i *IssuesResponse) GetAllClosedIssues() []Issue {
	var closedIssues []Issue
	for _, issue := range i.Items {
		if issue.State == "closed" {
			closedIssues = append(closedIssues, issue)
		}
	}
	return closedIssues
}

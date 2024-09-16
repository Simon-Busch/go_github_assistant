package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type PullRequestResponse struct {
	TotalCount int           `json:"total_count"`
	Items      []IssueOrPR `json:"items"`
}

// type PullRequest struct {
// 	Title     			string `json:"title"`
// 	Body      			string `json:"body"`
// 	URL       			string `json:"html_url"`
// 	State     			string `json:"state"`
// 	CreatedAt 			string `json:"created_at"`
// 	UpdatedAt 			string `json:"updated_at"`
// 	CommentsURL 		string `json:"comments_url"`
// 	Repository 			string `json:"repository"`
// 	Organization 		string `json:"organization"`
// 	RepositoryURL 	string `json:"repository_url"`
// }

func FetchReviewRequests(ghUsername, ghToken string) (*PullRequestResponse, error) {
	query := fmt.Sprintf("review-requested:%s+state:open", url.QueryEscape(ghUsername))
	apiURL := fmt.Sprintf("https://api.github.com/search/issues?q=%s&per_page=100", query)

	var allPullRequests PullRequestResponse
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

		var pullRequestsPage PullRequestResponse
		if err := json.NewDecoder(resp.Body).Decode(&pullRequestsPage); err != nil {
			log.Fatalf("Error decoding JSON: %v", err)
			return nil, err
		}

		allPullRequests.Items = append(allPullRequests.Items, pullRequestsPage.Items...)

		linkHeader := resp.Header.Get("Link")
		nextURL := getNextPageURL(linkHeader)
		if nextURL == "" {
			break // No more pages
		}
		pageURL = nextURL
	}

	return &allPullRequests, nil
}

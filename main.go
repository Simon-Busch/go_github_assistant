package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/Simon-Busch/go_github_assistant/github"
	"github.com/Simon-Busch/go_github_assistant/utils"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	ghUserName, ghToken, err := utils.InitEnv()
	if err != nil {
		panic("please provide a .env file with GITHUB_USERNAME and GITHUB_TOKEN")
	}
	issuesResponse, err := github.FetchIssues(ghUserName, ghToken)
	if err != nil {
		log.Fatalf("Error fetching issues: %v", err)
	}
	issues := issuesResponse.Items

	if err := ui.Init(); err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	termWidth, termHeight := ui.TerminalDimensions()

	actionsTabs := widgets.NewTabPane("Github Assistant", "Issues", "Pull Requests", "Help")
	actionsTabs.SetRect(0, 0, termWidth, 3)

	issuesList := widgets.NewList()
	issuesList.Title = "Issues"
	issuesList.Rows = make([]string, len(issues))
	for i, issue := range issues {
		var color string
		if issue.State == "open" {
			color = "green"
		} else {
			color = "red"
		}
		issuesList.Rows[i] = fmt.Sprintf("[%d] [%s](fg:%s)", i+1, issue.Title, color)
	}
	issuesList.SelectedRowStyle.Fg = ui.ColorYellow
	issuesList.SetRect(0, 3, termWidth/2, termHeight-4)

	issueDetails := widgets.NewParagraph()
	issueDetails.Title = "Issue Details"
	issueDetails.Text = "Select an issue to see details."
	issueDetails.SetRect(termWidth, 3, 100, termHeight-4)

	ui.Render(actionsTabs, issuesList, issueDetails)

	selectedIndex := 0
	showComments := false
	commentsText := ""

	updateIssueDetails := func(index int, showComments bool) {
		issue := issues[index]
		issueText := fmt.Sprintf(
			"Title: %s\n\nRepository: %s\nOrganization: %s\n\nState: %s\n\nURL: %s\nCreated At: %s\n\nDescription:\n\n%s",
			issue.Title, issue.Repository, issue.Organization, issue.State, issue.URL, issue.CreatedAt, issue.Body)
		if showComments && commentsText != "" {
			issueText += fmt.Sprintf("\n\nComments:\n%s", commentsText)
		}
		issueDetails.Text = issueText
		ui.Render(issuesList, issueDetails)
	}


	updateIssueDetails(selectedIndex, showComments)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q":// Quit on 'q' or Ctrl+C
			return
		case "<Down>":
			if selectedIndex < len(issues)-1 {
				selectedIndex++
				issuesList.ScrollDown()
				updateIssueDetails(selectedIndex, showComments)
			}
			if selectedIndex == len(issues)-1 {
				selectedIndex = 0
				issuesList.ScrollTop()
				updateIssueDetails(selectedIndex, showComments)
			}
			commentsText = ""
		case "<Up>":
			if selectedIndex > 0 {
				selectedIndex--
				issuesList.ScrollUp()
				updateIssueDetails(selectedIndex, showComments)
			}
			if selectedIndex == 0 {
				selectedIndex = len(issues) - 1
				issuesList.ScrollBottom()
				updateIssueDetails(selectedIndex, showComments)
			}
			commentsText = ""
		case "<Enter>":
			issue := issues[selectedIndex]
			err := openBrowser(issue.URL)
			if err != nil {
				log.Printf("Failed to open browser: %v", err)
			}

		case "<C-c>":
			issue := issues[selectedIndex]
			if !showComments {
				// Fetch comments if not already fetched
				comments, err := github.FetchComments(issue.CommentsURL, ghUserName, ghToken)
				if err != nil {
					log.Printf("Error fetching comments for issue %s: %v", issue.URL, err)
					commentsText = "Error fetching comments."
				} else {
					// Build the comments text
					commentsText = ""
					for _, comment := range comments {
						commentsText += fmt.Sprintf("Comment by %s at %s:\n%s\n\n",
							comment.User.Login, comment.CreatedAt, comment.Body)
					}
				}
			}
			showComments = !showComments // Toggle comments visibility
			updateIssueDetails(selectedIndex, showComments)
		}
		ui.Render(issuesList, issueDetails)
	}
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin": // macOS
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

package main

import (
	"fmt"
	"log"

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

	issuesList := widgets.NewList()
	issuesList.Title = "Issues"
	issuesList.Rows = make([]string, len(issues))
	for i, issue := range issues {
		issuesList.Rows[i] = fmt.Sprintf("%d: %s", i+1, issue.Title)
	}
	issuesList.SelectedRowStyle.Fg = ui.ColorYellow
	issuesList.SetRect(0, 0, termWidth/2, termHeight)

	issueDetails := widgets.NewParagraph()
	issueDetails.Title = "Issue Details"
	issueDetails.Text = "Select an issue to see details."
	issueDetails.SetRect(termWidth, 0, 100, termHeight)

	ui.Render(issuesList, issueDetails)

	selectedIndex := 0

	updateIssueDetails := func(index int) {
		issue := issues[index]
		issueDetails.Text = fmt.Sprintf(
			"Title: %s\nState: %s\nURL: %s\nCreated At: %s\n\nDescription:\n%s",
			issue.Title, issue.State, issue.URL, issue.CreatedAt, issue.Body)
		ui.Render(issuesList, issueDetails)
	}

	updateIssueDetails(selectedIndex)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>": // Quit on 'q' or Ctrl+C
			return
		case "<Down>":
			if selectedIndex < len(issues)-1 {
				selectedIndex++
				issuesList.ScrollDown()
				updateIssueDetails(selectedIndex)
			}
			if selectedIndex == len(issues)-1 {
				selectedIndex = 0
				issuesList.ScrollTop()
				updateIssueDetails(selectedIndex)
			}
		case "<Up>":
			if selectedIndex > 0 {
				selectedIndex--
				issuesList.ScrollUp()
				updateIssueDetails(selectedIndex)
			}
			if selectedIndex == 0 {
				selectedIndex = len(issues) - 1
				issuesList.ScrollBottom()
				updateIssueDetails(selectedIndex)
			}
		case "<Enter>":
			issue := issues[selectedIndex]
			comments, _ := github.FetchComments(issue.CommentsURL, ghUserName, ghToken)
			commentsText := "Comments:\n"
			for _, comment := range comments {
				commentsText += fmt.Sprintf("- %s: %s\n", comment.User.Login, comment.Body)
			}
			issueDetails.Text += "\n\n" + commentsText
		}

		ui.Render(issuesList, issueDetails)
	}
}

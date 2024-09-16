package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Simon-Busch/go_github_assistant/github"
	"github.com/Simon-Busch/go_github_assistant/ui"
	"github.com/Simon-Busch/go_github_assistant/utils"
	termui "github.com/gizak/termui/v3"
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
	prToReview, err := github.FetchReviewRequests(ghUserName, ghToken)
	if err != nil {
		log.Fatalf("Error fetching PRs: %v", err)
	}

	openIssues := issuesResponse.GetAllOpenedIssues()
	closedIssues := issuesResponse.GetAllClosedIssues()

	var (
		selectedIndex = 0
		showComments = false
		commentsText = ""
		showHelper = false
		showClosedIssues = false
		showPr = false
		help *widgets.Paragraph
		issuesList *widgets.List
		currentIssues []github.IssueOrPR
	)

	if err := termui.Init(); err != nil {
		log.Fatal(err)
	}
	defer termui.Close()

	ui.RenderWaitingScreen(len(openIssues), len(closedIssues), ghUserName, prToReview.TotalCount)
	time.Sleep(1 * time.Second)
	termui.Clear()

	actionsTabs := ui.RenderHeader(ghUserName)
	issueDetails := ui.RenderIssueDetails()
	footer := ui.RenderFooter(len(openIssues), len(closedIssues), ghUserName, prToReview.TotalCount)
	if showClosedIssues {
		currentIssues = closedIssues
		issuesList = ui.RenderIssues(closedIssues)
		ui.UpdateIssueDetails(currentIssues,selectedIndex, showComments, commentsText, issueDetails, issuesList)
	} else {
		currentIssues = openIssues
		issuesList = ui.RenderIssues(openIssues)
		ui.UpdateIssueDetails(currentIssues,selectedIndex, showComments, commentsText, issueDetails, issuesList)
	}

	termui.Render(actionsTabs, issuesList, issueDetails, footer)

	uiEvents := termui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q":// Quit on 'q'
			return
		case "<Down>":
			if selectedIndex < len(currentIssues)-1 {
				selectedIndex++
				issuesList.ScrollDown()
				ui.UpdateIssueDetails(currentIssues, selectedIndex, showComments, commentsText, issueDetails, issuesList)
			}
			if selectedIndex == len(currentIssues)-1 {
				selectedIndex = 0
				issuesList.ScrollTop()
				ui.UpdateIssueDetails(currentIssues, selectedIndex, showComments, commentsText, issueDetails, issuesList)
			}
			commentsText = ""
			showComments = false
		case "<Up>":
			if selectedIndex > 0 {
				selectedIndex--
				issuesList.ScrollUp()
				ui.UpdateIssueDetails(currentIssues, selectedIndex, showComments, commentsText, issueDetails, issuesList)
			}
			if selectedIndex == 0 {
				selectedIndex = len(currentIssues) - 1
				issuesList.ScrollBottom()
				ui.UpdateIssueDetails(currentIssues, selectedIndex, showComments, commentsText, issueDetails, issuesList)
			}
			commentsText = ""
			showComments = false
		case "<Enter>":
			issue := currentIssues[selectedIndex]
			err := ui.OpenBrowser(issue.URL)
			if err != nil {
				log.Printf("Failed to open browser: %v", err)
			}
		case "h":
			showHelper = !showHelper
			if showHelper {
				help = ui.RenderHelp()
			} else {
				help = nil // Clear the help when toggling off
			}
		case "<C-o>":
			showClosedIssues = !showClosedIssues
			if showClosedIssues {
				currentIssues = closedIssues
				issuesList = ui.RenderIssues(closedIssues)
			} else {
				currentIssues = openIssues
				issuesList = ui.RenderIssues(openIssues)
			}
			selectedIndex = 0
			ui.UpdateIssueDetails(currentIssues, selectedIndex, showComments, commentsText, issueDetails, issuesList)
			termui.Render(actionsTabs, issuesList, issueDetails, footer)
		case "<C-p>":
			showPr = !showPr
			if showPr {
				currentIssues = prToReview.Items
				issuesList = ui.RenderIssues(prToReview.Items)
			} else {
				currentIssues = openIssues
				issuesList = ui.RenderIssues(openIssues)
			}
			selectedIndex = 0
			ui.UpdateIssueDetails(currentIssues, selectedIndex, showComments, commentsText, issueDetails, issuesList)
			termui.Render(actionsTabs, issuesList, issueDetails, footer)
		case "<C-c>":
			issue := currentIssues[selectedIndex]
			if !showComments {
				comments, err := github.FetchComments(issue.CommentsURL, ghUserName, ghToken)
				if err != nil {
					log.Printf("Error fetching comments for issue %s: %v", issue.URL, err)
					commentsText = "Error fetching comments."
				} else {
					commentsText = ""
					for _, comment := range comments {
						if (comment.User.Login != "vercel[bot]") {
							commentsText += fmt.Sprintf("Comment by %s at %s:\n%s\n\n",
								comment.User.Login, comment.CreatedAt, comment.Body)
						}
					}
				}
			}
			showComments = !showComments
			ui.UpdateIssueDetails(currentIssues, selectedIndex, showComments, commentsText, issueDetails, issuesList)
		}
		termui.Render(issuesList, issueDetails)
		if showHelper && help != nil {
			termui.Render(help)
		}
	}
}

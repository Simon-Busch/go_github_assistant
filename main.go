package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"

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
	prToReview, err := github.FetchReviewRequests(ghUserName, ghToken)
	if err != nil {
		log.Fatalf("Error fetching issues: %v", err)
	}
	issuesRaw := issuesResponse.Items
	var issues []github.Issue
	var closedIssues []github.Issue
	for i, issue := range issuesResponse.Items {
		if issue.State == "open" {
			issues = append(issues, issuesRaw[i])
		} else {
			closedIssues = append(closedIssues, issuesRaw[i])
		}
	}

	if err := ui.Init(); err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	termWidth, termHeight := ui.TerminalDimensions()

	renderWaitingScreen(len(issues), len(closedIssues), ghUserName, prToReview.TotalCount)
	time.Sleep(2 * time.Second)
	ui.Clear()

	actionsTabs := renderHeader()
	issuesList := renderIssues(issues)
	issueDetails := renderIssueDetails()
	footer := renderFooter()

	ui.Render(actionsTabs, issuesList, issueDetails, footer)

	selectedIndex := 0
	showComments := false
	commentsText := ""
	showHelper := false
	var help *widgets.Paragraph

	updateIssueDetails(issues,selectedIndex, showComments, commentsText, issueDetails, issuesList)

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
				updateIssueDetails(issues, selectedIndex, showComments, commentsText, issueDetails, issuesList)
			}
			if selectedIndex == len(issues)-1 {
				selectedIndex = 0
				issuesList.ScrollTop()
				updateIssueDetails(issues, selectedIndex, showComments, commentsText, issueDetails, issuesList)
			}
			commentsText = ""
		case "<Up>":
			if selectedIndex > 0 {
				selectedIndex--
				issuesList.ScrollUp()
				updateIssueDetails(issues, selectedIndex, showComments, commentsText, issueDetails, issuesList)
			}
			if selectedIndex == 0 {
				selectedIndex = len(issues) - 1
				issuesList.ScrollBottom()
				updateIssueDetails(issues, selectedIndex, showComments, commentsText, issueDetails, issuesList)
			}
			commentsText = ""
		case "<Enter>":
			issue := issues[selectedIndex]
			err := openBrowser(issue.URL)
			if err != nil {
				log.Printf("Failed to open browser: %v", err)
			}
		case "h":
			showHelper = !showHelper
			if showHelper {
				help = widgets.NewParagraph()
				helpBoxWidth := termWidth / 2
				helpBoxHeight := termHeight / 4
				x0 := (termWidth - helpBoxWidth) / 2
				y0 := (termHeight - helpBoxHeight) / 2
				x1 := x0 + helpBoxWidth
				y1 := y0 + helpBoxHeight

				help.SetRect(x0, y0, x1, y1)
				help.Title = "Help"
				help.Text = "List of actions: \n\n'q' to quit \n<Enter> to open the issue in the browser \n<C-c> to toggle comments.\n'h' to open help "
			} else {
				help = nil // Clear the help when toggling off
			}
		case "<C-o>":
			//TODO - Not working
			issues = closedIssues
			issuesList := renderIssues(issues)
			ui.Render(issuesList)
		case "<C-c>":
			issue := issues[selectedIndex]
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
			updateIssueDetails(issues, selectedIndex, showComments, commentsText, issueDetails, issuesList)
		}
		ui.Render(issuesList, issueDetails)
		if showHelper && help != nil {
			ui.Render(help)
		}
	}
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

func updateIssueDetails(issues []github.Issue ,index int, showComments bool, commentsText string, issueDetails *widgets.Paragraph, issuesList *widgets.List) {
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


func renderIssues(issues []github.Issue) *widgets.List {
	termWidth, termHeight := ui.TerminalDimensions()
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
	issuesList.SetRect(0, 3, termWidth/2, termHeight-5)

	return issuesList
}

func renderHeader() *widgets.TabPane {
	termWidth, _ := ui.TerminalDimensions()
	actionsTabs := widgets.NewTabPane("Github Assistant", "Issues", "Pull Requests", "Help")
	actionsTabs.SetRect(0, 0, termWidth, 3)
	return actionsTabs
}

func renderFooter() *widgets.Paragraph {
	termWidth, termHeight := ui.TerminalDimensions()
	footer := widgets.NewParagraph()
	footer.SetRect(0, termHeight-4, termWidth, termHeight)
	footer.Title = "Help"
	footer.Text = "Press 'h' to show help."
	return footer
}

func renderIssueDetails() *widgets.Paragraph {
	termWidth, termHeight := ui.TerminalDimensions()
	issueDetails := widgets.NewParagraph()
	issueDetails.Title = "Issue Details"
	issueDetails.Text = "Select an issue to see details."
	issueDetails.SetRect(termWidth, 3, termWidth/2, termHeight-5)
	return issueDetails
}

func createAsciiFrames(text string) []string {
	frames := []string{}
	for i := 1; i <= len(text); i++ {
		frames = append(frames, text[:i])
	}
	return frames
}

func renderWaitingScreen(openedIssues, closedIssues int, name string, prToReview int) {
	termWidth, termHeight := ui.TerminalDimensions()

	title := "GitHub Assistant"
	open := fmt.Sprintf("Open Issues: %d", openedIssues)
	closed := fmt.Sprintf("Closed Issues: %d", closedIssues)
	pr := fmt.Sprintf("Pr to review: %d", prToReview)
	user := fmt.Sprintf("User: %s", name)
	frames := [][]string{
    createAsciiFrames(title),
    createAsciiFrames(open),
    createAsciiFrames(closed),
    createAsciiFrames(user),
    createAsciiFrames(pr),
	}

	helpBoxWidth := termWidth / 2
	helpBoxHeight := termHeight / 4
	x0 := (termWidth - helpBoxWidth) / 2
	y0 := (termHeight - helpBoxHeight) / 2
	x1 := x0 + helpBoxWidth
	y1 := y0 + helpBoxHeight

	asciiAnimation := widgets.NewParagraph()
	asciiAnimation.SetRect(x0, y0, x1, y1)
	asciiAnimation.TextStyle.Fg = ui.ColorGreen
	ui.Render(asciiAnimation)

	for _, frameSet := range frames {
		for _, frame := range frameSet {
			asciiAnimation.Text = frame
			ui.Render(asciiAnimation)
			time.Sleep(70 * time.Millisecond) // Control animation speed
		}
		time.Sleep(500 * time.Millisecond)
	}

	ui.Clear()
}

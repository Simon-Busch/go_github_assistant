package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/Simon-Busch/go_github_assistant/github"
	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func createAsciiFrames(text string) []string {
	frames := []string{}
	for i := 1; i <= len(text); i++ {
		frames = append(frames, text[:i])
	}
	return frames
}


func RenderWaitingScreen(openedIssues, closedIssues int, name string, prToReview int) {
	termWidth, termHeight := termui.TerminalDimensions()

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
	asciiAnimation.TextStyle.Fg = termui.ColorGreen
	termui.Render(asciiAnimation)

	for _, frameSet := range frames {
		for _, frame := range frameSet {
			asciiAnimation.Text = frame
			termui.Render(asciiAnimation)
			time.Sleep(70 * time.Millisecond) // Control animation speed
		}
		time.Sleep(500 * time.Millisecond)
	}

	termui.Clear()
}

func OpenBrowser(url string) error {
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

func RenderHeader() *widgets.TabPane {
	termWidth, _ := termui.TerminalDimensions()
	actionsTabs := widgets.NewTabPane("Github Assistant", "Issues", "Pull Requests", "Help")
	actionsTabs.SetRect(0, 0, termWidth, 3)
	return actionsTabs
}

func RenderIssues(issues []github.IssueOrPR) *widgets.List {
	termWidth, termHeight := termui.TerminalDimensions()
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
	issuesList.SelectedRowStyle.Fg = termui.ColorYellow
	issuesList.SetRect(0, 3, termWidth/2, termHeight-5)

	return issuesList
}

func RenderHelp() *widgets.Paragraph {
	termWidth, termHeight := termui.TerminalDimensions()
	help := widgets.NewParagraph()
	helpBoxWidth := termWidth / 2
	helpBoxHeight := termHeight / 4
	x0 := (termWidth - helpBoxWidth) / 2
	y0 := (termHeight - helpBoxHeight) / 2
	x1 := x0 + helpBoxWidth
	y1 := y0 + helpBoxHeight

	help.SetRect(x0, y0, x1, y1)
	help.Title = "Help"
	help.Text = "List of actions: \n\n'q' to quit \n<Enter> to open the issue in the browser \n<C-c> to toggle comments.\n<C-o> to toggle between open and closed issues.\n'h' to open help "
	return help;
}

func RenderFooter(openedIssues, closedIssues int, username string, prToReview int) *widgets.Paragraph {
	termWidth, termHeight := termui.TerminalDimensions()
	footer := widgets.NewParagraph()
	footer.SetRect(0, termHeight-4, termWidth, termHeight)
	footer.Title = "Summary"
	footer.Text = fmt.Sprintf("Open Issues: %d, Closed Issues: %d, User: %s, Pr to review: %d", openedIssues, closedIssues, username, prToReview)
	return footer
}

func RenderIssueDetails() *widgets.Paragraph {
	termWidth, termHeight := termui.TerminalDimensions()
	issueDetails := widgets.NewParagraph()
	issueDetails.Title = "Issue Details"
	issueDetails.Text = "Select an issue to see details."
	issueDetails.SetRect(termWidth, 3, termWidth/2, termHeight-5)
	return issueDetails
}

func UpdateIssueDetails(issues []github.IssueOrPR ,index int, showComments bool, commentsText string, issueDetails *widgets.Paragraph, issuesList *widgets.List) {
	if len(issues) == 0 {
		issueDetails.Text = "No issues to show."
		termui.Render(issuesList, issueDetails)
		return
	}
	issue := issues[index]
	issueText := fmt.Sprintf(
		"Title: %s\n\nRepository: %s\nOrganization: %s\n\nState: %s\n\nURL: %s\nCreated At: %s\n\nDescription:\n\n%s",
		issue.Title, issue.Repository, issue.Organization, issue.State, issue.URL, issue.CreatedAt, issue.Body)
	if showComments && commentsText != "" {
		issueText += fmt.Sprintf("\n\nComments:\n%s", commentsText)
	}
	issueDetails.Text = issueText
	termui.Render(issuesList, issueDetails)
}

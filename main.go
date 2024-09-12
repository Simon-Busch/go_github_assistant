package main

import (
	"fmt"
	"log"

	"github.com/Simon-Busch/go_github_assistant/github"
	"github.com/Simon-Busch/go_github_assistant/utils"
	// ui "github.com/gizak/termui/v3"
	// "github.com/gizak/termui/v3/widgets"
)

func main() {
	ghUserName, ghToken, err := utils.InitEnv()
	if err != nil {
		panic("please provide a .env file with GITHUB_USERNAME and GITHUB_TOKEN")
	}
	// if err := ui.Init(); err != nil {
	// 	log.Fatal(err)
	// }
	// defer ui.Close()

	// // Set up a simple UI element (e.g., a paragraph)
	// waitingScreen := widgets.NewParagraph()
	// waitingScreen.Text = "Github assistant\nPress 'q' to quit"

	// termWidth, termHeight := ui.TerminalDimensions()
	// waitingScreen.SetRect(0, 0, termWidth, termHeight)
	// ui.Render(waitingScreen)


	// // Event handler for key presses
	// for e := range ui.PollEvents() {
	// 	switch e.ID {
	// 	case "q", "<C-c>": // "q" or Ctrl+C to quit
	// 		return
	// 	}
	// }
	issues, err := github.FetchIssues(ghUserName, ghToken)
	if err != nil {
		log.Fatalf("Error fetching issues: %v", err)
		
	}

	fmt.Println("Issues assigned: %s",len(issues.Items))

	pr, err := github.FetchReviewRequests(ghUserName, ghToken)
	if err != nil {
		log.Fatalf("Error fetching pull requests: %v", err)
	}
	fmt.Println("Pull requests assigned: %s",len(pr.Items))
}

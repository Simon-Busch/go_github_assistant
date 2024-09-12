package main

import (
	// ui "github.com/gizak/termui/v3"
	// "github.com/gizak/termui/v3/widgets"
	"github.com/Simon-Busch/go_github_assistant/utils"
)

func main() {
	_, _, err := utils.InitEnv()
	if err != nil {
		panic("please provide a .env file with GITHUB_USERNAME and GITHUB_TOKEN")
	}
	// if err := ui.Init(); err != nil {
	// 	log.Fatal(err)
	// }
	// defer ui.Close()

	// // Set up a simple UI element (e.g., a paragraph)
	// p := widgets.NewParagraph()
	// p.Text = "Press 'q' to quit"
	// p.SetRect(0, 0, 50, 3)
	// ui.Render(p)

	// // Event handler for key presses
	// for e := range ui.PollEvents() {
	// 	switch e.ID {
	// 	case "q", "<C-c>": // "q" or Ctrl+C to quit
	// 		return
	// 	}
	// }
}

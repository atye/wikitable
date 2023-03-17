package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/atye/wikitable/internal/model"
	"github.com/atye/wikitable2json/pkg/client"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	userAgent := flag.String("user-agent", "github.com/atye/wikitable", "user agent for making Wikipedia API requests")
	flag.Parse()

	//log = newLogger()

	if _, err := tea.NewProgram(model.NewModel(client.NewTableGetter(*userAgent)), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("error running program:", err)
		os.Exit(1)
	}
}

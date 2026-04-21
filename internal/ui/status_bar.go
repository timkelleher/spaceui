package ui

import (
	"fmt"
	"time"

	"github.com/timkelleher/spaceui/internal/state"
)

func statusBarContent() string {
	currentShip := ""
	if state.HasActiveShip() {
		currentShip = fmt.Sprintf("[orange]%s[-]", state.ActiveShip().Symbol)
	}

	now := fmt.Sprintf("%s", state.FormattedTime(time.Now()))
	agent := state.Agent(false)
	agentSymbol := fmt.Sprintf("[red]%s[-]", agent.Symbol)
	credits := fmt.Sprintf("[green]%d[-]", agent.Credits)
	content := now + "\t" + agentSymbol + "\t" + currentShip + "\t" + credits

	err := state.GlobalError()
	if err != "" {
		content = fmt.Sprintf("[red]%s[-]", err)
	}
	return content
}

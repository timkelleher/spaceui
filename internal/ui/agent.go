package ui

import (
	"fmt"

	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/logger"
)

func AgentAction() {
	app.setCurrentTab("agent")
	app.setDataPaneData("Loading...")
	logger.Info("setting current tab to agents")

	go func() {
		agent, resp := api.GetAgent()
		if resp.Err != nil || resp.Resp.StatusCode() != 200 {
			return
		}

		content := "[yellow]Current Agent Information[-]\n" +
			fmt.Sprintf("[blue]Account ID:[-]\t\t%s\n", agent.Agent.AccountID) +
			fmt.Sprintf("[blue]Symbol:[-]\t\t\t%s\n", agent.Agent.Symbol) +
			fmt.Sprintf("[blue]Ship Count:[-]\t\t%d\n", agent.Agent.ShipCount) +
			fmt.Sprintf("[blue]Headquarters:[-]\t%s\n", agent.Agent.Headquarters) +
			fmt.Sprintf("[blue]Credits:[-]\t\t%d\n", agent.Agent.Credits)
		app.setDataPaneData(content)
	}()
}

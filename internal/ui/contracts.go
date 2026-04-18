package ui

import (
	"fmt"
	"time"

	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/logger"
)

func ContractsAction() {
	app.setCurrentTab("contracts")
	app.setDataPaneData("Loading...")
	logger.Info("setting current tab to contracts")

	contracts, resp := api.GetContracts()
	if resp.Err != nil || resp.Resp.StatusCode() != 200 {
		return
	}

	content := fmt.Sprintf("[yellow]Number of contracts:[-] %d\n", len(contracts.Contracts))
	content += "----------\n"
	for i, contract := range contracts.Contracts {
		content += fmt.Sprintf("[red]Contract %s[-]\n", contract.ID)
		content += fmt.Sprintf("[blue]%s[-] | %s | %s\n", contract.Type, contract.HasAccepted(), contract.HasFulfilled())
		content += fmt.Sprintf("Deadline to Accept: %s\n", contract.DeadlineToAccept.Format(time.RFC1123))
		content += fmt.Sprintf("Expiration: %s\n", contract.Expiration.Format(time.RFC1123))
		content += "----- Delivery Terms -----\n"
		content += fmt.Sprintf("Deadline: %s\n", contract.Terms.Deadline.Format(time.RFC1123))
		for _, deliver := range contract.Terms.Deliver {
			content += fmt.Sprintf("Progress: %d/%d %s\n", deliver.UnitsFulfilled, deliver.UnitsRequired, deliver.DestinationSymbol)
		}
		if i != len(contracts.Contracts)-1 {
			content += "\n"
		}
	}
	app.setDataPaneData(content)
}

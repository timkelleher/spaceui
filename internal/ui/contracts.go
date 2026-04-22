package ui

import (
	"fmt"

	"github.com/timkelleher/spaceui/internal/state"
)

const (
	PAGE_CONTRACTS = "contracts"
)

type ContractsPage struct {
}

func (cp ContractsPage) ID() string {
	return PAGE_CONTRACTS
}

func (cp ContractsPage) Menu() Menu {
	return MainMenu{}
}

func (cp ContractsPage) RequiredData() []string {
	return []string{state.DATA_CONTRACTS}
}

func (cp ContractsPage) Content() string {
	contracts := state.Contracts(false)

	content := ""
	if !state.HasActiveShip() {
		content += "\t[yellow]Please activate a ship in order to negotiate a new contract.[-]\n" +
			"\n"
	}

	content += fmt.Sprintf("[yellow]Number of contracts:[-] %d\n", len(contracts))
	content += "----------\n"

	for i, contract := range contracts {
		content += fmt.Sprintf("[red]ID:[-]\t\t\t %s[-]\n", contract.ID)
		content += fmt.Sprintf("[blue]Type:[-]\t\t %s\n", contract.Type)
		content += fmt.Sprintf("[orange]State:[-]\t\t %s\n", contract.State())
		if !contract.Accepted {
			content += fmt.Sprintf("[green]Deadline to Accept:[-] %s\n", state.FormattedTime(contract.DeadlineToAccept))
		}
		if !contract.Fulfilled {
			content += fmt.Sprintf("[green]Expiration:[-]\t %s\n", state.FormattedTime(contract.Expiration))
		}
		content += fmt.Sprintf("[green]Deadline:[-]\t %s\n", state.FormattedTime(contract.Terms.Deadline))
		content += "----- Delivery Terms -----\n"
		for _, deliver := range contract.Terms.Deliver {
			content += fmt.Sprintf("[orange]%s[-]\t (%d/%d)\n", deliver.DestinationSymbol, deliver.UnitsFulfilled, deliver.UnitsRequired)
		}
		if i != len(contracts)-1 {
			content += "\n"
		}
	}
	return content
}

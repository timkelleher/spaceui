package ui

import (
	"fmt"
)

/*
func (a *App) systemsContent() string {
	systems := state.Systems(false)

	content := fmt.Sprintf("[yellow]Number of systems scanned:[-] %d\n", len(systems))
	if len(systems) == 0 {
		return content
	}

	for _, system := range systems {
		content += fmt.Sprintf("[orange]%s[-] [blue]%s[-] [gold]%d,%d[-]\n", system.Symbol, system.Type, system.X, system.Y)
		waypointCounts := make(map[string]int)
		var waypointNames []string
		for _, waypoint := range system.Waypoints {
			waypointCounts[waypoint.Type]++
		}
		for waypointType, count := range waypointCounts {
			waypointNames = append(waypointNames, fmt.Sprintf("%d [purple]%s[-]", count, waypointType))
		}
		sort.Strings(waypointNames)
		if len(waypointNames) > 0 {
			content += strings.Join(waypointNames, " ") + "\n"
		}
	}
	return content
}
*/

func loadingContent(num int, objName string) string {
	if num == 0 {
		return fmt.Sprintf("Loading %s...\n", objName)
	}

	return fmt.Sprintf("Loading %d %s...\n", num, objName)
}

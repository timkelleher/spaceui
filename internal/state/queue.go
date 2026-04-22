package state

import (
	"time"
)

// /////////////////////////////////////
// Loading & Last Updated
// /////////////////////////////////////
var (
	refreshQueue []string

	loading     map[string]bool
	lastUpdated map[string]time.Time
)

// Newly proposed system:
// If currently loading, return true
// If not currently loading, but data has been update before, return false
// If not currently loading, but data has never been updated before, return true and trigger update
func IsRefreshing(id string) bool {
	loading, ok := loading[id]
	if !ok {
		loading = false
	}
	return loading
}

func AreRefreshing(ids []string) bool {
	for _, id := range ids {
		if IsRefreshing(id) {
			return true
		}
	}
	return false
}

func AnyRefreshing() bool {
	for _, isRefreshing := range loading {
		if isRefreshing {
			return true
		}
	}
	return false
}

func LastUpdated(id string) time.Time {
	lastUpdated, ok := lastUpdated[id]
	if !ok {
		return time.Time{}
	}
	return lastUpdated
}

func Fresh(id string) bool {
	lastUpdated, ok := lastUpdated[id]
	if !ok || lastUpdated.IsZero() {
		return false
	}
	return true
}

func MarkStale(id string) {
	if !IsRefreshing(id) {
		lastUpdated[id] = time.Time{}
	}
}

func Poll() {
	for {
		currentlyLoading := false
		for datatype, isLoading := range loading {
			if isLoading && lastUpdated[datatype].IsZero() {
				currentlyLoading = true
			}
		}

		mu.Lock()
		if len(refreshQueue) > 0 && !currentlyLoading {
			refresh(refreshQueue[0])
			refreshQueue = refreshQueue[1:]
		}
		mu.Unlock()

		time.Sleep(500 * time.Millisecond)
	}
}

func refresh(datatype string) {
	if !IsRefreshing(datatype) {
		switch datatype {
		case DATA_AGENT:
			go Agent(true)
		case DATA_CONTRACTS:
			go Contracts(true)
		case DATA_SHIPS:
			go Ships(true)
		case DATA_WAYPOINTS:
			go Waypoints(true)
		}
	}
}

func Queue(datatype string) {
	mu.Lock()
	defer mu.Unlock()

	if IsRefreshing(datatype) {
		return
	}
	for _, item := range refreshQueue {
		if item == datatype {
			return
		}
	}

	refreshQueue = append(refreshQueue, datatype)
}

package state

import (
	"time"
)

// /////////////////////////////////////
// Active Page
// /////////////////////////////////////
var (
	activePage string
)

func ActivePage() string {
	return activePage
}

func SetActivePage(id string) {
	activePage = id
}

// /////////////////////////////////////
// Utils
// /////////////////////////////////////
var loc *time.Location

func SetLoc(l *time.Location) {
	loc = l
}

func FormattedTime(t time.Time) string {
	if loc == nil {
		return t.Format(time.RFC1123)
	}
	return t.In(loc).Format(time.RFC1123)
}

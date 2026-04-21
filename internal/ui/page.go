package ui

import "github.com/rivo/tview"

type Page interface {
	ID() string
	Menu() Menu
	Content() string
}

type Menu interface {
	ID() string
	Menu() *tview.List
}

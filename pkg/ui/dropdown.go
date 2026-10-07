package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func searchHook(list *tview.List, values []string, pages *tview.Pages, infoPanel *tview.TextView, eventPanel *tview.TextView, root *tview.TreeNode) func(string) {
	return func(s string) {
		list.Clear()
		searchTerm := strings.ToLower(s)

		for _, value := range values {
			if strings.Contains(strings.ToLower(value), searchTerm) {
				list.AddItem(value, "", 0, dropDownSelectedHook(pages, infoPanel, root, eventPanel, value))
			}
		}
	}
}

func createSearchableList(app *tview.Application, pages *tview.Pages, infoPanel *tview.TextView, eventPanel *tview.TextView, root *tview.TreeNode, values []string) *tview.Flex {
	input := tview.NewInputField().SetLabel(" Search: ")
	list := tview.NewList().ShowSecondaryText(false)

	input.SetChangedFunc(searchHook(list, values, pages, infoPanel, eventPanel, root))

	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDown {
			app.SetFocus(list)
			return nil
		} else if event.Key() == tcell.KeyEsc {
			pages.RemovePage("dropdown")
			return nil
		}
		return event
	})

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyUp && list.GetCurrentItem() == 0 {
			app.SetFocus(input)
			return nil
		} else if event.Key() == tcell.KeyEsc {
			pages.RemovePage("dropdown")
			return nil
		}
		return event
	})

	container := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(input, 1, 0, true).
		AddItem(list, 0, 1, false)

	container.SetBorder(true).SetTitle(" Directories ")

	return container
}

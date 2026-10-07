package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Harichandra-Prasath/Tman/pkg/tman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func handleDropdown(workDir string, app *tview.Application, pages *tview.Pages, root *tview.TreeNode, infoPanel *tview.TextView, eventPanel *tview.TextView) error {
	entries, err := os.ReadDir(workDir)
	if err != nil {
		return fmt.Errorf("error walking workDir: %v", err)
	}
	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	dropDown := createSearchableList(workDir, app, pages, infoPanel, eventPanel, root, dirs)

	overlay := tview.NewGrid().
		SetColumns(0, 30, 0).
		SetRows(0, 10, 0).
		AddItem(dropDown, 1, 1, 1, 1, 0, 0, true)

	pages.AddPage("dropdown", overlay, true, true)
	app.SetFocus(dropDown)
	return nil
}

func searchHook(workDir string, list *tview.List, values []string, pages *tview.Pages, infoPanel *tview.TextView, eventPanel *tview.TextView, root *tview.TreeNode) func(string) {
	return func(s string) {
		list.Clear()
		searchTerm := strings.ToLower(s)

		for _, value := range values {
			if strings.Contains(strings.ToLower(value), searchTerm) {
				list.AddItem(value, "", 0, dropDownSelectedHook(workDir, pages, infoPanel, root, eventPanel, value))
			}
		}
	}
}

func createSearchableList(workDir string, app *tview.Application, pages *tview.Pages, infoPanel *tview.TextView, eventPanel *tview.TextView, root *tview.TreeNode, values []string) *tview.Flex {
	input := tview.NewInputField().SetLabel(" Search: ")
	list := tview.NewList().ShowSecondaryText(false)

	input.SetChangedFunc(searchHook(workDir, list, values, pages, infoPanel, eventPanel, root))

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

func dropDownSelectedHook(workDir string, pages *tview.Pages, infoPanel *tview.TextView, root *tview.TreeNode, eventPanel *tview.TextView, selectedText string) func() {
	return func() {
		defer pages.RemovePage("dropdown")
		rootComp := root.GetReference().(*TmanTreeNode).Component.(*tman.Root)

		sessionDir := filepath.Join(workDir, selectedText)
		session, err := tman.CreateSessionComponent(sessionDir, selectedText, rootComp)
		if err != nil {
			infoPanel.SetText(fmt.Sprintf("error creating session: %v", err)).SetTextColor(tcell.ColorRed)
			return
		}
		node := &TmanTreeNode{Parent: root, Component: session}
		addNodes(root, []*TmanTreeNode{node}, tcell.ColorBlue)
		eventPanel.SetText(fmt.Sprintf("New Session Created: %s", session.Name())).SetTextColor(tcell.ColorGreen)
		rootComp.SetChildCount(rootComp.GetChildCount() + 1)
	}
}

package ui

import (
	"fmt"

	"github.com/Harichandra-Prasath/Tman/pkg/tman"
	"github.com/gdamore/tcell/v2"
	tview "github.com/rivo/tview"
)

func StartUI() error {
	app := tview.NewApplication()
	app.SetInputCapture(appKeyHooks(app))
	mainFlex, err := buildUI()
	if err != nil {
		return fmt.Errorf("building ui: %v", err)
	}

	if err := app.SetRoot(mainFlex, true).Run(); err != nil {
		return fmt.Errorf("running ui: %v", err)
	}
	return nil
}

func buildUI() (*tview.Flex, error) {
	// Info Panel
	infoPanel := tview.NewTextView().SetDynamicColors(true)
	infoPanel.SetTitle("Details").SetBorder(true)

	tree := tview.NewTreeView()
	tree.SetInputCapture(treeKeyHooks(tree, infoPanel))
	tree.SetTitle("Tman - tmux Manager").SetBorder(true).SetTitleColor(tcell.ColorWhiteSmoke)
	tree.SetChangedFunc(hoverTreeHook(infoPanel))
	tree.SetSelectedFunc(selectedNodeHook(infoPanel))

	compFlex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(tree, 0, 2, true).AddItem(infoPanel, 0, 1, false)

	rootComp, sessions, err := tman.GetRootAndSessions()
	if err != nil {
		return nil, fmt.Errorf("retrieving root and sessions: %v", err)
	}
	root := tview.NewTreeNode(rootComp.Name()).SetReference(&tman.TmuxTreeNode{Parent: nil, Component: rootComp}).SetColor(tcell.ColorGreen)
	tree.SetRoot(root).SetCurrentNode(root)

	var rootChilds []*tman.TmuxTreeNode
	for _, session := range sessions {
		rootChilds = append(rootChilds, &tman.TmuxTreeNode{Parent: root, Component: session})
	}
	addNodes(root, rootChilds, tcell.ColorBlue)

	centerRowFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(compFlex, 20, 1, true).
		AddItem(nil, 0, 1, false)
	centerColFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(centerRowFlex, 60, 1, true).
		AddItem(nil, 0, 1, false)

	return centerColFlex, nil
}

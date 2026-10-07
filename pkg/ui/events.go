package ui

import (
	"fmt"
	"os"

	"github.com/Harichandra-Prasath/Tman/pkg/tman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func deleteNode(tree *tview.TreeView, node *tview.TreeNode) {
	ref := node.GetReference().(*tman.TmuxTreeNode)
	comp := ref.Component

	parentNode := ref.Parent
	if parentNode == nil {
		// Root Node [Destroy the entire tree]
		tree.SetRoot(nil)
		return
	}

	parentNode.RemoveChild(node)

	parent := comp.GetParent()
	if parent == nil {
		// Should not exectute, but to be safer
		tree.SetRoot(nil)
		return
	}
	currChildren := parent.GetChildCount()
	parent.SetChildCount(currChildren - 1)

	// Get the parent first
	if parent.GetChildCount() == 0 {
		// remove the parent as well
		deleteNode(tree, parentNode)
	}
}

func dropDownSelectedHook(pages *tview.Pages, infoPanel *tview.TextView, root *tview.TreeNode, eventPanel *tview.TextView, selectedText string) func() {
	return func() {
		defer pages.RemovePage("dropdown")
		rootComp := root.GetReference().(*tman.TmuxTreeNode).Component.(*tman.Root)

		session, err := tman.CreateSessionComponent(selectedText, rootComp)
		if err != nil {
			infoPanel.SetText(fmt.Sprintf("error creating session: %v", err)).SetTextColor(tcell.ColorRed)
			return
		}
		node := &tman.TmuxTreeNode{Parent: root, Component: session}
		addNodes(root, []*tman.TmuxTreeNode{node}, tcell.ColorBlue)
		eventPanel.SetText(fmt.Sprintf("New Session Created: %s", session.Name())).SetTextColor(tcell.ColorGreen)
		rootComp.SetChildCount(rootComp.GetChildCount() + 1)
	}
}

func handleDropdown(app *tview.Application, pages *tview.Pages, root *tview.TreeNode, infoPanel *tview.TextView, eventPanel *tview.TextView) error {
	workDir := tman.GlobalTmanConfig.WorkDir
	entries, err := os.ReadDir(workDir)
	if err != nil {
		return fmt.Errorf("error creating session: %v", err)
	}
	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	dropDown := createSearchableList(app, pages, infoPanel, eventPanel, root, dirs)

	overlay := tview.NewGrid().
		SetColumns(0, 30, 0).
		SetRows(0, 10, 0).
		AddItem(dropDown, 1, 1, 1, 1, 0, 0, true)

	pages.AddPage("dropdown", overlay, true, true)
	app.SetFocus(dropDown)
	return nil
}

func refreshTree(tree *tview.TreeView) *tview.TreeNode {
	root := tree.GetRoot()
	tree.SetRoot(nil)
	if root == nil {
		// create a dummy root and attach [Eventually will be replaced]
		root = tview.NewTreeNode("Tmux Server")
		root.SetReference(&tman.TmuxTreeNode{Component: &tman.Root{}})
	}
	root.ClearChildren()
	_selectedWebhook(root)
	return root
}

func globalKeyHooks(app *tview.Application, pages *tview.Pages, tree *tview.TreeView, eventPanel *tview.TextView, infoPanel *tview.TextView) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyEscape {
			app.Stop()
			return nil
		} else if event.Rune() == 'r' {
			root := refreshTree(tree)

			// Possible failure case
			if len(root.GetChildren()) != 0 {
				tree.SetRoot(root)
				eventPanel.SetText("Tree Synced with Tmux Server").SetTextColor(tcell.ColorGreen)
			}
		} else if event.Rune() == 'c' {
			root := tree.GetRoot()
			err := handleDropdown(app, pages, root, infoPanel, eventPanel)
			if err != nil {
				infoPanel.SetText(fmt.Sprintf("error creating session: %v", err)).SetTextColor(tcell.ColorRed)
				return nil
			}
		}
		return event
	}
}

func treeKeyHooks(tree *tview.TreeView, infoPanel *tview.TextView, eventPanel *tview.TextView, app *tview.Application) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		node := tree.GetCurrentNode()
		_ref := node.GetReference()
		if _ref == nil {
			return event
		}
		ref := _ref.(*tman.TmuxTreeNode)
		comp := ref.Component
		switch event.Rune() {
		case 'd':
			msg, err := tman.DeleteTmuxComponent(comp)
			if err != nil {
				infoPanel.SetText(fmt.Sprintf("deleting node: %v", err)).SetTextColor(tcell.ColorRed)
				return event
			}
			eventPanel.SetText(msg).SetTextColor(tcell.ColorGreen)
			// Remove the nodes
			infoPanel.Clear()
			deleteNode(tree, node)
		case 's':
			msg, err := tman.SwitchTmuxComponent(comp)
			if err != nil {
				infoPanel.SetText(fmt.Sprintf("switching to target: %v", err)).SetTextColor(tcell.ColorRed)
				return event
			}
			if msg != "" {
				eventPanel.SetText(msg).SetTextColor(tcell.ColorGreen)
				app.Stop()
			}
		default:
			return event
		}

		return event
	}
}

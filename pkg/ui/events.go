package ui

import (
	"fmt"

	"github.com/Harichandra-Prasath/Tman/pkg/tman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func deleteNode(node *tview.TreeNode) {
	ref := node.GetReference().(*tman.TmuxTreeNode)
	comp := ref.Component

	parentNode := ref.Parent
	if parentNode == nil {
		return
	}
	parentNode.RemoveChild(node)

	parent := comp.GetParent()
	if parent == nil {
		return
	}
	currChildren := parent.GetChildCount()
	parent.SetChildCount(currChildren - 1)

	// Get the parent first
	if parent.GetChildCount() == 0 {
		// remove the parent as well
		deleteNode(parentNode)
	}
}

func appKeyHooks(app *tview.Application, tree *tview.TreeView, eventPanel *tview.TextView) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyEscape {
			app.Stop()
			return nil
		} else if event.Rune() == 'r' {
			root := tree.GetRoot()
			root.SetChildren([]*tview.TreeNode{})
			_selectedWebhook(root)
			tree.SetCurrentNode(root)
			eventPanel.SetText("Tree Synced with Tmux Server").SetTextColor(tcell.ColorGreen)
		}
		return event
	}
}

func treeKeyHooks(tree *tview.TreeView, infoPanel *tview.TextView, eventPanel *tview.TextView) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		node := tree.GetCurrentNode()
		ref := node.GetReference().(*tman.TmuxTreeNode)
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
			deleteNode(node)
		default:
			return event
		}

		return event
	}
}

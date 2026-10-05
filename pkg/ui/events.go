package ui

import (
	"fmt"

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

func appKeyHooks(app *tview.Application, tree *tview.TreeView, eventPanel *tview.TextView) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyEscape {
			app.Stop()
			return nil
		} else if event.Rune() == 'r' {
			root := tree.GetRoot()
			tree.SetRoot(nil)
			if root == nil {
				// create a dummy root and attach [Eventually will be replaced]
				root = tview.NewTreeNode("Tmux Server")
				root.SetReference(&tman.TmuxTreeNode{Component: &tman.Root{}})
			}
			root.ClearChildren()
			_selectedWebhook(root)

			// Possible failure case
			if len(root.GetChildren()) != 0 {
				tree.SetRoot(root)
				eventPanel.SetText("Tree Synced with Tmux Server").SetTextColor(tcell.ColorGreen)
			}
		}
		return event
	}
}

func treeKeyHooks(tree *tview.TreeView, infoPanel *tview.TextView, eventPanel *tview.TextView) func(*tcell.EventKey) *tcell.EventKey {
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
		default:
			return event
		}

		return event
	}
}

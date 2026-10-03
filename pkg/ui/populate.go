package ui

import (
	"fmt"

	"github.com/Harichandra-Prasath/Tman/pkg/tman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createNode(name string, ref any) *tview.TreeNode {
	return tview.NewTreeNode(name).SetReference(ref).SetSelectable(true)
}

func addComponents[T tman.TmuxComponent](target *tview.TreeNode, components []T) {
	for _, component := range components {
		target.AddChild(createNode(component.Name(), component))
	}
}

func hoverTreeHook(infoPanel *tview.TextView) func(*tview.TreeNode) {
	return func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			infoPanel.Clear()
			return
		}

		data := ref.(tman.TmuxComponent)
		infoPanel.SetText(data.Details()).SetTextColor(tcell.ColorGreenYellow)
	}
}

func selectedNodeHook(infoPanel *tview.TextView) func(*tview.TreeNode) {
	return func(node *tview.TreeNode) {
		ref := node.GetReference().(tman.TmuxComponent)
		if ref == nil {
			return
		}
		children := node.GetChildren()

		if len(children) == 0 {
			switch ref := ref.(type) {
			case *tman.Session:
				windows, err := tman.GetWindows(ref)
				if err != nil {
					infoPanel.SetText(fmt.Sprintf("error getting windows: %v", err)).SetTextColor(tcell.ColorRed)
				}
				addComponents(node, windows)
			default:
				return
			}
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	}
}

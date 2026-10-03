package ui

import (
	"fmt"

	"github.com/Harichandra-Prasath/Tman/pkg/tman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createNode(name string, ref any, color tcell.Color) *tview.TreeNode {
	return tview.NewTreeNode(name).SetReference(ref).SetSelectable(true).SetColor(color)
}

func addComponents[T tman.TmuxComponent](target *tview.TreeNode, components []T, color tcell.Color) {
	for _, component := range components {
		target.AddChild(createNode(component.Name(), component, color))
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
				addComponents(node, windows, tcell.ColorAqua)
			case *tman.Window:
				panes, err := tman.GetPanes(ref)
				if err != nil {
					infoPanel.SetText(fmt.Sprintf("error getting panes: %v", err)).SetTextColor(tcell.ColorRed)
				}
				addComponents(node, panes, tcell.ColorPurple)
			default:
				return
			}
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	}
}

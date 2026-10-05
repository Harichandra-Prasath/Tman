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

func addNodes(target *tview.TreeNode, nodes []*tman.TmuxTreeNode, color tcell.Color) {
	for _, node := range nodes {
		target.AddChild(createNode(node.Component.Name(), node, color))
	}
}

func hoverTreeHook(infoPanel *tview.TextView) func(*tview.TreeNode) {
	return func(node *tview.TreeNode) {
		ref := node.GetReference().(*tman.TmuxTreeNode)
		if ref == nil {
			infoPanel.Clear()
			return
		}

		data := ref.Component
		infoPanel.SetText(data.Details()).SetTextColor(tcell.ColorGreenYellow)
	}
}

func selectedNodeHook(infoPanel *tview.TextView) func(*tview.TreeNode) {
	return func(node *tview.TreeNode) {
		ref := node.GetReference().(*tman.TmuxTreeNode)
		if ref == nil {
			return
		}
		children := node.GetChildren()

		if len(children) == 0 {
			var nodes []*tman.TmuxTreeNode
			switch ref := ref.Component.(type) {
			case *tman.Session:
				windows, err := tman.GetWindows(ref)
				if err != nil {
					infoPanel.SetText(fmt.Sprintf("error getting windows: %v", err)).SetTextColor(tcell.ColorRed)
				}

				for _, window := range windows {
					nodes = append(nodes, &tman.TmuxTreeNode{Parent: node, Component: window})
				}

				addNodes(node, nodes, tcell.ColorAqua)
			case *tman.Window:
				panes, err := tman.GetPanes(ref)
				if err != nil {
					infoPanel.SetText(fmt.Sprintf("error getting panes: %v", err)).SetTextColor(tcell.ColorRed)
				}
				for _, pane := range panes {
					nodes = append(nodes, &tman.TmuxTreeNode{Parent: node, Component: pane})
				}
				addNodes(node, nodes, tcell.ColorPurple)
			default:
				return
			}
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	}
}

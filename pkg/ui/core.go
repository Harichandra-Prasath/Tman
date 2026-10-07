package ui

import (
	"github.com/Harichandra-Prasath/Tman/pkg/tman"
	"github.com/rivo/tview"
)

type TmanTreeNode struct {
	Parent    *tview.TreeNode
	Component tman.TmuxComponent
}

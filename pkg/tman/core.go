package tman

import (
	"fmt"

	"github.com/rivo/tview"
)

type TmanConfig struct {
	WorkDir string
}

type TmuxTreeNode struct {
	Parent    *tview.TreeNode
	Component TmuxComponent
}

type TmuxComponent interface {
	Name() string
	Details() string
	GetParent() TmuxComponent
	GetChildCount() int
	SetChildCount(int)
}
type Root struct {
	sessionsCount int
	version       string
	socketPath    string
	user          string
}

func (R *Root) Details() string {
	return fmt.Sprintf("Tmux Version: %s\nNo of Sessions: %d\nServer Socket Path: %s\nServer User: %s\n", R.version, R.sessionsCount, R.socketPath, R.user)
}

func (R *Root) Name() string {
	return "Tmux Server"
}

func (R *Root) GetParent() TmuxComponent {
	return nil
}

func (R *Root) GetChildCount() int {
	return R.sessionsCount
}

func (R *Root) SetChildCount(count int) {
	R.sessionsCount = count
}

type Session struct {
	root         *Root
	sessionName  string
	windowsCount int
	createdAt    string
	directory    string
}

func (S *Session) Details() string {
	return fmt.Sprintf("Session Name: %s\nNo of Windows: %d\nSession Created At: %s\nSession Directory: %s\n", S.sessionName, S.windowsCount, S.createdAt, S.directory)
}

func (S *Session) Name() string {
	return S.sessionName
}

func (S *Session) GetParent() TmuxComponent {
	return S.root
}

func (S *Session) GetChildCount() int {
	return S.windowsCount
}

func (S *Session) SetChildCount(count int) {
	S.windowsCount = count
}

type Window struct {
	parentSession *Session
	windowIndex   string
	panesCount    int
	windowName    string
}

func (W *Window) Details() string {
	return fmt.Sprintf("Window Index: %s\nNo of Panes: %d\nWindow Name (Focused Pane): %s\n", W.windowIndex, W.panesCount, W.windowName)
}

func (W *Window) Name() string {
	return W.windowName
}

func (W *Window) GetParent() TmuxComponent {
	return W.parentSession
}

func (W *Window) GetChildCount() int {
	return W.panesCount
}

func (W *Window) SetChildCount(count int) {
	W.panesCount = count
}

type Pane struct {
	parentWindow       *Window
	paneIndex          string
	paneCurrentCommand string
	paneCurrentPath    string
}

func (P *Pane) Details() string {
	return fmt.Sprintf("Pane Index: %s\nPane Current Command: %s\nPane Current Path: %s\n", P.paneIndex, P.paneCurrentCommand, P.paneCurrentPath)
}

func (P *Pane) Name() string {
	return P.paneCurrentCommand
}

func (P *Pane) GetParent() TmuxComponent {
	return P.parentWindow
}

func (P *Pane) GetChildCount() int {
	// N/A, should not be confused with zero children
	return -1
}

func (P *Pane) SetChildCount(count int) {}

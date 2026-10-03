package tman

import "fmt"

type TmuxComponent interface {
	Name() string
	Details() string
}

type Session struct {
	SessionName  string
	WindowsCount int
	CreatedAt    string
	Directory    string
}

func (S *Session) Details() string {
	return fmt.Sprintf("Session Name: %s\nNo of Windows: %d\nSession Created At: %s\nSession Directory: %s\n", S.SessionName, S.WindowsCount, S.CreatedAt, S.Directory)
}

func (S *Session) Name() string {
	return S.SessionName
}

type Window struct {
	ParentSessionName string
	WindowIndex       string
	PanesCount        int
	WindowName        string
}

func (W *Window) Details() string {
	return fmt.Sprintf("Window Index: %s\nNo of Panes: %d\nWindow Name (Focused Pane): %s\n", W.WindowIndex, W.PanesCount, W.WindowName)
}

func (W *Window) Name() string {
	return W.WindowName
}

type Pane struct {
	PaneIndex          string
	PaneCurrentCommand string
	PaneCurrentPath    string
}

func (P *Pane) Details() string {
	return fmt.Sprintf("Pane Index: %s\nPane Current Command: %s\nPane Current Path: %s\n", P.PaneIndex, P.PaneCurrentCommand, P.PaneCurrentPath)
}

func (P *Pane) Name() string {
	return P.PaneCurrentCommand
}

package tman

import "fmt"

type TmuxComponent interface {
	Name() string
	Details() string
}

type Session struct {
	SessionId    string
	WindowsCount int
	CreatedAt    string
}

type Window struct {
	WindowId    string
	PanesCount  int
	FocusedPane string
}

func (S *Session) Details() string {
	return fmt.Sprintf("Session ID: %s\nNo of Windows: %d\nSession Created At: %s\n", S.SessionId, S.WindowsCount, S.CreatedAt)
}

func (S *Session) Name() string {
	return S.SessionId
}

func (W *Window) Details() string {
	return fmt.Sprintf("Window Id: %s\nNo of Panes: %d\nFocused Pane: %s\n", W.WindowId, W.PanesCount, W.FocusedPane)
}

func (W *Window) Name() string {
	return W.WindowId
}

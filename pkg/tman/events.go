package tman

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func DeleteTmuxComponent(comp TmuxComponent) (string, error) {
	var cmd *exec.Cmd
	var msg string
	switch ref := comp.(type) {
	case *Root:
		msg = "Tmux Server Killed"
		cmd = exec.Command("tmux", "kill-server")
	case *Session:
		target := ref.sessionName
		msg = "Target Session Killed: " + target
		cmd = exec.Command("tmux", "kill-session", "-t", target)
	case *Window:
		sessionName := ref.parentSession.sessionName
		target := sessionName + ":" + ref.windowIndex
		msg = "Target Window Killed: " + target
		cmd = exec.Command("tmux", "kill-window", "-t", target)

	case *Pane:
		SessionName := ref.parentWindow.parentSession.sessionName
		windoIndex := ref.parentWindow.windowIndex
		target := SessionName + ":" + windoIndex + "." + ref.paneIndex
		msg = "Target Pane Killed: " + target
		cmd = exec.Command("tmux", "kill-pane", "-t", target)
	}

	_, err := runStandard(cmd)
	if err != nil {
		return "", fmt.Errorf("executing kill(delete): %v", err)
	}

	return msg, nil
}

func SwitchTmuxComponent(comp TmuxComponent) (string, error) {
	var cmd *exec.Cmd
	var msg string
	switch ref := comp.(type) {
	case *Root:
		return "", nil

	case *Session:
		target := ref.sessionName
		msg = "Switched to Session: " + target
		cmd = exec.Command("tmux", "switch-client", "-t", target)
	case *Window:
		sessionName := ref.parentSession.sessionName
		target := sessionName + ":" + ref.windowIndex
		msg = "Switched to Window: " + target
		cmd = exec.Command("tmux", "switch-client", "-t", target)

	case *Pane:
		SessionName := ref.parentWindow.parentSession.sessionName
		windoIndex := ref.parentWindow.windowIndex
		target := SessionName + ":" + windoIndex + "." + ref.paneIndex
		msg = "Switched to Pane: " + target
		cmd = exec.Command("tmux", "switch-client", "-t", target)
	}

	_, err := runStandard(cmd)
	if err != nil {
		return "", fmt.Errorf("executing switch: %v", err)
	}

	return msg, nil
}

func CreateSessionComponent(name string, rootComp *Root) (*Session, error) {
	sessionDir := filepath.Join(GlobalTmanConfig.WorkDir, name)

	// drop the . for hidden folders
	name = strings.TrimPrefix(name, ".")

	cmd := exec.Command("tmux", "new-session", "-d", "-c", sessionDir, "-n", name, "-s", name)
	_, err := runStandard(cmd)
	if err != nil {
		return nil, fmt.Errorf("executing new session: %v", err)
	}

	t := time.Now()

	session := &Session{
		root:         rootComp,
		sessionName:  name,
		windowsCount: 1,
		createdAt:    t.Format("Jan 02 Mon 15:04"),
		directory:    sessionDir,
	}

	return session, nil
}

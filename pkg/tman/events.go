package tman

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func targetFor(comp TmuxComponent) string {
	var target string
	switch ref := comp.(type) {
	case *Session:
		target = ref.sessionName
	case *Window:
		sessionName := ref.parentSession.sessionName
		target = sessionName + ":" + ref.windowIndex
	case *Pane:
		sessionName := ref.parentWindow.parentSession.sessionName
		windowIndex := ref.parentWindow.windowIndex
		target = sessionName + ":" + windowIndex + "." + ref.paneIndex
	}
	return target
}

func DeleteTmuxComponent(comp TmuxComponent) (string, error) {
	var cmd *exec.Cmd
	var msg string

	target := targetFor(comp)
	switch comp.(type) {
	case *Root:
		msg = "Tmux Server Killed"
		cmd = exec.Command("tmux", "kill-server")
	case *Session:
		msg = "Target Session Killed: " + target
		cmd = exec.Command("tmux", "kill-session", "-t", target)
	case *Window:
		msg = "Target Window Killed: " + target
		cmd = exec.Command("tmux", "kill-window", "-t", target)

	case *Pane:
		msg = "Target Pane Killed: " + target
		cmd = exec.Command("tmux", "kill-pane", "-t", target)
	default:
		return "", fmt.Errorf("not a valid node for delete")
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
	target := targetFor(comp)
	if target == "" {
		return "", nil
	}
	cmd = exec.Command("tmux", "switch-client", "-t", target)
	msg = "Switched to target: " + target
	_, err := runStandard(cmd)
	if err != nil {
		return "", fmt.Errorf("executing switch: %v", err)
	}

	return msg, nil
}

func CreateSessionComponent(sessionDir string, name string, rootComp *Root) (*Session, error) {
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

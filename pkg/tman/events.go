package tman

import (
	"fmt"
	"os/exec"
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

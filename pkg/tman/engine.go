package tman

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runStandard(cmd *exec.Cmd) (*bufio.Scanner, error) {
	var out strings.Builder
	var er strings.Builder

	cmd.Stdout = &out
	cmd.Stderr = &er

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error in run: %v", er.String())
	}

	// Expand the output for safe exits
	scanner := bufio.NewScanner(strings.NewReader(out.String()))

	return scanner, nil
}

func IsReachable() bool {
	cmd := exec.Command("tmux", "-V")
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func parseSessionLine(lineText string) (*Session, error) {
	details := strings.Split(lineText, ":")
	if len(details) != 4 {
		return nil, fmt.Errorf("parse error: unexpected line text")
	}

	sessionID := details[0]

	windowCounts, _ := strconv.Atoi(details[1])

	creationUnix := details[2]
	unixInt, _ := strconv.ParseInt(creationUnix, 10, 64)
	t := time.Unix(unixInt, 0)
	createdAt := t.Format("Jan 02 Mon 15:04")

	sessionDirectory := details[3]

	return &Session{SessionName: sessionID, WindowsCount: windowCounts, CreatedAt: createdAt, Directory: sessionDirectory}, nil
}

func GetSessions() ([]*Session, error) {
	var sessions []*Session

	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}:#{session_windows}:#{session_created}:#{session_path}")
	scanner, err := runStandard(cmd)
	if err != nil {
		return nil, fmt.Errorf("executing list sessions: %v", err)
	}

	for scanner.Scan() {
		lineText := scanner.Text()
		session, err := parseSessionLine(lineText)
		if err != nil {
			return nil, fmt.Errorf("line text parsing: %v", err)
		}
		sessions = append(sessions, session)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %v", err)
	}

	return sessions, nil
}

func parseWindowLine(lineText string, sessionName string) (*Window, error) {
	details := strings.Split(lineText, ":")
	if len(details) != 3 {
		return nil, fmt.Errorf("parse error: unexpected line text")
	}

	windowIndex := details[0]

	paneCount, _ := strconv.Atoi(details[1])

	windowName := details[2]

	return &Window{ParentSessionName: sessionName, WindowIndex: windowIndex, PanesCount: paneCount, WindowName: windowName}, nil
}

func GetWindows(session *Session) ([]*Window, error) {
	var windows []*Window

	cmd := exec.Command("tmux", "list-windows", "-t", session.SessionName, "-F", "#{window_index}:#{window_panes}:#{window_name}")
	scanner, err := runStandard(cmd)
	if err != nil {
		return nil, fmt.Errorf("executing list windows: %v", err)
	}

	for scanner.Scan() {
		lineText := scanner.Text()
		window, err := parseWindowLine(lineText, session.SessionName)
		if err != nil {
			return nil, fmt.Errorf("line text parsing: %v", err)
		}
		windows = append(windows, window)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %v", err)
	}

	return windows, nil
}

func parsePaneLine(lineText string) (*Pane, error) {
	details := strings.Split(lineText, ":")
	if len(details) != 3 {
		return nil, fmt.Errorf("parse error: unexpected line text")
	}

	paneIndex := details[0]

	paneCurrentCommand := details[1]

	paneCurrentPath := details[2]

	return &Pane{PaneIndex: paneIndex, PaneCurrentCommand: paneCurrentCommand, PaneCurrentPath: paneCurrentPath}, nil
}

func GetPanes(window *Window) ([]*Pane, error) {
	var panes []*Pane

	cmd := exec.Command("tmux", "list-panes", "-t", window.ParentSessionName+":"+window.WindowIndex, "-F", "#{pane_index}:#{pane_current_command}:#{pane_current_path}")
	scanner, err := runStandard(cmd)
	if err != nil {
		return nil, fmt.Errorf("executing list panes: %v", err)
	}

	for scanner.Scan() {
		lineText := scanner.Text()
		pane, err := parsePaneLine(lineText)
		if err != nil {
			return nil, fmt.Errorf("line text parsing: %v", err)
		}
		panes = append(panes, pane)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %v", err)
	}

	return panes, nil
}

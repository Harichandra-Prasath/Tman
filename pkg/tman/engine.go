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
		if detail := strings.TrimSpace(er.String()); detail != "" {
			return nil, fmt.Errorf("%s (%v)", detail, err)
		}
		return nil, fmt.Errorf("%v", err)
	}

	// Expand the output for safe exits
	scanner := bufio.NewScanner(strings.NewReader(out.String()))

	return scanner, nil
}

func parseSessionLine(lineText string, rootBuild bool) (*Session, *Root, error) {
	details := strings.Split(lineText, ":")
	if rootBuild && len(details) != 8 {
		return nil, nil, fmt.Errorf("parse error: unexpected line text")
	}

	if len(details) < 4 {
		return nil, nil, fmt.Errorf("parse error: unexpected line text: %v", lineText)
	}

	sessionID := details[0]

	windowCounts, _ := strconv.Atoi(details[1])

	creationUnix := details[2]
	unixInt, _ := strconv.ParseInt(creationUnix, 10, 64)
	t := time.Unix(unixInt, 0)
	createdAt := t.Format("Jan 02 Mon 15:04")

	sessionDirectory := details[3]

	root := &Root{}
	if rootBuild {
		version := details[4]
		socketPath := details[5]
		user := details[6]
		sessionsCount, _ := strconv.Atoi(details[7])

		root.version = version
		root.socketPath = socketPath
		root.user = user
		root.sessionsCount = sessionsCount
	}

	return &Session{sessionName: sessionID, windowsCount: windowCounts, createdAt: createdAt, directory: sessionDirectory}, root, nil
}

func GetRootAndSessions() (*Root, []*Session, error) {
	var sessions []*Session
	var root *Root

	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}:#{session_windows}:#{session_created}:#{session_path}:#{version}:#{socket_path}:#{user}:#{server_sessions}")
	scanner, err := runStandard(cmd)
	if err != nil {
		return nil, nil, fmt.Errorf("executing list sessions: %v", err)
	}

	rootBuild := true
	for scanner.Scan() {
		lineText := scanner.Text()
		session, tempRoot, err := parseSessionLine(lineText, rootBuild)
		if rootBuild {
			root = tempRoot
			rootBuild = false
		}
		if err != nil {
			return nil, nil, fmt.Errorf("line text parsing: %v", err)
		}
		session.root = root
		sessions = append(sessions, session)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("scanner error: %v", err)
	}

	return root, sessions, nil
}

func parseWindowLine(lineText string) (*Window, error) {
	details := strings.Split(lineText, ":")
	if len(details) != 3 {
		return nil, fmt.Errorf("parse error: unexpected line text")
	}

	windowIndex := details[0]

	paneCount, _ := strconv.Atoi(details[1])

	windowName := details[2]

	return &Window{windowIndex: windowIndex, panesCount: paneCount, windowName: windowName}, nil
}

func GetWindows(session *Session) ([]*Window, error) {
	var windows []*Window

	cmd := exec.Command("tmux", "list-windows", "-t", session.sessionName, "-F", "#{window_index}:#{window_panes}:#{window_name}")
	scanner, err := runStandard(cmd)
	if err != nil {
		return nil, fmt.Errorf("executing list windows: %v", err)
	}

	for scanner.Scan() {
		lineText := scanner.Text()
		window, err := parseWindowLine(lineText)
		if err != nil {
			return nil, fmt.Errorf("line text parsing: %v", err)
		}
		window.parentSession = session

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

	return &Pane{paneIndex: paneIndex, paneCurrentCommand: paneCurrentCommand, paneCurrentPath: paneCurrentPath}, nil
}

func GetPanes(window *Window) ([]*Pane, error) {
	var panes []*Pane

	cmd := exec.Command("tmux", "list-panes", "-t", window.parentSession.sessionName+":"+window.windowIndex, "-F", "#{pane_index}:#{pane_current_command}:#{pane_current_path}")
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
		pane.parentWindow = window
		panes = append(panes, pane)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %v", err)
	}

	return panes, nil
}

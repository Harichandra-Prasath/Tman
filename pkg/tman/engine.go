package tman

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
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
	colonIdx := strings.IndexByte(lineText, ':')
	if colonIdx == -1 {
		return nil, fmt.Errorf("parse error: unexpected line text")
	}
	sessionID := lineText[:colonIdx]
	startCount := colonIdx + 2
	endCount := startCount + strings.IndexByte(lineText[startCount:], ' ')
	windows := lineText[startCount:endCount]

	windowCounts, _ := strconv.Atoi(windows)

	dateStr := "(created "
	startDate := strings.Index(lineText, dateStr)
	if startDate == -1 {
		return nil, fmt.Errorf("parse error: unexpected line text")
	}
	startDate += len(dateStr)
	endDate := startDate + strings.IndexByte(lineText[startDate:], ')')
	createdAt := lineText[startDate:endDate]

	return &Session{SessionId: sessionID, WindowsCount: windowCounts, CreatedAt: createdAt}, nil
}

func GetSessions() ([]*Session, error) {
	var sessions []*Session

	cmd := exec.Command("tmux", "list-sessions")
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

func parseWindowLine(lineText string) (*Window, error) {
	colonIdx := strings.IndexByte(lineText, ':')
	if colonIdx == -1 {
		return nil, fmt.Errorf("parse error: unexpected line text")
	}

	windowId := lineText[:colonIdx]
	startCount := strings.IndexByte(lineText, '(')
	endStr := " panes"
	endCount := strings.Index(lineText[startCount:], endStr)
	if startCount == -1 || endCount == -1 {
		return nil, fmt.Errorf("parse error: unexpected line text")
	}

	panesCounts := lineText[startCount+1 : startCount+endCount]
	paneCount, _ := strconv.Atoi(panesCounts)

	return &Window{WindowId: windowId, PanesCount: paneCount}, nil
}

func GetWindows(session *Session) ([]*Window, error) {
	var windows []*Window

	cmd := exec.Command("tmux", "list-windows", "-t", session.SessionId)
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
		windows = append(windows, window)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %v", err)
	}

	return windows, nil
}

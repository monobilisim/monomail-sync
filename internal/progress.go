package internal

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type TaskProgress struct {
	TaskID        int    `json:"task_id"`
	Percent       int    `json:"percent"`
	CurrentFolder string `json:"current_folder"`
	FoldersDone   int    `json:"folders_done"`
	FoldersTotal  int    `json:"folders_total"`
	MsgsDone      int    `json:"msgs_done"`
	MsgsTotal     int    `json:"msgs_total"`
	Speed         string `json:"speed"`
	BytesCopied   string `json:"bytes_copied"`
	EtaSeconds    int    `json:"eta_seconds"`
}

var (
	progressStore   = make(map[int]*TaskProgress)
	progressStoreMu sync.RWMutex
)

var (
	reFoldersFound = regexp.MustCompile(`Host1: found (\d+) folders`)
	reTotalMsgs    = regexp.MustCompile(`Host1 Nb messages:\s+(\d+) messages`)
	reFolderLine   = regexp.MustCompile(`^Folder\s+(\d+)/(\d+)\s+\[(.+?)\]`)
	reMsgsLeft     = regexp.MustCompile(`(\d+)/(\d+) msgs left`)
	reSpeed        = regexp.MustCompile(`([\d.]+) msgs/s`)
	reCopied       = regexp.MustCompile(`([\d.]+)\s+(KiB|MiB|GiB) copied`)
	reEtaSecs      = regexp.MustCompile(`\s+(\d+) s\s+\d+/\d+ msgs left`)
)

func SetTaskProgress(taskID int, p *TaskProgress) {
	progressStoreMu.Lock()
	defer progressStoreMu.Unlock()
	progressStore[taskID] = p
}

func GetTaskProgress(taskID int) *TaskProgress {
	progressStoreMu.RLock()
	defer progressStoreMu.RUnlock()
	if p, ok := progressStore[taskID]; ok {
		cp := *p
		return &cp
	}
	return nil
}

func DeleteTaskProgress(taskID int) {
	progressStoreMu.Lock()
	defer progressStoreMu.Unlock()
	delete(progressStore, taskID)
}

func ParseImapsyncOutput(taskID int, r io.Reader, passthrough io.Writer) {
	p := &TaskProgress{TaskID: taskID}
	SetTaskProgress(taskID, p)

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 256*1024), 256*1024)

	for scanner.Scan() {
		line := scanner.Text()

		if passthrough != nil {
			passthrough.Write([]byte(line + "\n"))
		}

		if m := reFoldersFound.FindStringSubmatch(line); m != nil {
			p.FoldersTotal, _ = strconv.Atoi(m[1])
			SetTaskProgress(taskID, p)
			continue
		}

		if m := reTotalMsgs.FindStringSubmatch(line); m != nil {
			p.MsgsTotal, _ = strconv.Atoi(m[1])
			SetTaskProgress(taskID, p)
			continue
		}

		if m := reFolderLine.FindStringSubmatch(line); m != nil {
			done, _ := strconv.Atoi(m[1])
			total, _ := strconv.Atoi(m[2])
			p.FoldersDone = done - 1
			p.FoldersTotal = total
			p.CurrentFolder = m[3]
			SetTaskProgress(taskID, p)
			continue
		}

		if !strings.Contains(line, "msgs left") {
			continue
		}

		if m := reMsgsLeft.FindStringSubmatch(line); m != nil {
			left, _ := strconv.Atoi(m[1])
			total, _ := strconv.Atoi(m[2])
			p.MsgsTotal = total
			p.MsgsDone = total - left
			if total > 0 {
				pct := p.MsgsDone * 100 / total
				if pct > 99 {
					pct = 99
				}
				p.Percent = pct
			}
		}

		if m := reSpeed.FindStringSubmatch(line); m != nil {
			p.Speed = m[1] + " msgs/s"
		}

		if m := reCopied.FindStringSubmatch(line); m != nil {
			p.BytesCopied = m[1] + " " + m[2]
		}

		if m := reEtaSecs.FindStringSubmatch(line); m != nil {
			p.EtaSeconds, _ = strconv.Atoi(m[1])
		}

		SetTaskProgress(taskID, p)
	}

	p.Percent = 100
	SetTaskProgress(taskID, p)
}

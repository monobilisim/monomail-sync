package internal

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseImapsyncOutput_BasicProgress(t *testing.T) {
	input := strings.Join([]string{
		"Host1: found 12 folders",
		"Host1 Nb messages:  436 messages",
		"Folder  3/12 [INBOX]",
		"  220/436 msgs left  11.39 msgs/s  45.2 MiB copied  19 s  220/436 msgs left",
	}, "\n")

	var buf bytes.Buffer
	ParseImapsyncOutput(9999, strings.NewReader(input), &buf)

	p := GetTaskProgress(9999)
	if p == nil {
		t.Fatal("expected progress to be set")
	}

	if p.FoldersTotal != 12 {
		t.Errorf("FoldersTotal = %d, want 12", p.FoldersTotal)
	}
	if p.FoldersDone != 2 {
		t.Errorf("FoldersDone = %d, want 2", p.FoldersDone)
	}
	if p.CurrentFolder != "INBOX" {
		t.Errorf("CurrentFolder = %q, want INBOX", p.CurrentFolder)
	}
	if p.MsgsTotal != 436 {
		t.Errorf("MsgsTotal = %d, want 436", p.MsgsTotal)
	}
	if p.MsgsDone != 216 {
		t.Errorf("MsgsDone = %d, want 216", p.MsgsDone)
	}
	if p.Speed != "11.39 msgs/s" {
		t.Errorf("Speed = %q, want '11.39 msgs/s'", p.Speed)
	}
	if p.BytesCopied != "45.2 MiB" {
		t.Errorf("BytesCopied = %q, want '45.2 MiB'", p.BytesCopied)
	}
	if p.EtaSeconds != 19 {
		t.Errorf("EtaSeconds = %d, want 19", p.EtaSeconds)
	}
	if p.Percent != 100 {
		t.Errorf("Percent after EOF = %d, want 100", p.Percent)
	}

	if !strings.Contains(buf.String(), "Host1: found 12 folders") {
		t.Error("passthrough buffer missing expected line")
	}

	DeleteTaskProgress(9999)
}

func TestParseImapsyncOutput_EmptyInput(t *testing.T) {
	ParseImapsyncOutput(8888, strings.NewReader(""), nil)

	p := GetTaskProgress(8888)
	if p == nil {
		t.Fatal("expected progress entry even with empty input")
	}
	if p.Percent != 100 {
		t.Errorf("Percent = %d, want 100 after EOF", p.Percent)
	}

	DeleteTaskProgress(8888)
}

func TestParseImapsyncOutput_PercentCappedAt99(t *testing.T) {
	input := "  1/100 msgs left  5.00 msgs/s  10 s  1/100 msgs left"
	ParseImapsyncOutput(7777, strings.NewReader(input), nil)

	p := GetTaskProgress(7777)
	if p == nil {
		t.Fatal("expected progress")
	}
	if p.Percent != 100 {
		t.Errorf("Percent = %d, want 100 (after EOF)", p.Percent)
	}

	DeleteTaskProgress(7777)
}

func TestProgressStore_SetGetDelete(t *testing.T) {
	p := &TaskProgress{TaskID: 1234, Percent: 42, CurrentFolder: "Sent"}
	SetTaskProgress(1234, p)

	got := GetTaskProgress(1234)
	if got == nil {
		t.Fatal("expected non-nil")
	}
	if got.Percent != 42 {
		t.Errorf("Percent = %d, want 42", got.Percent)
	}
	if got.CurrentFolder != "Sent" {
		t.Errorf("CurrentFolder = %q, want Sent", got.CurrentFolder)
	}

	got.Percent = 99
	original := GetTaskProgress(1234)
	if original.Percent != 42 {
		t.Error("store returned a reference instead of copy")
	}

	DeleteTaskProgress(1234)
	if GetTaskProgress(1234) != nil {
		t.Error("expected nil after delete")
	}
}

func TestProgressStore_NonExistent(t *testing.T) {
	if GetTaskProgress(999999) != nil {
		t.Error("expected nil for non-existent task")
	}
}

func TestParseImapsyncOutput_MultipleFolder(t *testing.T) {
	input := strings.Join([]string{
		"Host1: found 5 folders",
		"Folder  1/5 [INBOX]",
		"  50/100 msgs left  8.00 msgs/s  6 s  50/100 msgs left",
		"Folder  2/5 [Sent]",
		"  10/50 msgs left  12.50 msgs/s  1 s  10/50 msgs left",
		"Folder  3/5 [Drafts]",
	}, "\n")

	ParseImapsyncOutput(6666, strings.NewReader(input), nil)

	p := GetTaskProgress(6666)
	if p == nil {
		t.Fatal("expected progress")
	}
	if p.CurrentFolder != "Drafts" {
		t.Errorf("CurrentFolder = %q, want Drafts", p.CurrentFolder)
	}
	if p.FoldersDone != 2 {
		t.Errorf("FoldersDone = %d, want 2", p.FoldersDone)
	}
	if p.FoldersTotal != 5 {
		t.Errorf("FoldersTotal = %d, want 5", p.FoldersTotal)
	}

	DeleteTaskProgress(6666)
}

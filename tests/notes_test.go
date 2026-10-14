package tests

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type CheckItem struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

type Note struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	Content    string      `json:"content"`
	Category   string      `json:"category"`
	Tags       []string    `json:"tags"`
	Color      string      `json:"color"`
	Priority   string      `json:"priority"`
	IsPinned   bool        `json:"isPinned"`
	IsFavorite bool        `json:"isFavorite"`
	IsArchived bool        `json:"isArchived"`
	IsTrashed  bool        `json:"isTrashed"`
	Checklist  []CheckItem `json:"checklist"`
	CreatedAt  int64       `json:"createdAt"`
	UpdatedAt  int64       `json:"updatedAt"`
}

func TestNotesCRUD(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "test_notes.json")
	defer os.Remove(tmpFile)

	notes := []Note{
		{
			ID:       "note-test-1",
			Title:    "টেস্ট নোট ১",
			Content:  "এটি একটি টেস্ট নোট",
			Category: "work",
			Tags:     []string{"টেস্ট", "নীলাং"},
			Color:    "indigo",
			IsPinned: true,
			Checklist: []CheckItem{
				{ID: "c1", Text: "টাস্ক ১", Completed: true},
				{ID: "c2", Text: "টাস্ক ২", Completed: false},
			},
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
	}

	data, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	readData, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	var loaded []Note
	if err := json.Unmarshal(readData, &loaded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if len(loaded) != 1 {
		t.Fatalf("Expected 1 note, got %d", len(loaded))
	}
	if loaded[0].Title != "টেস্ট নোট ১" {
		t.Errorf("Expected title 'টেস্ট নোট ১', got '%s'", loaded[0].Title)
	}
	if len(loaded[0].Checklist) != 2 {
		t.Errorf("Expected 2 checklist items, got %d", len(loaded[0].Checklist))
	}
}

func TestQEMUBootImage(t *testing.T) {
	imgPath := filepath.Join("..", "qemu", "notes-os.img")
	info, err := os.Stat(imgPath)
	if err != nil {
		t.Fatalf("QEMU boot image not found at %s: %v", imgPath, err)
	}
	if info.Size() < 512 {
		t.Errorf("Expected boot image >= 512 bytes, got %d", info.Size())
	}

	// Verify boot sector magic 0x55, 0xAA
	data, err := os.ReadFile(imgPath)
	if err != nil {
		t.Fatalf("Failed to read boot image: %v", err)
	}
	if len(data) < 512 {
		t.Fatalf("Boot image too small: %d", len(data))
	}
	if data[510] != 0x55 || data[511] != 0xAA {
		t.Errorf("Invalid boot sector signature: got 0x%02X 0x%02X, want 0x55 0xAA", data[510], data[511])
	}
}

func TestNilangSourceFiles(t *testing.T) {
	files := []string{
		filepath.Join("..", "src", "main.nil"),
		filepath.Join("..", "src", "notes_app.nil"),
		filepath.Join("..", "src", "domain", "types.nil"),
		filepath.Join("..", "src", "domain", "checklist.nil"),
		filepath.Join("..", "src", "domain", "markdown.nil"),
		filepath.Join("..", "src", "domain", "theme.nil"),
		filepath.Join("..", "src", "domain", "storage.nil"),
		filepath.Join("..", "src", "domain", "notes.nil"),
		filepath.Join("..", "src", "ui", "Header.nil"),
		filepath.Join("..", "src", "ui", "Sidebar.nil"),
		filepath.Join("..", "src", "ui", "NoteList.nil"),
		filepath.Join("..", "src", "ui", "NoteEditor.nil"),
		filepath.Join("..", "src", "ui", "StatsModal.nil"),
	}

	for _, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("Required Nilang source file missing: %s", f)
		}
	}
}

func TestNilangExecution(t *testing.T) {
	nilBin := filepath.Join("..", "..", "nilLang", "bin", "nil.exe")
	if _, err := os.Stat(nilBin); os.IsNotExist(err) {
		nilBin = "nil"
	}

	cmd := exec.Command(nilBin, "run", filepath.Join("..", "src", "main.nil"))
	cmd.Stdin = strings.NewReader("q\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("nil run failed: %v, output: %s", err, string(out))
	}

	outStr := string(out)
	if !strings.Contains(outStr, "NILANG NOTES") && !strings.Contains(outStr, "নোটবুক") {
		t.Errorf("Expected banner in output, got: %s", outStr)
	}
}

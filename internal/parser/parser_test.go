package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile_Valid(t *testing.T) {
	// Create a temporary file with valid participants
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "participants.txt")

	content := "Abe\nBeck\nCharlie\nFrank\n"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	names, err := ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	expected := []string{"Abe", "Beck", "Charlie", "Frank"}
	if len(names) != len(expected) {
		t.Fatalf("ParseFile() returned %d names, want %d", len(names), len(expected))
	}

	for i, name := range names {
		if name != expected[i] {
			t.Errorf("ParseFile() name[%d] = %s, want %s", i, name, expected[i])
		}
	}
}

func TestParseFile_WithWhitespace(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "participants.txt")

	content := "  Abe  \n\nBeck\n  Charlie\nFrank  \n\n"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	names, err := ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	expected := []string{"Abe", "Beck", "Charlie", "Frank"}
	if len(names) != len(expected) {
		t.Fatalf("ParseFile() returned %d names, want %d", len(names), len(expected))
	}

	for i, name := range names {
		if name != expected[i] {
			t.Errorf("ParseFile() name[%d] = %s, want %s", i, name, expected[i])
		}
	}
}

func TestParseFile_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "participants.txt")

	if err := os.WriteFile(tmpFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := ParseFile(tmpFile)
	if err == nil {
		t.Error("ParseFile() expected error for empty file, got nil")
	}
}

func TestParseFile_OnlyWhitespace(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "participants.txt")

	if err := os.WriteFile(tmpFile, []byte("  \n\n  \n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := ParseFile(tmpFile)
	if err == nil {
		t.Error("ParseFile() expected error for file with only whitespace, got nil")
	}
}

func TestParseFile_SingleParticipant(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "participants.txt")

	if err := os.WriteFile(tmpFile, []byte("Abe\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := ParseFile(tmpFile)
	if err == nil {
		t.Error("ParseFile() expected error for single participant, got nil")
	}
}

func TestParseFile_DuplicateNames(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "participants.txt")

	content := "Abe\nBeck\nAbe\n"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := ParseFile(tmpFile)
	if err == nil {
		t.Error("ParseFile() expected error for duplicate names, got nil")
	}
}

func TestParseFile_FileNotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/to/file.txt")
	if err == nil {
		t.Error("ParseFile() expected error for nonexistent file, got nil")
	}
}

func TestParseFile_TwoParticipants(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "participants.txt")

	content := "Abe\nBeck\n"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	names, err := ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile() error = %v, want nil (2 participants is valid)", err)
	}

	if len(names) != 2 {
		t.Errorf("ParseFile() returned %d names, want 2", len(names))
	}
}

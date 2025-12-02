package output

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFormatText(t *testing.T) {
	assignments := map[string]string{
		"Abe":     "Frank",
		"Beck":    "Abe",
		"Charlie": "Beck",
		"Frank":   "Charlie",
	}
	names := []string{"Abe", "Beck", "Charlie", "Frank"}

	result := FormatText(assignments, names)

	// Check that result contains all expected assignments
	expectedLines := []string{
		"Abe -> Frank",
		"Beck -> Abe",
		"Charlie -> Beck",
		"Frank -> Charlie",
	}

	for _, expected := range expectedLines {
		if !strings.Contains(result, expected) {
			t.Errorf("FormatText() missing expected line: %s\nGot:\n%s", expected, result)
		}
	}

	// Check that output has correct number of lines
	lines := strings.Split(result, "\n")
	if len(lines) != len(expectedLines) {
		t.Errorf("FormatText() returned %d lines, want %d", len(lines), len(expectedLines))
	}
}

func TestFormatText_Sorted(t *testing.T) {
	assignments := map[string]string{
		"Zoe":     "Alice",
		"Alice":   "Bob",
		"Bob":     "Zoe",
	}
	names := []string{"Zoe", "Alice", "Bob"}

	result := FormatText(assignments, names)

	// Output should be sorted alphabetically by giver
	expected := "Alice -> Bob\nBob -> Zoe\nZoe -> Alice"
	if result != expected {
		t.Errorf("FormatText() = %q, want %q", result, expected)
	}
}

func TestFormatJSON(t *testing.T) {
	assignments := map[string]string{
		"Abe":     "Frank",
		"Beck":    "Abe",
		"Charlie": "Beck",
		"Frank":   "Charlie",
	}

	result, err := FormatJSON(assignments)
	if err != nil {
		t.Fatalf("FormatJSON() error = %v", err)
	}

	// Parse the JSON to verify it's valid
	var output AssignmentsOutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("FormatJSON() produced invalid JSON: %v", err)
	}

	// Check that all assignments are present
	if len(output.Assignments) != len(assignments) {
		t.Errorf("FormatJSON() returned %d assignments, want %d", len(output.Assignments), len(assignments))
	}

	// Verify each assignment
	foundAssignments := make(map[string]string)
	for _, a := range output.Assignments {
		foundAssignments[a.Giver] = a.Receiver
	}

	for giver, receiver := range assignments {
		foundReceiver, ok := foundAssignments[giver]
		if !ok {
			t.Errorf("FormatJSON() missing giver: %s", giver)
			continue
		}
		if foundReceiver != receiver {
			t.Errorf("FormatJSON() giver %s -> %s, want %s", giver, foundReceiver, receiver)
		}
	}
}

func TestFormatJSON_Sorted(t *testing.T) {
	assignments := map[string]string{
		"Zoe":   "Alice",
		"Alice": "Bob",
		"Bob":   "Zoe",
	}

	result, err := FormatJSON(assignments)
	if err != nil {
		t.Fatalf("FormatJSON() error = %v", err)
	}

	var output AssignmentsOutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("FormatJSON() produced invalid JSON: %v", err)
	}

	// Verify assignments are sorted by giver name
	if len(output.Assignments) != 3 {
		t.Fatalf("Expected 3 assignments, got %d", len(output.Assignments))
	}

	if output.Assignments[0].Giver != "Alice" {
		t.Errorf("First assignment giver = %s, want Alice", output.Assignments[0].Giver)
	}
	if output.Assignments[1].Giver != "Bob" {
		t.Errorf("Second assignment giver = %s, want Bob", output.Assignments[1].Giver)
	}
	if output.Assignments[2].Giver != "Zoe" {
		t.Errorf("Third assignment giver = %s, want Zoe", output.Assignments[2].Giver)
	}
}

func TestFormatJSON_ValidStructure(t *testing.T) {
	assignments := map[string]string{
		"Abe":  "Beck",
		"Beck": "Abe",
	}

	result, err := FormatJSON(assignments)
	if err != nil {
		t.Fatalf("FormatJSON() error = %v", err)
	}

	// Verify the JSON has the expected structure
	if !strings.Contains(result, `"assignments"`) {
		t.Error("FormatJSON() missing 'assignments' key")
	}
	if !strings.Contains(result, `"giver"`) {
		t.Error("FormatJSON() missing 'giver' key")
	}
	if !strings.Contains(result, `"receiver"`) {
		t.Error("FormatJSON() missing 'receiver' key")
	}
}

func TestFormatText_EmptyAssignments(t *testing.T) {
	assignments := map[string]string{}
	names := []string{}

	result := FormatText(assignments, names)

	if result != "" {
		t.Errorf("FormatText() with empty assignments = %q, want empty string", result)
	}
}

func TestFormatJSON_EmptyAssignments(t *testing.T) {
	assignments := map[string]string{}

	result, err := FormatJSON(assignments)
	if err != nil {
		t.Fatalf("FormatJSON() error = %v", err)
	}

	var output AssignmentsOutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("FormatJSON() produced invalid JSON: %v", err)
	}

	if len(output.Assignments) != 0 {
		t.Errorf("FormatJSON() with empty assignments returned %d items, want 0", len(output.Assignments))
	}
}

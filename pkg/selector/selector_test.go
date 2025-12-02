package selector

import (
	"testing"
)

func TestGenerateAssignments_Deterministic(t *testing.T) {
	names := []string{"Abe", "Beck", "Charlie", "Frank"}
	seed := make([]byte, SeedBytes)
	for i := range seed {
		seed[i] = byte(i)
	}

	// Generate assignments twice with the same seed
	assignments1, err := GenerateAssignments(names, seed)
	if err != nil {
		t.Fatalf("GenerateAssignments() error = %v", err)
	}

	assignments2, err := GenerateAssignments(names, seed)
	if err != nil {
		t.Fatalf("GenerateAssignments() second call error = %v", err)
	}

	// Should produce identical results
	if len(assignments1) != len(assignments2) {
		t.Fatalf("Assignments have different lengths: %d vs %d", len(assignments1), len(assignments2))
	}

	for giver, receiver1 := range assignments1 {
		receiver2, ok := assignments2[giver]
		if !ok {
			t.Errorf("Giver %s missing in second assignment", giver)
			continue
		}
		if receiver1 != receiver2 {
			t.Errorf("Different assignments for %s: %s vs %s", giver, receiver1, receiver2)
		}
	}
}

func TestGenerateAssignments_NoSelfAssignment(t *testing.T) {
	names := []string{"Abe", "Beck", "Charlie", "Frank"}
	seed := make([]byte, SeedBytes)

	// Test with multiple different seeds
	for seedVal := 0; seedVal < 100; seedVal++ {
		for i := range seed {
			seed[i] = byte(seedVal + i)
		}

		assignments, err := GenerateAssignments(names, seed)
		if err != nil {
			t.Fatalf("GenerateAssignments() error = %v", err)
		}

		// Check no self-assignments
		for giver, receiver := range assignments {
			if giver == receiver {
				t.Errorf("Self-assignment found: %s -> %s", giver, receiver)
			}
		}
	}
}

func TestGenerateAssignments_AllParticipate(t *testing.T) {
	names := []string{"Abe", "Beck", "Charlie", "Frank"}
	seed := make([]byte, SeedBytes)

	assignments, err := GenerateAssignments(names, seed)
	if err != nil {
		t.Fatalf("GenerateAssignments() error = %v", err)
	}

	// Check all names appear as givers
	if len(assignments) != len(names) {
		t.Errorf("Expected %d assignments, got %d", len(names), len(assignments))
	}

	for _, name := range names {
		if _, ok := assignments[name]; !ok {
			t.Errorf("Name %s does not appear as a giver", name)
		}
	}

	// Check all names appear as receivers exactly once
	receivers := make(map[string]int)
	for _, receiver := range assignments {
		receivers[receiver]++
	}

	for _, name := range names {
		count, ok := receivers[name]
		if !ok {
			t.Errorf("Name %s does not appear as a receiver", name)
		} else if count != 1 {
			t.Errorf("Name %s appears as receiver %d times, want 1", name, count)
		}
	}
}

func TestGenerateAssignments_SingleCycle(t *testing.T) {
	names := []string{"Abe", "Beck", "Charlie", "Frank"}
	seed := make([]byte, SeedBytes)

	assignments, err := GenerateAssignments(names, seed)
	if err != nil {
		t.Fatalf("GenerateAssignments() error = %v", err)
	}

	// Follow the chain and verify it forms a single cycle
	visited := make(map[string]bool)
	current := names[0]
	cycleLength := 0

	for {
		if visited[current] {
			break
		}
		visited[current] = true
		current = assignments[current]
		cycleLength++
	}

	// Should visit all participants in one cycle
	if cycleLength != len(names) {
		t.Errorf("Cycle length = %d, want %d (single cycle)", cycleLength, len(names))
	}

	// Should return to the starting point
	if current != names[0] {
		t.Errorf("Cycle does not return to start: ended at %s, started at %s", current, names[0])
	}
}

func TestGenerateAssignments_TwoParticipants(t *testing.T) {
	names := []string{"Abe", "Beck"}
	seed := make([]byte, SeedBytes)

	assignments, err := GenerateAssignments(names, seed)
	if err != nil {
		t.Fatalf("GenerateAssignments() error = %v", err)
	}

	// With 2 participants, each must give to the other
	if assignments["Abe"] != "Beck" {
		t.Errorf("Expected Abe -> Beck, got Abe -> %s", assignments["Abe"])
	}
	if assignments["Beck"] != "Abe" {
		t.Errorf("Expected Beck -> Abe, got Beck -> %s", assignments["Beck"])
	}
}

func TestGenerateAssignments_LargeGroup(t *testing.T) {
	// Create 100 participants
	names := make([]string, 100)
	for i := 0; i < 100; i++ {
		names[i] = string(rune('A' + (i % 26))) + string(rune('0' + (i / 26)))
	}

	seed := make([]byte, SeedBytes)
	assignments, err := GenerateAssignments(names, seed)
	if err != nil {
		t.Fatalf("GenerateAssignments() error = %v", err)
	}

	// Verify basic properties
	if len(assignments) != 100 {
		t.Errorf("Expected 100 assignments, got %d", len(assignments))
	}

	// Check no self-assignments
	for giver, receiver := range assignments {
		if giver == receiver {
			t.Errorf("Self-assignment found: %s -> %s", giver, receiver)
		}
	}

	// Verify single cycle
	visited := make(map[string]bool)
	current := names[0]
	cycleLength := 0

	for {
		if visited[current] {
			break
		}
		visited[current] = true
		current = assignments[current]
		cycleLength++
	}

	if cycleLength != 100 {
		t.Errorf("Cycle length = %d, want 100", cycleLength)
	}
}

func TestGenerateAssignments_ZeroParticipants(t *testing.T) {
	names := []string{}
	seed := make([]byte, SeedBytes)

	_, err := GenerateAssignments(names, seed)
	if err == nil {
		t.Error("GenerateAssignments() expected error for 0 participants, got nil")
	}
}

func TestGenerateAssignments_OneParticipant(t *testing.T) {
	names := []string{"Abe"}
	seed := make([]byte, SeedBytes)

	_, err := GenerateAssignments(names, seed)
	if err == nil {
		t.Error("GenerateAssignments() expected error for 1 participant, got nil")
	}
}

func TestGenerateAssignments_DifferentSeeds(t *testing.T) {
	names := []string{"Abe", "Beck", "Charlie", "Frank"}

	seed1 := make([]byte, SeedBytes)
	seed2 := make([]byte, SeedBytes)
	seed2[0] = 1

	assignments1, err := GenerateAssignments(names, seed1)
	if err != nil {
		t.Fatalf("GenerateAssignments() error = %v", err)
	}

	assignments2, err := GenerateAssignments(names, seed2)
	if err != nil {
		t.Fatalf("GenerateAssignments() error = %v", err)
	}

	// Different seeds should (very likely) produce different assignments
	same := true
	for giver, receiver1 := range assignments1 {
		receiver2 := assignments2[giver]
		if receiver1 != receiver2 {
			same = false
			break
		}
	}

	if same {
		t.Error("Different seeds produced identical assignments (extremely unlikely)")
	}
}

func TestGenerateAssignments_ThreeParticipants(t *testing.T) {
	names := []string{"Abe", "Beck", "Charlie"}
	seed := make([]byte, SeedBytes)

	assignments, err := GenerateAssignments(names, seed)
	if err != nil {
		t.Fatalf("GenerateAssignments() error = %v", err)
	}

	// Verify it's a valid cycle of 3
	if len(assignments) != 3 {
		t.Errorf("Expected 3 assignments, got %d", len(assignments))
	}

	// No self-assignments
	for giver, receiver := range assignments {
		if giver == receiver {
			t.Errorf("Self-assignment found: %s -> %s", giver, receiver)
		}
	}

	// Forms a single cycle
	visited := make(map[string]bool)
	current := names[0]
	cycleLength := 0

	for {
		if visited[current] {
			break
		}
		visited[current] = true
		current = assignments[current]
		cycleLength++
	}

	if cycleLength != 3 {
		t.Errorf("Cycle length = %d, want 3", cycleLength)
	}
}

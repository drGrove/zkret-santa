package selector

import (
	"fmt"
	"math/rand"
)

// GenerateAssignments creates a Secret Santa assignment where each person gives to exactly one other person
// and receives from exactly one other person, forming a single cycle with no self-assignments.
// Uses Sattolo's algorithm to guarantee a single cycle.
func GenerateAssignments(names []string, seed []byte) (map[string]string, error) {
	if len(names) < 2 {
		return nil, fmt.Errorf("need at least 2 participants, got %d", len(names))
	}

	// Create a seeded random number generator
	rng := rand.New(rand.NewSource(SeedToInt64(seed)))

	// Copy the names slice to avoid modifying the input
	cycle := make([]string, len(names))
	copy(cycle, names)

	// Apply Sattolo's algorithm to create a single random cycle
	// This guarantees no fixed points (no self-assignments)
	for i := len(cycle) - 1; i > 0; i-- {
		j := rng.Intn(i) // Note: Intn(i) not Intn(i+1) - this is Sattolo's key difference
		cycle[i], cycle[j] = cycle[j], cycle[i]
	}

	// Build the assignment map: giver -> receiver
	// Each person at index i gives to the person at index (i+1) % len
	assignments := make(map[string]string)
	for i := 0; i < len(names); i++ {
		giver := names[i]
		receiver := cycle[i]
		assignments[giver] = receiver
	}

	return assignments, nil
}

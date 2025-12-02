package output

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Assignment represents a single giver-receiver pair
type Assignment struct {
	Giver    string `json:"giver"`
	Receiver string `json:"receiver"`
}

// AssignmentsOutput represents the JSON output structure
type AssignmentsOutput struct {
	Assignments []Assignment `json:"assignments"`
}

// FormatText formats assignments as human-readable text
// Output format: "Giver -> Receiver" with one assignment per line
// Assignments are sorted by giver name for consistent output
func FormatText(assignments map[string]string, names []string) string {
	// Sort names to ensure consistent output order
	sortedNames := make([]string, len(names))
	copy(sortedNames, names)
	sort.Strings(sortedNames)

	var lines []string
	for _, giver := range sortedNames {
		receiver := assignments[giver]
		lines = append(lines, fmt.Sprintf("%s -> %s", giver, receiver))
	}

	return strings.Join(lines, "\n")
}

// FormatJSON formats assignments as JSON
// Returns a JSON string with an "assignments" array containing giver-receiver pairs
func FormatJSON(assignments map[string]string) (string, error) {
	// Convert map to sorted slice for consistent JSON output
	var assignmentList []Assignment

	// Get all givers and sort them
	givers := make([]string, 0, len(assignments))
	for giver := range assignments {
		givers = append(givers, giver)
	}
	sort.Strings(givers)

	// Build sorted assignment list
	for _, giver := range givers {
		assignmentList = append(assignmentList, Assignment{
			Giver:    giver,
			Receiver: assignments[giver],
		})
	}

	output := AssignmentsOutput{
		Assignments: assignmentList,
	}

	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(jsonBytes), nil
}

package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ParseFile reads a file and returns a list of participant names.
// Each line should contain one name. Whitespace is trimmed.
// Returns an error if the file is empty, has only one participant, or contains duplicates.
func ParseFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}
	defer file.Close()

	var names []string
	seen := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for duplicates
		if seen[line] {
			return nil, fmt.Errorf("duplicate participant: %s", line)
		}

		seen[line] = true
		names = append(names, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Validate minimum participants
	if len(names) == 0 {
		return nil, fmt.Errorf("no participants found")
	}

	if len(names) == 1 {
		return nil, fmt.Errorf("need at least 2 participants")
	}

	return names, nil
}

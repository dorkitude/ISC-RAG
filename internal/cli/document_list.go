package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readDocumentList(path string) ([]string, error) {
	if path == "" {
		return nil, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	seen := map[string]bool{}
	var out []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if seen[line] {
			return nil, fmt.Errorf("duplicate document ID %q", line)
		}
		seen[line] = true
		out = append(out, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("document list %q has no IDs", path)
	}
	return out, nil
}

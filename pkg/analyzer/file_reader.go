package analyzer

import (
	"bufio"
	"os"
	"strings"
)

// ReadFile lit un fichier et retourne son contenu et ses lignes
func ReadFile(filename string) (string, []string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	var lines []string
	var content strings.Builder
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		content.WriteString(line)
		content.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", nil, err
	}

	return content.String(), lines, nil
}
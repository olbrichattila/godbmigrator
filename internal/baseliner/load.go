package baseliner

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (b *baselilner) Load(migrationFilePath string) error {
	filename := migrationFilePath + "/baseline.sql"

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("file opening error %s Error:%v", filename, err)

	}
	defer file.Close()

	var statementBuilder strings.Builder
	scanner := bufio.NewScanner(file)
	isDelimiterSeparation := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "--") || line == "" {
			continue
		}

		if line == openingDelimiter {
			isDelimiterSeparation = true
			continue
		}

		if line != closingDelimiter {
			statementBuilder.WriteString(line + "\n")
		}

		if b.detectStatementEnd(line, isDelimiterSeparation) {
			query := statementBuilder.String()
			statementBuilder.Reset()

			_, err := b.db.Exec(query)
			if err != nil {
				return fmt.Errorf("SQL Execution Error: %v query: %s", err, query)
			}
		}
	}

	return scanner.Err()
}

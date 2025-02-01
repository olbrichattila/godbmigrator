package migrator

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func newBaseliner(specificBaseliner SpecificBaseliner) Baseliner {
	return &baselilner{
		specificBaseliner: specificBaseliner,
	}
}

type baselilner struct {
	specificBaseliner SpecificBaseliner
}

func (b *baselilner) Save(migrationFilePath string) error {
	filename := migrationFilePath + "/baseline.sql"
	
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
		
	}
	defer file.Close()

	err = b.specificBaseliner.GetSchemaData(func(schemaDef string, useDelimiter bool) error {
		if useDelimiter {
			_, err := file.WriteString("DELIMITER ;;\n")
			if err != nil {
				return fmt.Errorf("cannot save baseline %v", err)
			}	
		}

		_, err := file.WriteString(schemaDef + ";\n\n")
		if err != nil {
			return fmt.Errorf("cannot save baseline %v", err)
		}

		if useDelimiter {
			_, err := file.WriteString("DELIMITER ;\n")
			if err != nil {
				return fmt.Errorf("cannot save baseline %v", err)
			}	
		}

		return nil
	});

	if err != nil {
		return err;
	}

	return nil;

}

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

		if line == "DELIMITER ;;" {
			isDelimiterSeparation = true
			continue
		}

		if line != "DELIMITER ;" {
			statementBuilder.WriteString(line + "\n")	
		}
		
		if b.detectStatementEnd(line, isDelimiterSeparation) {
			query := statementBuilder.String()
			statementBuilder.Reset()

			
			_, err := b.specificBaseliner.GetDb().Exec(query)
			if err != nil {
				return fmt.Errorf("SQL Execution Error: %v query: %s", err, query)
			}
		}
	}

	return scanner.Err()
}

func (*baselilner) detectStatementEnd(line string, isDelimiterSeparation bool) bool {
	if isDelimiterSeparation {
		return line == "DELIMITER ;"
	}

	return strings.HasSuffix(line, ";")
}

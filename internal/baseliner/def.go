package baseliner

import (
	"fmt"

	"github.com/olbrichattila/godbmigrator/internal/dbtypemanager"
)

// getEngineSpecificInstructions, if you implement a new database, please add it here with the corresponding dev-<my-database>.go
func (b *baselilner) getEngineSpecificInstructions() (*baselineInstruction, error) {
	dbType, err := dbtypemanager.GetDiverType(b.db)
	if err != nil {
		return nil, fmt.Errorf("getEngineSpecificInstructions: error %v", err)
	}

	var instructions *baselineInstruction

	switch dbType {
	case dbtypemanager.DbTypeSqlite:
		instructions = b.getSQLiteInstruction()
	case dbtypemanager.DbTypeMySQL:
		instructions = b.getMySQLInstruction()
	case dbtypemanager.DbTypePostgres:
		instructions = b.getPostgreSQLInstruction()
	case dbtypemanager.DbTypeFirebird:
		instructions = b.getFirebirdSQLInstruction()
	default:
		return nil, fmt.Errorf("the driver used %s does not match any known driver by the application", dbType)
	}

	if len(instructions.execute) == 0 {
		return nil, fmt.Errorf("the baseline feature not implemented for %s database type", dbType)
	}

	return instructions, nil
}

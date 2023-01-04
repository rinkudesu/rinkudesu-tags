package Migrations

import (
	"embed"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var (
	//go:embed sql
	migrations embed.FS
)

func getMigrationDefinitions(start int, end int) *[]string {
	definitions := make([]string, end-start+1)
	for i := start; i <= end; i++ {
		file, err := migrations.ReadFile(fmt.Sprintf("sql/%d.sql", i))
		if err != nil {
			log.Panicf("Failed to open migration file for migration %d: %s", i, err.Error())
		}
		definitions[start-i] = string(file)
	}
	return &definitions
}

func getLatestAvailableMigration() (latest int) {
	dir, err := migrations.ReadDir("sql")
	if err != nil {
		log.Panic("Failed to get latest migration")
	}
	latest = -1
	for _, entry := range dir {
		if found, err := getMigrationVersionFromFileName(entry.Name()); err == nil && found > latest {
			latest = found
		}
	}
	return
}

func getMigrationVersionFromFileName(fileName string) (int, error) {
	fileName = strings.TrimPrefix(fileName, "sql/")
	fileName = strings.TrimSuffix(fileName, ".sql")
	return strconv.Atoi(fileName)
}

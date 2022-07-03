package Migrations

import (
	"log"
	"rinkudesu-tags/Data"
)

const currentVersion = 0

type Executor struct {
	connection Data.DbConnection
	migrations []func(executor Executor)
}

func NewExecutor(connection Data.DbConnection) Executor {
	var executor = Executor{
		connection: connection,
	}
	executor.initialiseMigrationFunctions()
	return executor
}

func (e *Executor) Migrate() {
	latest := e.getLatestMigration()
	if latest == currentVersion {
		return
	}
	if latest < 0 {
		latest = 0
	}

	for i := latest; i <= currentVersion; i++ {
		e.migrations[i](*e)
	}
}

func (e *Executor) getLatestMigration() int {
	row, err := e.connection.QueryRow("SELECT id FROM migrations ORDER BY id DESC LIMIT 1;")
	if err != nil {
		return -1
	}

	var latest int
	err = row.Scan(&latest)
	if err != nil {
		return -1
	}
	return latest
}

func (e *Executor) initialiseMigrationFunctions() {
	e.migrations = []func(executor Executor){
		initialMigration,
	}
}

func initialMigration(executor Executor) {
	_, err := executor.connection.Query("CREATE TABLE migrations (id integer PRIMARY KEY);")
	if err != nil {
		log.Panicln("Unable to create migrations table")
	}

	_, err = executor.connection.Query("INSERT INTO migrations VALUES (0);")
	if err != nil {
		log.Panicln("Unable to insert migration record")
	}
}

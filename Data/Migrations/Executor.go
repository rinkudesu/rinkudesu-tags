package Migrations

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"rinkudesu-tags/Data"
)

type Executor struct {
	connection Data.DbConnector
	migrations []func(executor Executor) error
}

func NewExecutor(connection Data.DbConnector) *Executor {
	var executor = Executor{
		connection: connection,
	}
	return &executor
}

func (e *Executor) Migrate() {
	latestAvailable := getLatestAvailableMigration()
	latestApplied, _ := e.getLatestAppliedMigration()
	if latestAvailable == latestApplied {
		return
	}

	toApply := getMigrationDefinitions(latestApplied+1, latestAvailable)
	for i, migrationSqlBody := range *toApply {
		migrationIndex := latestApplied + 1 + i
		migrationSql := fmt.Sprintf("begin transaction;\n"+
			"%s\n"+
			"insert into migrations (id) values (%d);\n"+
			"commit;",
			migrationSqlBody,
			migrationIndex)
		_, err := e.connection.Exec(migrationSql)
		if err != nil {
			log.Panicf("Failed to migrate database: %s", err.Error())
		}
	}
}

func (e *Executor) getLatestAppliedMigration() (int, error) {
	row, err := e.connection.QueryRow("SELECT id FROM migrations ORDER BY id DESC LIMIT 1;")
	if err != nil {
		log.Info("Migrations table not found, assuming no migrations ever performed")
		return -1, err
	}

	var latest int
	err = row.Scan(&latest)
	if err != nil {
		log.Warning("Failed to read latest performed migration, assuming no migrations ever performed")
		return -1, err
	}
	log.Infof("Last applied migration: %d", latest)
	return latest, nil
}

func (e *Executor) IsMigrated() (bool, error) {
	current, err := e.getLatestAppliedMigration()
	if err != nil {
		return false, err
	}
	return current == getLatestAvailableMigration(), nil
}

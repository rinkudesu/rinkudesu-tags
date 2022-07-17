package Migrations

import (
	log "github.com/sirupsen/logrus"
	"rinkudesu-tags/Data"
)

const currentVersion = 0

type Executor struct {
	connection Data.DbConnector
	migrations []func(executor Executor) error
}

func NewExecutor(connection Data.DbConnector) Executor {
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
		log.Infof("Running migration %d", i)
		err := e.migrations[i](*e)
		if err != nil {
			log.Panicf("Failed to apply migration %d because: %s", i, err.Error())
		}
	}
}

func (e *Executor) getLatestMigration() int {
	row, err := e.connection.QueryRow("SELECT id FROM migrations ORDER BY id DESC LIMIT 1;")
	if err != nil {
		log.Info("Migrations table not found, assuming no migrations ever performed")
		return -1
	}

	var latest int
	err = row.Scan(&latest)
	if err != nil {
		log.Warning("Failed to read latest performed migration, assuming no migrations ever performed")
		return -1
	}
	log.Infof("Last applied migration: %d", latest)
	return latest
}

func (e *Executor) initialiseMigrationFunctions() {
	e.migrations = []func(executor Executor) error{
		initialMigration,
	}
}

func initialMigration(executor Executor) error {
	_, err := executor.connection.Exec("BEGIN TRANSACTION;" +
		"" +
		"CREATE TABLE migrations (id integer PRIMARY KEY);" +
		"" +
		"CREATE TABLE tags (" +
		"id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid()," +
		"name CHARACTER VARYING(50) NOT NULL," +
		"user_id UUID NOT NULL);" +
		"" +
		"CREATE TABLE links (" +
		"id UUID PRIMARY KEY NOT NULL);" +
		"" +
		"CREATE TABLE link_tags (" +
		"id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid()," +
		"link_id UUID NOT NULL REFERENCES links(id) ON DELETE CASCADE," +
		"tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE);" +
		"" +
		"CREATE UNIQUE INDEX idx_tags_name_user on tags(name, user_id);" +
		"" +
		"INSERT INTO migrations VALUES (0);\n" +
		"" +
		"COMMIT;")
	return err
}

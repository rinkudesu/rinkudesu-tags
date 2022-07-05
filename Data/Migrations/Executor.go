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
		"link_id UUID NOT NULL," +
		"tag_id UUID NOT NULL," +
		"CONSTRAINT fk_link FOREIGN KEY(link_id) REFERENCES links(id)," +
		"CONSTRAINT fk_tag FOREIGN KEY(tag_id) REFERENCES tags(id));" +
		"" +
		"CREATE UNIQUE INDEX idx_tags_name_user on tags(name, user_id);" +
		"" +
		"INSERT INTO migrations VALUES (0);" +
		"" +
		"COMMIT;")
	if err != nil {
		log.Panicln("Unable to apply migration")
	}
}

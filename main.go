package main

import (
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Data/Migrations"
)

func main() {
	connection := Data.DbConnection{}
	_ = connection.Initialise("postgres://postgres:postgres@localhost:5432/postgres")
	defer connection.Close()
	migrator := Migrations.NewExecutor(connection)
	migrator.Migrate()
}

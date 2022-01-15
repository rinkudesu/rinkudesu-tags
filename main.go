package main

import "rinkudesu-tags/Data"

func main() {
	connection := Data.DbConnection{}
	_ = connection.Initialise("postgres://postgres:postgres@localhost:5432/postgres")
	defer connection.Close()

	result, err := connection.QueryRow("select * from test")
	if err != nil {
		panic(err.Error())
	}
	var id string
	err = result.Scan(&id)
	if err != nil {
		panic(err.Error())
	}

	println(id)
}

package data

type DbConnector interface {
	Initialise(connectionString string) error
	QueryRow(sql string, args ...interface{}) (Row, error)
	QueryRows(sql string, args ...interface{}) (Rows, error)
	Query(sql string) (Rows, error)
	Exec(sql string, args ...interface{}) (ExecResult, error)
	Close()
}

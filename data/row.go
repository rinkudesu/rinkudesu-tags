package data

type Row interface {
	Scan(...interface{}) error
}

type Rows interface {
	Scan(...interface{}) error
	Next() bool
	Close()
}

type ExecResult interface {
	RowsAffected() int64
}

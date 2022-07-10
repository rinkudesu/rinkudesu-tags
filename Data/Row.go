package Data

type Row interface {
	Scan(...interface{}) error
}

type Rows interface {
	Scan(...interface{}) error
	Next() bool
	RawValues() [][]byte
}

type ExecResult interface {
	RowsAffected() int64
}

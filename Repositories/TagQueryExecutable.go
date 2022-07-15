package Repositories

import (
	"github.com/gofrs/uuid"
	"rinkudesu-tags/Data"
	"rinkudesu-tags/Models"
)

type TagQueryExecutable interface {
	GetAll() (Data.Rows, error)
	GetSingleById(id uuid.UUID) (Data.Row, error)
	Insert(tag *Models.Tag) (Data.Row, error)
	Update(tag *Models.Tag) (Data.ExecResult, error)
	Delete(id uuid.UUID) (Data.ExecResult, error)
	ScanIntoTag(row Data.Row, id uuid.UUID) (*Models.Tag, error)
}

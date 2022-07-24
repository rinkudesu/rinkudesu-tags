package Services

import "rinkudesu-tags/Data"

type GlobalState struct {
	DbConnection Data.DbConnector
}

func NewGlobalState(dbConnection Data.DbConnector) *GlobalState {
	return &GlobalState{DbConnection: dbConnection}
}

package services

import "rinkudesu-tags/data"

type GlobalState struct {
	DbConnection data.DbConnector
}

func NewGlobalState(dbConnection data.DbConnector) *GlobalState {
	return &GlobalState{DbConnection: dbConnection}
}

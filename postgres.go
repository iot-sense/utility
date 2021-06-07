package utility

import (
	"database/sql"

	_ "github.com/lib/pq"
)

//PgClient struct
type PgClient struct {
	IsConnected bool
	DB          *sql.DB
}

var G_PG_CLIENT *PgClient

//NewPgClient func
func NewPgClient() *PgClient {
	pg := PgClient{false, nil}
	Logger.Info("<< START CONNECT POSTGRESDB >>")
	pgURI := G_CONFIGER.GetString("postgres.uri")
	Logger.Info("postgres.uri :", pgURI)
	// open connection
	db, err := sql.Open("postgres", pgURI)
	if err != nil {
		Logger.Error(err)
	}
	// check connection
	err = db.Ping()
	if err != nil {
		Logger.Error(err)
		return nil
	}
	pg.IsConnected = true
	pg.DB = db
	Logger.Info("<< CONNECT POSTGRESDB SUCCESSFULLY >>")
	return &pg
}

//PgQuery func: select
func PgQuery(db *sql.DB, query string, args ...interface{}) *sql.Rows {
	rows, err := db.Query(query, args...)
	if err != nil {
		Logger.Error(err)
	}
	return rows
}

//PgExec func: insert, update, delete
func PgExec(db *sql.DB, query string, args ...interface{}) sql.Result {
	result, err := db.Exec(query, args...)
	if err != nil {
		Logger.Error(err)
	}
	return result
}

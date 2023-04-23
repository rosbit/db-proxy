package db

import (
	"db-proxy/utils"
	"context"
)

func (txId Tx) Prepare(query string) (stmtId Stmt, err error) {
	tx, e := txId.getTx()
	if e != nil {
		err = e
		return
	}
	stmt, e := tx.Conn.PrepareContext(context.Background(), query)
	if e != nil {
		err = e
		return
	}
	sId := utils.NewObjId(stmt)
	stmtId = Stmt(sId)
	return
}

func (txId Tx) Ping() (err error) {
	tx, e := txId.getTx()
	if e != nil {
		err = e
		return
	}
	return tx.Conn.PingContext(context.Background())
}


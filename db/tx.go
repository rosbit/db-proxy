package db

import (
	"db-proxy/utils"
	"database/sql"
	"fmt"
)

type Tx uint32
type TxData struct {
	Conn *sql.Conn
	Tx   *sql.Tx
}

func (txId Tx) Rollback() (err error) {
	tx, e := txId.getTx()
	if e != nil {
		err = e
		return
	}
	utils.FreeObjId(uint32(txId))
	defer tx.Conn.Close()
	return tx.Tx.Rollback()
}

func (txId Tx) Commit() (err error) {
	tx, e := txId.getTx()
	if e != nil {
		err = e
		return
	}
	utils.FreeObjId(uint32(txId))
	defer tx.Conn.Close()
	return tx.Tx.Commit()
}

func (txId Tx) getTx() (tx *TxData, err error) {
	obj := utils.GetObjById(uint32(txId))
	if obj == nil {
		err = fmt.Errorf("no Tx found")
		return
	}
	t, ok := obj.(*TxData)
	if !ok {
		err = fmt.Errorf("bad obj")
		return
	}
	tx = t
	return
}

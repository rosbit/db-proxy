package db

import (
	_ "github.com/go-sql-driver/mysql"
	"db-proxy/utils"
	"database/sql"
	"strings"
	"fmt"
	"context"
)

type DB uint32

func Open(dsn string) (dbId DB, err error) {
	var cId uint32
	cId, err = utils.IncRef(dsn, connectToDB, disconnectDB)
	dbId = DB(cId)
	return
}

func (dbId DB) Close() {
	utils.DecRefByObjId(uint32(dbId))
}

func (dbId DB) Ping() (err error) {
	dbInst, e := dbId.getConn()
	if e != nil {
		err = e
		return
	}
	return dbInst.Ping()
}

func (dbId DB) Prepare(query string) (stmtId Stmt, err error) {
	dbInst, e := dbId.getConn()
	if e != nil {
		err = e
		return
	}
	stmt, e := dbInst.Prepare(query)
	if e != nil {
		err = e
		return
	}
	sId := utils.NewObjId(stmt)
	stmtId = Stmt(sId)
	return
}

func (dbId DB) BeginTx(opts *sql.TxOptions) (txId Tx, err error) {
	dbConn, e := dbId.getConn()
	if e != nil {
		err = e
		return
	}
	// always get a connection from pool, then begin a transaction
	conn, e := dbConn.Conn(context.Background())
	if e != nil {
		err = e
		return
	}
	tx, e := conn.BeginTx(context.Background(), opts)
	if e != nil {
		conn.Close()
		err = e
		return
	}
	tId := utils.NewObjId(&TxData{Conn: conn, Tx: tx})
	txId = Tx(tId)
	return
}

func connectToDB(p interface{}) (res interface{}, err error) {
	dsn, _ := p.(string)
	s := strings.SplitN(dsn, ":", 2)
	if len(s) != 2 {
		err = fmt.Errorf("bad dsn")
		return
	}
	driverName, realDSN := s[0], s[1]
	switch driverName {
	default:
		err = fmt.Errorf("driver %s unsupported")
		return
	case "mysql":
	}

	dbConn, e := sql.Open(driverName, realDSN)
	if e != nil {
		err = e
		return
	}
	res = dbConn
	return
}

func disconnectDB(res interface{}) {
	dbConn, _ := res.(*sql.DB)
	dbConn.Close()
}

func (dbId DB) getConn() (dbInst *sql.DB, err error) {
	obj := utils.GetRefByObjId(uint32(dbId))
	if obj == nil {
		err = fmt.Errorf("no connection found")
		return
	}
	c, ok := obj.(*sql.DB)
	if !ok {
		err = fmt.Errorf("bad obj")
		return
	}
	dbInst = c
	return
}


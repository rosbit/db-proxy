package db

import (
	"db-proxy/utils"
	"database/sql"
	"fmt"
)

type Stmt uint32

func (stmtId Stmt) Close() (err error) {
	stmt, e := stmtId.getStmt()
	if e != nil {
		err = e
		return
	}
	utils.FreeObjId(uint32(stmtId))
	return stmt.Close()
}

func (stmtId Stmt) Exec(args ...any) (lastInsertId, rowsAffected int64, err error) {
	stmt, e := stmtId.getStmt()
	if e != nil {
		err = e
		return
	}

	enterQ("Exec")
	defer leaveQ("Exec")

	res, e := stmt.Exec(args...)
	if e != nil {
		err = e
		return
	}
	lastInsertId, _ = res.LastInsertId()
	rowsAffected, _ = res.RowsAffected()
	return
}

func (stmtId Stmt) Query(args ...any) (columns []string, columnTypes []*sql.ColumnType, it <-chan interface{}, err error) {
	stmt, e := stmtId.getStmt()
	if e != nil {
		err = e
		return
	}

	enterQ("Query")
	defer leaveQ("Query")
	rows, e := stmt.Query(args...)
	if e != nil {
		err = e
		return
	}
	return getRows(rows)
}

func getRows(rows *sql.Rows) (columns []string, columnTypes []*sql.ColumnType, it <-chan interface{}, err error) {
	if columns, err = rows.Columns(); err != nil {
		rows.Close()
		return
	}
	if columnTypes, err = rows.ColumnTypes(); err != nil {
		rows.Close()
		return
	}

	colNum := len(columns)
	scanArgs := make([]interface{}, colNum)
	res := make(chan interface{})
	go func() {
		for rows.Next() {
			row := make([]interface{}, colNum)
			for i, _ := range row {
				scanArgs[i] = &row[i]
			}

			if e := rows.Scan(scanArgs...); e != nil {
				break
			}
			for i, _ := range row {
				col := &row[i]
				switch v := (*col).(type) {
				case []byte:
					switch columnTypes[i].DatabaseTypeName() {
					case "VARCHAR", "TEXT", "NVARCHAR", "CHAR", "JSON":
						*col = string(v)
					default:
					}
				default:
				}
			}
			res <- row
		}
		close(res)
		rows.Close()
	}()

	it = res
	return
}

func (stmtId Stmt) getStmt() (stmt *sql.Stmt, err error) {
	obj := utils.GetObjById(uint32(stmtId))
	if obj == nil {
		err = fmt.Errorf("no Stmt found")
		return
	}
	s, ok := obj.(*sql.Stmt)
	if !ok {
		err = fmt.Errorf("bad obj")
		return
	}
	stmt = s
	return
}

package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func NewOpen(user, passwd, host, port, dbName string) *sql.DB {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, passwd, host, port, dbName)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic("open db fail:" + conn)
	}
	return db
}

func Qeury(db *sql.DB, query string) ([]map[string]interface{}, error) {
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Dao:query error")
	}

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return records, errors.Wrap(err, "row scan error")
		}
		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records, rows.Err()
}

package main

import (
	"database/sql"
	"fmt"
	"geekbangGo/workError/dao"
	"github.com/pkg/errors"
	"log"
)

func main() {
	db := dao.NewOpen("root", "123456", "127.0.0.1", "3306", "db_test")
	defer db.Close()
	query := "select * from user;"
	records, err := dao.Qeury(db, query)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("user no record")
	} else {
		log.Fatalf("other errors row scan,casue type: %T\ncasue vaule:%v \ncasue trace: %v\n",
			errors.Cause(err), errors.Cause(err), err)
	}
	fmt.Println(records)
}

package main

import (
	"database/sql"
	"fmt"
	"geekbangGo/workError/dao"
	"github.com/pkg/errors"
	"log"
)

func main() {
	db := dao.NewOpen("root", "", "127.0.0.1", "3306", "mysql")
	defer db.Close()
	query := "select * from user;"
	records, err := dao.Qeury(db, query)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("user no record")
	}
	if err != nil {
		log.Fatalf("other errors row scan,casue type: %T\ncasue vaule:%v \ncasue trace: %+v\n",
			errors.Cause(err), errors.Cause(err), err)
	}
	for _, v := range records {
		for key, value := range v {
			fmt.Printf("%v:%s\n", key, value)
		}
	}
}

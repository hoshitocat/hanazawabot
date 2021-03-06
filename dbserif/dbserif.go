package dbserif

import (
	"fmt"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Serif interface {
	Add(string) bool
	SelectRandom() (tweet string, error string)
}

type DBSerif struct {
}

func (u *DBSerif) Add(body string) bool {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("insert into serifs (body, created_at) values (?, ?)", body, time.Now())
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
	}

	defer db.Close()

	return true
}

func (u *DBSerif) SelectRandom() (tweet string, error string) {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("select * from serifs order by rand() limit 1;")
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
	}

	defer db.Close()

	id, body, created_at, updated_at := 0, "", "", ""
	for rows.Next() {
		err = rows.Scan(&id, &body, &created_at, &updated_at)
		if err != nil{
			return "", err.Error()
		}
	}

	return body, ""
}

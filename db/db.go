package ecom

import (
	_ "github.com/mattn/go-oci8"
	"database/sql"
	"errors"
)

var db *sql.DB

func ConnectOracle(openString string) (err error) {
	db, err = sql.Open("oci8", openString)
	if err != nil {
		return
	}
	if db == nil {
		return errors.New("DB is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	err = db.PingContext(ctx)
	cancel()
	if err != nil {
		fmt.Println("PingContext error is not nil:", err)
	}
	return
}

func CloseOracle() (err error) {
	err = db.Close()
	if err != nil {
		fmt.Println("Close error is not nil:", err)
	}
	return
}
package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() *sql.DB {
	dsn := "root@tcp(127.0.0.1:4000)/order_matching_engine?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Failed to open connection to TiDB: %v\n", err)
		return nil
	}

	// Check if the connection is alive
	if err := db.Ping(); err != nil {
		fmt.Printf("Failed to connect to TiDB: %v", err)
	}

	fmt.Println("Connected to TiDB successfully")
	return db
}

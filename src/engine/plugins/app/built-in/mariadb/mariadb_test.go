package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

func main() {
	// Create the database handle, confirm driver is present
	//user:password@tcp(localhost:5555)
	//db, _ := sql.Open("mysql", "userTO7:46P4mKgqd5AjqDac@tcp(localhost:3306)/sampledb")
	//defer db.Close()
	// LlmrCcCF6J528ufu
	//conn, err := sql.Open("mysql", "userTO7:46P4mKgqd5AjqDac@tcp(localhost:3306)/sampledb")
	//conn, err := sql.Open("mysql", "root:LlmrCcCF6J528ufu@tcp(localhost:3306)/sampledb")
	conn, err := sql.Open("mysql", "root@tcp(localhost:3306)/sampledb")
	checkErr(err)
	defer conn.Close()

	//statement, err := conn.Prepare("select title from posts limit 10") 
	query,err := conn.Query("flush tables with read lock")
	checkErr(err)

	fmt.Println("Locked",query)
	time.Sleep(10 * time.Second)

	query,err = conn.Query("unlock tables")
	checkErr(err)
	// Connect and check the server version
	var version string
	conn.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)

	res, _ := conn.Query("SHOW TABLES")
	var table string
	
	for res.Next() {
		res.Scan(&table)
		fmt.Println(table)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
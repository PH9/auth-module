package main

import (
	"fmt"
	"net/http"

	"database/sql"

	_ "github.com/lib/pq"
)

func authenticated(w http.ResponseWriter, r *http.Request) {
	var token = r.Header.Get("X-Authorization")

	if len(token) < 32 {
		panic(http.StatusUnauthorized)
	}

	if !isFoundTokenInDB(token) {
		panic(http.StatusUnauthorized)
	}

	fmt.Fprint(w, "{\"status\":\"success\"}")
}

func isFoundTokenInDB(token string) bool {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var findValidTokenQueryString = "SELECT token FROM login WHERE token = $1 and token_expire_date_time > NOW()"
	statement, err := db.Prepare(findValidTokenQueryString)
	if err != nil {
		panic(err)
	}

	if statement == nil {
		panic("findValidTokenQueryString statement is nil, Database may not able to connect.")
	}
	defer statement.Close()

	result, err := statement.Query(token)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	return result.Next()
}

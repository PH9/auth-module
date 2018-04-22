package main

import (
	"fmt"
	"net/http"

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
	var findValidTokenQueryString = "SELECT token FROM login WHERE token = $1 and token_expire_date_time > NOW()"
	rows, err := db.Query(findValidTokenQueryString, token)
	if err != nil {
		panic(err)
	}
	rows.Close()

	return rows.Next()
}

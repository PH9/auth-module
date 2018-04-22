package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/satori/go.uuid"
)

var dataSourceName string

func main() {
	fmt.Println("[I] Auth by golang is starting....")

	setupLogger()
	setupDatabase()
	setupTeminator()
	setupAPI()
}

func setupAPI() {
	config := configInstance()

	if os.Getenv("APP_NET_HTTP_PPROF") == "true" {
		go func() {
			fmt.Println(http.ListenAndServe(":6060", nil))
		}()
	}

	http.Handle("/authenticated", panicHandler{http.HandlerFunc(authenticated)})
	http.Handle("/version", panicHandler{http.HandlerFunc(version)})
	fmt.Println("[I] Auth by golang is started in port " + config.ApplicationPort)
	fmt.Println(http.ListenAndServe(":"+config.ApplicationPort, nil))
}

type panicHandler struct {
	http.Handler
}

func (h panicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	beginTime := time.Now()
	u, _ := uuid.NewV4()
	transactionID := strings.Replace(u.String(), "-", "", -1)
	r.Header.Set("Transaction-UUID", transactionID)

	logStep(1, r, beginTime, beginTime)

	defer func() {
		err := recover()
		if err != nil {
			if err == http.StatusUnauthorized {
				e := LogError{
					HTTPStatus: http.StatusUnauthorized,
					Step:       ErrorResponseToClient,
					Code:       StatusUnauthorized,
					Type:       BussinessError,
				}

				w.WriteHeader(http.StatusUnauthorized)
				logStepWithLogError(e, r, beginTime, time.Now())
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			logStepWithError(r, beginTime, time.Now(), fmt.Sprintf("%v", err))
			return
		}

		logStep(4, r, beginTime, time.Now())
	}()

	h.Handler.ServeHTTP(w, r)
}

func setupLogger() {
	log.SetFlags(0)
	c := configInstance()
	log.SetOutput(&lumberjack.Logger{
		Filename:  c.LogName,
		MaxSize:   c.LogMaxSizeBeforeArchiveInMB,
		MaxAge:    c.LogMaxAgeBeforeDeleteInDays,
		LocalTime: c.LogUseLocalTime,
		Compress:  c.LogCompress,
	})
}

var db *sql.DB

func setupDatabase() {
	fmt.Println("[I] Setting up database")
	config := configInstance()

	sslDisable := ""
	if os.Getenv("APP_DISABLE_DB_SSL") == "true" {
		sslDisable = "?sslmode=disable"
	}
	dataSourceName = fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s%s",
		config.Database.DatabaseUsername, config.Database.DatabasePassword,
		config.Database.DatabaseHost, config.Database.DatabasePort,
		config.Database.DatabaseName, sslDisable)

	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(5)
}

func setupTeminator() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Print("[X] Auth by golang is terminated!")
		os.Exit(-1)
	}()
}

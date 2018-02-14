package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func logStep(step LogStep, r *http.Request, beginTime time.Time, happenTime time.Time) {
	write(generateLogBody(step, r, beginTime, happenTime))
}

func logStepWithError(r *http.Request, beginTime time.Time, errorTime time.Time, e string) {
	write(fmt.Sprintf("%s|%s", generateLogError(r, beginTime, errorTime), e))
}

func logStepWithLogError(e LogError, r *http.Request, beginTime time.Time, errorTime time.Time) {
	write(generateLogLine(e, r, beginTime, errorTime))
}

func generateLogBody(step LogStep, r *http.Request, beginTime time.Time, endTime time.Time) string {
	e := LogError{
		Step: step,
	}
	return generateLogLine(e, r, beginTime, endTime)
}

func generateLogError(r *http.Request, beginTime time.Time, endTime time.Time) string {
	e := LogError{
		Step: ErrorResponseToClient,
		Type: TechnicalError,
	}

	return generateLogLine(e, r, beginTime, endTime)
}

func generateLogLine(e LogError, r *http.Request, beginTime time.Time, endTime time.Time) string {
	hostname, _ := os.Hostname()
	clientUDID := r.Header.Get("UDID")
	cliNumber, _ := encryptByConfig(r.Header.Get("Cli-Number"))
	usernameLogin, _ := encryptByConfig(r.Header.Get("User-Name-Login"))
	clientIPAddress := r.RemoteAddr // TODO: Is using X-Real-IP in header or request.RemoteAddr
	clientApplicationVersion := r.Header.Get("App-Version")
	clientOS := r.Header.Get("Client-OS")
	clientJailStatus := r.Header.Get("Jail-Status")
	clientModel := r.Header.Get("Client-Model")
	clientNetworkType := r.Header.Get("Network-Type")
	rtrCode := r.Header.Get("RTR-Code")
	retailerMobileNumber, _ := encryptByConfig(r.Header.Get("Sim-R-Number"))
	token := subToken(r.Header.Get("X-Authorization"))
	requestID := r.Header.Get("Transaction-UUID")
	requestURL := r.URL.Path
	legacyName := ""
	legacyReturnCode := ""
	deviceLocation, _ := encryptByConfig(r.Header.Get("Lat-Long"))
	customerNumber, _ := encryptByConfig(r.Header.Get("Customer-Mobile"))

	return fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%f|%s|%s|%s|%s|%s|%s",
		hostname,
		clientUDID,
		cliNumber,
		usernameLogin,
		clientIPAddress,
		clientApplicationVersion,
		clientOS,
		clientJailStatus,
		clientModel,
		clientNetworkType,
		rtrCode,
		retailerMobileNumber,
		token,
		requestID,
		requestURL,
		e.Step,
		endTime.Sub(beginTime).Seconds(),
		legacyName,
		legacyReturnCode,
		e.Type,
		emptyIfZero(e.Code),
		deviceLocation,
		customerNumber)
}

func write(text string) {
	writeLogWithTime(time.Now, text)
}

func writeError(err error) {
	write(err.Error())
}

func writeLogWithTime(logTime func() time.Time, text string) {
	log.Println(getTime(logTime) + "|" + text)
}

func getTime(now func() time.Time) string {
	if now == nil {
		now = time.Now
	}

	return fmt.Sprint(now().Format("2006-01-02 15:04:05.9999999"))
}

func subToken(s string) string {
	if len(s) >= 20 {
		return s[len(s)-20:]
	}

	return s
}

func emptyIfZero(i ErrorCode) string {
	if i == 0 {
		return ""
	}

	return fmt.Sprintf("%d", i)
}

package main

import (
	"testing"
	"time"
)

func Test_getTime(t *testing.T) {
	stubNow := func() time.Time { return time.Unix(1499999999, 99999) }
	result := getTime(stubNow)

	if result == "" {
		t.Error("getTime() should not be empty")
		return
	}

	expected := "2017-07-14 09:39:59.0000999"
	if result != expected {
		t.Error(
			"getTime()",
			"Expected: "+expected,
			"Actual: "+result)
		return
	}
}

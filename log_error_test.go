package main

import (
	"fmt"
	"testing"
)

func Test_mrtr_error_log_step_const(t *testing.T) {
	if ReceiveFromClient != 1 {
		t.Error("ReceiveFromClient should ne 1")
	}

	if CallThirdParty != 2 {
		fmt.Println(CallThirdParty)
		t.Error("CallThirdParty should be 2")
	}

	if ResponseFromThirdParty != 3 {
		t.Error("ResponseFromThirdParty should be 3")
	}

	if SuccessResponseToClient != 4 {
		t.Error("SuccessResponseToClient should be 4")
	}

	if ErrorResponseToClient != 5 {
		t.Error("SuccessResponseToClient should be 5")
	}
}

func Test_mrtr_error_type_const(t *testing.T) {
	if BussinessError != "B" {
		t.Error("SuccessResponseToClient should be B")
	}

	if TechnicalError != "T" {
		t.Error("SuccessResponseToClient should be T")
	}
}

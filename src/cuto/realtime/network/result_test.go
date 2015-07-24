package network

import "testing"

func TestSuccessResult(t *testing.T) {
	expected := `{"status":0,"message":"Success.","pid":1234,"network":{"instance":321,"name":"test"}}`
	actual := SuccessResult(1234, 321, "test")
	if actual != expected {
		t.Log("Unexpected result.")
		t.Log("EXPECTED:")
		t.Log(expected)
		t.Log("ACTUAL:")
		t.Log(actual)
		t.Fail()
	}
}

func TestMasterErrorResult(t *testing.T) {
	expected := `{"status":1,"message":"test","pid":1234,"network":{"instance":0,"name":""}}`
	actual := MasterErrorResult("test", 1234)
	if actual != expected {
		t.Log("Unexpected result.")
		t.Log("EXPECTED:")
		t.Log(expected)
		t.Log("ACTUAL:")
		t.Log(actual)
		t.Fail()
	}
}

func TestRealtimeErrorResult(t *testing.T) {
	expected := `{"status":2,"message":"test","pid":0,"network":{"instance":0,"name":""}}`
	actual := RealtimeErrorResult("test")
	if actual != expected {
		t.Log("Unexpected result.")
		t.Log("EXPECTED:")
		t.Log(expected)
		t.Log("ACTUAL:")
		t.Log(actual)
		t.Fail()
	}
}

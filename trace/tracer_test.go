package trace

import (
	"bytes"
	"testing"
)


func TestNew(t *testing.T){
	var buf bytes.Buffer
	tracer := New(&buf)

	if tracer == nil {
		t.Error("return from New should not be nil")
	} else {
		tracer.Trace("Test string to trace.")
		if buf.String() != "Test string to trace.\n" {
			t.Error("trace should not write '$'",buf.String())
		}
	}

}
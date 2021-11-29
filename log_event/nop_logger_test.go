package log_event

import (
	"context"
	"testing"
)

func TestNopLogger(t *testing.T) {
	t.Parallel()
	nopLogEventVec := NewNopLogEventVec()
	jaegerLogEventVec := NewJaegerLogEventVec(context.Background(), "")
	m := make(map[string]string)
	m["Label1"] = "tag1"
	m["Label2"] = "tag2"
	got, want := jaegerLogEventVec, nopLogEventVec
	if got != want {
		t.Errorf("LogEventVec.NewLogEventVec returned %#v, want %#v", got, want)
	}
	got1, want1:= jaegerLogEventVec.WithLabelValues(m), nopLogEventVec.WithLabelValues(m)
	if got1 != want1 {
		t.Errorf("LogEventVec.WithLabelValues returned %#v, want %#v", got1, want1)
	}

	jaegerLogEventVec.WithLabelValues(m).Log("")
	nopLogEventVec.WithLabelValues(m).Log("")

}
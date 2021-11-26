package logger_test
import (
    "context"
    "gitlab.mobvista.com/mtech/tracelog"
    "testing"
)

func TestTraceLogLogger(t *testing.T) {
    
    testLogger := tracelog.NewLoggerVec(tracelog.LoggerOpts{
        Name: "testlog",
    }, []string{"label1", "label2"})

    ctx := context.Background()
    testLogger.WithLabelValues("tag1", "tag2").LogContext(ctx, "logstr")
}

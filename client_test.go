package tracelog_test 
import (
    "context"
    "gitlab.mobvista.com/mtech/tracelog"
)

func TestTraceLogLogger(t *testing.T) {
    
    testLogger := tracelog.NewLoggerVec(tracelog.LoggerOpts{
        Name: "testlog",
    }, []string{"label1", "label2"})

    ctx = context.Context()
    testLogger.WithLabelValues("tag1", "tag2").LogContext(ctx, "logstr")
}

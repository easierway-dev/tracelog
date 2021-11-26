package logger

import (
	"github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/mtech/tracelog"
	"testing"
)

func TestInitTracingLog(t *testing.T) {
	convey.Convey("test init TracingLog", t, func() {
		// init tracee
		traceLogPath := "tracing.log"
		sampleRatio := 1.0
		traceConfig, err := tracelog.NewConfig(tracelog.WithServiceName("dsp_server"),
			tracelog.WithTraceLogPath(traceLogPath),
			tracelog.WithSampleRatio(sampleRatio),
		)
		if err != nil {
			panic(err)
		}
		if err := tracelog.Start(traceConfig); err != nil {
			panic(err)
		}
		defer tracelog.Shutdown(traceConfig)
	})
}

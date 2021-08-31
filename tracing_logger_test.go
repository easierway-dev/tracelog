package tracinglog

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInitTracingLog(t *testing.T) {
	convey.Convey("test init TracingLog", t, func() {
		// init tracee
		traceLogPath := "tracing.log"
		sampleRatio := 1.0
		traceConfig, err := NewConfig(WithServiceName("dsp_server"),
			WithTraceLogPath(traceLogPath),
			WithSampleRatio(sampleRatio),
		)
		if err != nil {
			panic(err)
		}
		if err := Start(traceConfig); err != nil {
			panic(err)
		}
		defer Shutdown(traceConfig)
	})
}

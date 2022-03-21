package logevent

import (
	"context"
	"gitlab.mobvista.com/mtech/tracelog/ctxutil"
)

type logEventVec interface {
	WithLabelValues(map[string]string) logEvent
}

type logEvent interface {
	Log(interface{})
}

func WithContext(ctx context.Context, name string) logEventVec {
	// if context span sampled
	// return jaeger_log (hard code)
	if ctxutil.IsSampledFromContext(ctx) {
		return NewLogrusLogEventVec(ctx, name)
	}

	// context span not sampled
	// return nopLogEvent
	return NewNopLogEventVec()
}

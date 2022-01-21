package log_event

import (
    "context"
    "gitlab.mobvista.com/mtech/tracelog"
)

type logEventVec interface{
    WithLabelValues(map[string]string) logEvent
}

type logEvent interface{
    Log(string)
}

func WithContext(ctx context.Context, name string) logEventVec {
    // if context span sampled
    // return jaeger_log (hard code)
    if tracelog.IsSampledFromContext(ctx) {
    return NewLogrusLogEventVec(ctx, name)
    }

    // context span not sampled
    // return nopLogEvent
    return NewNopLogEventVec()
}



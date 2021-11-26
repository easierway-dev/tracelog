package logger


import (
    "gitlab.mobvista.com/mtech/tracelog"
)

type logEvent interface{
    WithName(string)
    WithLabelValues(map[string]string)
    Log(string)
}

func WithContext(ctx context.Context) logEvent {
    // context span not sampled
    // return nopLogEvent
    if ! tracelog.SpanSampleFlagFromContext(ctx) {
        return NewNopLogEvent()
    }

    // if context span sampled
    // return jaeger_log (hard code)
    return NewJaegerLogEvent(ctx)
    


    // if context not sample
    // return nop_logger
    // else 
    // else 
    // return jaeger_logger (hard code)

}



package trace_log 

import (
    
)

/*
    参考prometheus的client
        prometheus的label类似于给日志打标签, 比较好用, 后续可以无缝切换到loki等日志系统

    NewLoggerVec 返回 LoggerVec结构体, LoggerVec结构体包含的方法有WithLabel
    
    
    NewTracerLoggerVec【待定】
    NewTracerVec【待定】

*/

type LoggerVec {
    
}

type Logger interface{
    Log(string)
    LogContext(context.Context, string)
}

type logger struct {

}

func NewLoggerVec(opts LoggerOpts, labelNames []string) *LoggerVec {
    return LoggerVec{
    }
}

func (v *LoggerVec)WithLabelValues(lvs ...string) Logger {

}

func (l *logger) Log(logStr string) {

}

func (l *logger) LogContext(ctx context.Context, logStr string) {

    LogContext(ctx, lostr, attrs)
}


/*
func NewTracerLoggerVec(opts TraceLogOpts, labelNames []string) *TracerLoggerVec {
    return &TraceLogVec {
        return newTracerLogger()
    }
}

func (v *TracerLoggerVec)WithLabelValues(lvs ...string) TracerLogger{

}

func (c *tracerLogger) Log(logStr string) {

}

func (c *tracerLogger) Trace() {

}
*/

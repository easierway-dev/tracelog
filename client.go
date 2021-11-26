package trace_log

import (
    "context"
    "encoding/json"
    "fmt"
    "go.opentelemetry.io/otel/attribute"
    "unicode/utf8"
)

/*
   参考prometheus的client
       prometheus的label类似于给日志打标签, 比较好用, 后续可以无缝切换到loki等日志系统

   NewLoggerVec 返回 LoggerVec结构体, LoggerVec结构体包含的方法有WithLabel


   NewTracerLoggerVec【待定】
   NewTracerVec【待定】

*/

type LoggerOpts Opts
type Opts struct {
    Name string
}
type LoggerVec struct {
    opts LoggerOpts
    labelNames []string
}

type Loggers interface{
    Log(string)
    LogContext(context.Context, string)
}

type logger struct {
    Name string
    Labels map[string]string
    keyValues []attribute.KeyValue
}

func NewLoggerVec(opts LoggerOpts, labelNames []string) *LoggerVec {
    for _, val := range labelNames {
        if !utf8.ValidString(val) {
            fmt.Errorf("label value %q is not valid UTF-8", val)
        }
    }
    return &LoggerVec{
        opts: opts,
        labelNames: labelNames,
    }
}

func (v *LoggerVec) WithLabelValues(lvs ...string) Loggers {
    if len(v.labelNames) != len(lvs){
        fmt.Errorf("LabelNames value and LabelValues is not equal ")
        panic(lvs)
    }
    for _, val := range lvs {
        if !utf8.ValidString(val) {
            fmt.Errorf("label value %q is not valid UTF-8", val)
        }
    }
    LabelValues := make(map[string]string)
    for  key,_:= range v.labelNames{
        LabelValues[v.labelNames[key]] = lvs[key]
    }
    return &logger{Name:v.opts.Name,Labels: LabelValues}
}

func (l *logger) Log(logStr string) {
    for key, value := range l.Labels {
        l.keyValues = append(l.keyValues, attribute.KeyValue{
            Key:   attribute.Key(key),
            Value: attribute.StringValue(value),
        })
    }
}
func TraceLog(ctx context.Context, name string, attrs map[string]string, logstr string){
    indent, _:= json.MarshalIndent(attrs, "", "\t\t")
    fmt.Printf(
        "LogInfo{\n\tname: %s,\n\tlogstr: %s,\n\tattrs: %s\n}\n",
        name,
        logstr,
        string(indent),
    )
}
func (l *logger) LogContext(ctx context.Context, logStr string) {
    l.Log(logStr)
    TraceLog(ctx,l.Name,l.Labels,logStr)
    LogContext(ctx, logStr,l.keyValues)
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

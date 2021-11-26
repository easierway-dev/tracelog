参考了:
1. prometheus的kv label的metrics设计和代码
2. go-kit/go的context结构体和log接口 https://github.com/go-kit/log/blob/main/logfmt_logger.go
3. 几个logger: 
* noplogger -- 啥也不干的logger
* jaegerlogger -- jaeger作为logger日志上报系统
* lokilogger -- loki的logger, [TODO]

注意:
* logger的选择权不需要暴露给用户
* 与其他的logger不同,由于是跟context结合的logger,因此上下文是必须的

有帮助的概念:
* Contextual logging
* zlog.Event

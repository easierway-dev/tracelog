package logger 

// ref:
// https://github.com/go-kit/log/blob/main/nop_logger.go

type nopLogEvent struct{}

func NewNopLogEvent() LogEvnet {
    return nopLogEvent{}
}

func (nopLogEvent) WithName(string) {return}
func (nopLogEvent) WithLabelValues(map[string]string) {return}
func (nopLogEvent) Log(string) {return}

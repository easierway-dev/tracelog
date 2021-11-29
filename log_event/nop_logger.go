package log_event

// ref:
// https://github.com/go-kit/log/blob/main/nop_logger.go

type nopLogEventVec struct{}

type nopLogEvent strut {}
func NewNopLogEventVec() logEventVec {
    return nopLogEventVec{}
}

func (nopLogEventVec) WithLabelValues(map[string]string) logEvent {return nopLogEvent{}}
func (nopLogEvent) Log(string) {return}

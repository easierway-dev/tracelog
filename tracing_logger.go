package tracelog

import (
	"github.com/sirupsen/logrus"
	"gitlab.mobvista.com/mtech/zlog"
	"io"
)

type Logger struct {
	logger zlog.Logger
}

// 初始化：设置日志打印路径、zlog的初始化
func initLogger(logPath string) (*Logger, error) {
	var tracingOps = &zlog.Ops{
		Path:         logPath,
		Format:       "TimeFormatter",
		ReportCaller: false,
	}

	logger, err := zlog.NewZLog(tracingOps)
	if err != nil {
		return nil, err
	}
	return &Logger{logger}, nil
}

func (l *Logger) Writer() io.Writer {
	return l.logger.(*logrus.Logger).Writer()
}

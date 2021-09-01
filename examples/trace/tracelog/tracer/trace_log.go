package tracer

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

type PlainFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := fmt.Sprintf(entry.Time.Format(f.TimestampFormat))
	return []byte(fmt.Sprintf("%s %s %s\n", f.LevelDesc[entry.Level], timestamp, entry.Message)), nil

}

func initLogger(logPath string) (*logrus.Logger, error) {
	var tracingLog = logrus.New()

	var plainFormatter = &PlainFormatter{}
	plainFormatter.TimestampFormat = "2006-01-02 15:04:05"
	plainFormatter.LevelDesc = []string{"PANC", "FATL", "ERRO", "WARN", "INFO", "DEBG"}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrap(err, logPath)
	}

	tracingLog.SetOutput(f)
	// tracingLog.SetOutput(os.Stdout)
	tracingLog.SetFormatter(plainFormatter)
	return tracingLog, nil
}

package log_event

import (
	lkh "github.com/gfremex/logrus-kafka-hook"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
	"os"
)


func AddES(url string) *logrus.Logger{
	logger := logrus.New()
	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		log.Panic(err)
	}
	// TODO url 还待解析
	hook, err := elogrus.NewAsyncElasticHook(client, url, logrus.DebugLevel, "mylog")
	if err != nil {
		log.Panic(err)
	}
	logger.Hooks.Add(hook)

	return logger
}
func AddKafka(url string) *logrus.Logger{
	// Create a new KafkaHook
	hook, err := lkh.NewKafkaHook(
		"kh",
		[]logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel},
		&logrus.JSONFormatter{},
		[]string{url},
	)

	if err != nil {
		panic(err)
	}

	// Create a new logrus.Logger
	logger := logrus.New()

	// Add hook to logger
	logger.Hooks.Add(hook)
	//l := logger.WithField("topics", []string{"first_topic"})
	return logger
}
func AddStdout() *logrus.Logger{
	logger := logrus.New()
	// 设置日志格式为json格式
	logger.SetFormatter(&log.JSONFormatter{})
	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	logger.SetOutput(os.Stdout)
	// 设置日志级别为warn以上
	logger.SetLevel(log.InfoLevel)
	return logger
}

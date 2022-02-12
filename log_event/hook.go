package log_event

import (
	"fmt"
	lkh "github.com/gfremex/logrus-kafka-hook"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
    "io/ioutil"
	log "github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
	"net/url"
	"os"
	"strings"
	"time"
)

func GetIndexNameFunc(key string) elogrus.IndexNameFunc {
	return func() string {
		return key + "-" + time.Now().Format("20060102")
	}
}
func AddES(Url string) *logrus.Logger{
	u, err := url.Parse(Url)
	if err != nil{
		fmt.Println("invalid url:",err.Error())
        return nil
	}
	// 设置了ES的健康检查味false
	client, err := elastic.NewClient(elastic.SetHealthcheck(false),elastic.SetSniff(false),elastic.SetURL(Url))
	if err != nil {
		fmt.Println("invalid client log event:",err.Error())
        return nil 
	}
	host := strings.Split(u.Host, ":")
	hook, err := elogrus.NewAsyncElasticHookWithFunc(client,host[0], log.DebugLevel, GetIndexNameFunc("trace_log"))
	if err != nil {
		fmt.Println("invalid hook log event:",err.Error())
        return nil
	}
	logger := logrus.New()
    logger.SetOutput(ioutil.Discard)
	logger.Hooks.Add(hook)
	return logger
}
func AddKafka(Url string) *logrus.Logger{
	// Create a new KafkaHook
	hook, err := lkh.NewKafkaHook(
		"kh",
		[]logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel},
		&logrus.JSONFormatter{},
		[]string{Url},
	)

	if err != nil {
		fmt.Println("invalid hook log event:",err.Error())
	}

	// Create a new logrus.Logger
	logger := logrus.New()
    log.SetOutput(ioutil.Discard)

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

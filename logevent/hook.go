package logevent

import (
	"context"
	"fmt"
	lkh "github.com/gfremex/logrus-kafka-hook"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"
)

// 按自定义时间格式标准输出
func GetIndexNameFunc(key string) elogrus.IndexNameFunc {
	return func() string {
		return key + "-" + time.Now().Format("20060102")
	}
}

// 添加ES日志配置
func AddES(Url, ESUserName, ESPassword string) *logrus.Logger {
	u, err := url.Parse(Url)
	if err != nil {
		fmt.Println("invalid url:", err.Error())
		return nil
	}
	esOpts := make([]elastic.ClientOptionFunc, 0)
	esOpts = append(esOpts, elastic.SetHealthcheck(false))
	esOpts = append(esOpts, elastic.SetURL(Url))
	esOpts = append(esOpts, elastic.SetSniff(false))
	if ESUserName != "" && ESPassword != "" {
		esOpts = append(esOpts, elastic.SetBasicAuth(ESUserName, ESPassword))
	}
	// 设置ES的健康检查为false
	client, err := elastic.NewClient(esOpts...)
	if err != nil {
		fmt.Println("invalid client log event:", err.Error())
		return nil
	}
	// 获取ES的主机地址
	client.IndexExists("trace_log").Do(context.Background())
	host := strings.Split(u.Host, ":")
	// 异步方法，当es出问题，不会影响到主流程的业务
	hook, err := elogrus.NewAsyncElasticHookWithFunc(client, host[0], log.DebugLevel, GetIndexNameFunc("trace_log"))
	if err != nil {
		fmt.Println("invalid hook log event:", err.Error())
		return nil
	}
	logger := logrus.New()
	// 设置本地不打印
	logger.SetOutput(ioutil.Discard)
	logger.Hooks.Add(hook)
	return logger
}

// 添加Kafka日志配置
func AddKafka(Url string) *logrus.Logger {
	// Create a new KafkaHook
	hook, err := lkh.NewKafkaHook(
		"kh",
		[]logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel},
		&logrus.JSONFormatter{},
		[]string{Url},
	)
	if err != nil {
		fmt.Println("invalid hook log event:", err.Error())
	}
	// Create a new logrus.Logger
	logger := logrus.New()
	// 设置本地不打印
	logger.SetOutput(ioutil.Discard)
	// Add hook to logger
	logger.Hooks.Add(hook)
	//l := logger.WithField("topics", []string{"first_topic"})
	return logger
}

// 添加Stdout日志配置
func AddStdout() *logrus.Logger {
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

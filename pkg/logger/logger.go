package logger

import (
	"fmt"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"os"
	"time"
)

const (
	InfoLevel  string = "info"
	WarnLevel  string = "warn"
	DebugLevel string = "debug"
	ErrorLevel string = "error"
	FatalLevel string = "fatal"
)

type Provider interface {
	// log at different level. Depending on the performance, might need to use the batch
	// list later.
	Infow(msg string, keysAndValues ...string)
	Warnw(msg string, keysAndValues ...string)
	Debugw(msg string, keysAndValues ...string)
	Errorw(msg string, keysAndValues ...string)
	// fatal logs a fatal error and exit with status 1.
	Fatal(msg string, err error)
	// close the logger safely.
	Close()
}

type provider struct {
	cfg      config.LoggerConfig
	producer *producer.Producer
}

func NewProvider(cfg config.LoggerConfig) Provider {
	stage := os.Getenv(constants.Stage)
	if stage == "" {
		stage = constants.StageLocal
	}

	if stage == constants.StageProd || stage == constants.StageStaging {
		return NewAliyunProvider(cfg)
	} else {
		return NewMockProvider()
	}
}

func NewAliyunProvider(cfg config.LoggerConfig) Provider {
	log := provider{cfg: cfg}
	log.Setup()
	return log
}

func (p provider) Setup() {
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = p.cfg.Endpoint
	producerConfig.AccessKeyID = p.cfg.AccessKeyID
	producerConfig.AccessKeySecret = p.cfg.AccessKeySecret
	instance := producer.InitProducer(producerConfig)
	p.producer = instance

	// start the instance
	instance.Start()
}

func (p provider) Infow(msg string, keysAndValues ...string) {
	p.log(InfoLevel, msg, keysAndValues...)
}

func (p provider) Warnw(msg string, keysAndValues ...string) {
	p.log(WarnLevel, msg, keysAndValues...)
}

func (p provider) Debugw(msg string, keysAndValues ...string) {
	p.log(DebugLevel, msg, keysAndValues...)
}

func (p provider) Errorw(msg string, keysAndValues ...string) {
	p.log(ErrorLevel, msg, keysAndValues...)
}

func (p provider) Fatal(msg string, err error) {
	p.log(FatalLevel, msg, err.Error())
	os.Exit(1)
}

func (p provider) Close() {
	p.producer.SafeClose()
}

/** util **/
func (p provider) log(topic string, msg string, keysAndValues ...string) {
	logtime := uint32(time.Now().Unix())

	num := len(keysAndValues) / 2
	if len(keysAndValues)%2 == 1 {
		num++
	}

	content := make([]*sls.LogContent, 0, num+1)

	// initial message
	msgKey := "message"
	content = append(content, &sls.LogContent{
		Key:   &msgKey,
		Value: &msg,
	})

	var key, value string
	i := -1
	for _, val := range keysAndValues {
		i++
		if i%2 == 0 {
			key = val
		} else {
			value = val
		}
		if i != 0 && i%2 == 0 {
			content = append(content, &sls.LogContent{
				Key:   &key,
				Value: &value,
			})
		}
	}

	// if there're more than even number messages, add it as unknown.
	if i%2 == 1 {
		unknownKey := "unknown"
		content = append(content, &sls.LogContent{
			Key: &unknownKey,
			// key is value now.
			Value: &key,
		})
	}

	// send it to logger service at aliyun.
	// topic might be different service name later. It is log level at this moment.
	// for debug purpose, use p.producer.SendLogWithContext().
	p.producer.SendLog(p.cfg.Project, p.cfg.Logstore, topic, "http://google.com", &sls.Log{
		Time:     &logtime,
		Contents: content,
	})
}

/* callback */
type Callback struct{}

func (callback Callback) Success(result *producer.Result) {
	attemptList := result.GetReservedAttempts() // 遍历获得所有的发送记录
	for _, attempt := range attemptList {
		fmt.Println(attempt)
	}
}

func (callback Callback) Fail(result *producer.Result) {
	fmt.Println(result.IsSuccessful())        // 获得发送日志是否成功
	fmt.Println(result.GetErrorCode())        // 获得最后一次发送失败错误码
	fmt.Println(result.GetErrorMessage())     // 获得最后一次发送失败信息
	fmt.Println(result.GetReservedAttempts()) // 获得producerBatch 每次尝试被发送的信息
	fmt.Println(result.GetRequestId())        // 获得最后一次发送失败请求Id
	fmt.Println(result.GetTimeStampMs())      // 获得最后一次发送失败请求时间
}

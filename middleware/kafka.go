package middleware

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"

	"github.com/labstack/echo"
	"github.com/shibingli/realclouds_go/utils"

	"crypto/tls"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	//DefaultKafkaVersion *
	DefaultKafkaVersion = sarama.V0_9_0_1
)

//KafkaDuMessage *
type KafkaDuMessage struct {
	Type int `json:"type" xml:"type" structs:"type" gorm:"column:type"` //1-公共消息  2-私有消息
	//1-评论 2-回复 3-点赞 4-关注 5-播客 6-下架 7-审核 8-金币 9-爱心值 10-推荐 11-活动 12-音频 13-文本 14-专辑
	Type2 int `json:"type_2" xml:"type_2" structs:"type_2" gorm:"column:type_2"`
	//0-通用 1-音频 2-文本 3-专辑 4-评论
	Type3 int `json:"type_3" xml:"type_3" structs:"type_3" gorm:"column:type_3"`
	//1-收到 2-发送 3-新增 4-更新 5-收入 6-支出 7-线上 8-线下 9-升级 10-降级 11-通过 12-驳回 13-审核中 14-首页 15-横幅
	Type4       int    `json:"type_4" xml:"type_4" structs:"type_4" gorm:"column:type_4"`
	Type5       int    `json:"type_5" xml:"type_5" structs:"type_5" gorm:"column:type_5"` //预留类型
	Type6       int    `json:"type_6" xml:"type_6" structs:"type_6" gorm:"column:type_6"` //预留类型
	ResourceID  string `json:"resource_id,omitempty" xml:"resource_id,omitempty" structs:"resource_id" gorm:"column:resource_id"`
	ResourceID2 string `json:"resource_id_2,omitempty" xml:"resource_id_2,omitempty" structs:"resource_id_2" gorm:"column:resource_id_2"`
	ResourceID3 string `json:"resource_id_3,omitempty" xml:"resource_id_3,omitempty" structs:"resource_id_3" gorm:"column:resource_id_3"`
	ResourceID4 string `json:"resource_id_4,omitempty" xml:"resource_id_4,omitempty" structs:"resource_id_4" gorm:"column:resource_id_4"`
	ResourceID5 string `json:"resource_id_5,omitempty" xml:"resource_id_5,omitempty" structs:"resource_id_5" gorm:"column:resource_id_5"`
	ResourceID6 string `json:"resource_id_6,omitempty" xml:"resource_id_6,omitempty" structs:"resource_id_6" gorm:"column:resource_id_6"`
	Value       string `json:"value,omitempty" xml:"value,omitempty" structs:"value" gorm:"column:value"`
	Value2      string `json:"value_2,omitempty" xml:"value_2,omitempty" structs:"value_2" gorm:"column:value_2"`
	Value3      string `json:"value_3,omitempty" xml:"value_3,omitempty" structs:"value_3" gorm:"column:value_3"`
	Value4      string `json:"value_4,omitempty" xml:"value_4,omitempty" structs:"value_4" gorm:"column:value_4"`
	Value5      string `json:"value_5,omitempty" xml:"value_5,omitempty" structs:"value_5" gorm:"column:value_5"`
	Value6      string `json:"value_6,omitempty" xml:"value_6,omitempty" structs:"value_6" gorm:"column:value_6"`
	Read        bool   `json:"read" xml:"read" structs:"read" gorm:"column:read"`
	State       int    `json:"state" xml:"state" structs:"state" gorm:"column:state"`
}

//KafkaMsg *
type KafkaMsg struct {
	Receiver string         `json:"receiver" xml:"receiver"`
	Message  KafkaDuMessage `json:"message" xml:"message"`
	Data     interface{}    `json:"data,omitempty" xml:"data,omitempty"`
	encoded  []byte         `json:"-" xml:"-"`
	err      error          `json:"-" xml:"-"`
}

func (k *KafkaMsg) ensureEncoded() {
	if k.encoded == nil && k.err == nil {
		k.encoded, k.err = json.Marshal(k)
	}
}

//Length *
func (k *KafkaMsg) Length() int {
	k.ensureEncoded()
	return len(k.encoded)
}

//Encode *
func (k *KafkaMsg) Encode() ([]byte, error) {
	k.ensureEncoded()
	return k.encoded, k.err
}

//Kafka *
type Kafka struct {
	BrokerList             []string
	SyncProducerCollector  sarama.SyncProducer
	AsyncProducerCollector sarama.AsyncProducer
}

//NewKafka *
func NewKafka(brokerList []string) (kafka *Kafka, err error) {

	if len(brokerList) == 0 {
		return nil, fmt.Errorf("%s", "Invalid broker data.")
	}

	kafka = &Kafka{
		BrokerList:             brokerList,
		SyncProducerCollector:  newSyncProducerCollector(brokerList),
		AsyncProducerCollector: newASyncProducerCollector(brokerList),
	}

	return kafka, nil
}

//Close *
func (k *Kafka) Close() error {
	if err := k.SyncProducerCollector.Close(); err != nil {
		log.Errorf("Failed to shut down sync producer collector cleanly", err)
	}

	if err := k.AsyncProducerCollector.Close(); err != nil {
		log.Errorf("Failed to shut down async producer collector cleanly", err)
	}
	return nil
}

//SyncSendMessage *
func (k *Kafka) SyncSendMessage(topic string, msg KafkaMsg, key ...string) (partition int32, offset int64, err error) {
	topic = strings.TrimSpace(topic)

	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: &msg,
	}

	if len(key) > 0 {
		producerMessage.Key = sarama.StringEncoder(utils.StringUtils(key[0]).MD5())
	} else {
		producerMessage.Key = sarama.StringEncoder(utils.StringUtils(utils.GenerateUUID()).MD5())
	}

	partition, offset, err = k.SyncProducerCollector.SendMessage(producerMessage)

	return
}

//ASyncSendMessage *
func (k *Kafka) ASyncSendMessage(topic string, msg KafkaMsg, key ...string) {
	topic = strings.TrimSpace(topic)

	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: &msg,
	}

	if len(key) > 0 {
		producerMessage.Key = sarama.StringEncoder(utils.StringUtils(key[0]).MD5())
	} else {
		producerMessage.Key = sarama.StringEncoder(utils.StringUtils(utils.GenerateUUID()).MD5())
	}

	k.AsyncProducerCollector.Input() <- producerMessage
}

func newSyncProducerCollector(brokerList []string) sarama.SyncProducer {

	config := sarama.NewConfig()

	tlsConfig := getKafkaTLSConfigByEnv()
	if tlsConfig != nil {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}

	config.Version = DefaultKafkaVersion
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
	}

	return producer
}

func newASyncProducerCollector(brokerList []string) sarama.AsyncProducer {
	config := sarama.NewConfig()

	tlsConfig := getKafkaTLSConfigByEnv()
	if tlsConfig != nil {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}

	config.Version = DefaultKafkaVersion
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
	}

	go func() {
		for err := range producer.Errors() {
			log.Errorf("Failed to write access log entry: %v", err)
		}
	}()

	return producer
}

func getKafkaTLSConfigByEnv() (t *tls.Config) {
	certFile := utils.GetENV("KAFKA_TLS_CERT")
	keyFile := utils.GetENV("KAFKA_TLS_KEY")
	caFile := utils.GetENV("KAFKA_TLS_CA")
	verifySSL := utils.GetENVToBool("KAFKA_TLS_VERIFYSSL")

	t = utils.CreateTLSConfig(certFile, keyFile, caFile, verifySSL)

	return
}

//MwKafka Kafa middleware
func (k *Kafka) MwKafka(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("kafka", k)
		return next(c)
	}
}

//Subscription *
func (k *Kafka) Subscription(topics []string, group string,
	onMessage func(topic string, partition int32, offset int64, key, value []byte) error) {

	config := cluster.NewConfig()
	config.Version = DefaultKafkaVersion
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	consumer, err := cluster.NewConsumer(k.BrokerList, group, topics, config)
	if err != nil {
		log.Printf("Kafka New consumer error: %s\n", err.Error())
	}
	defer consumer.Close()

	go func() {
		for err := range consumer.Errors() {
			log.Printf("Kafka New consumer errors: %s\n", err.Error())
		}
	}()

	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Kafka consumer notifications: %+v\n", ntf)
		}
	}()

	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				if nil != onMessage {
					err := onMessage(msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
					if err != nil {
						log.Printf("Kafka callback method error: %s\n", err.Error())
					} else {
						consumer.MarkOffset(msg, "")
					}
				}
			}
		}
	}
}

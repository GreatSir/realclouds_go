package middleware

import (
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
func (k *Kafka) SyncSendMessage(topic, msg string, key ...string) (partition int32, offset int64, err error) {
	topic = strings.TrimSpace(topic)
	msg = strings.TrimSpace(msg)

	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	if len(key) > 0 {
		producerMessage.Key = sarama.StringEncoder(utils.StringUtils(key[0]).MD5())
	} else {
		producerMessage.Key = sarama.StringEncoder(utils.StringUtils(msg).MD5())
	}

	partition, offset, err = k.SyncProducerCollector.SendMessage(producerMessage)

	return
}

//ASyncSendMessage *
func (k *Kafka) ASyncSendMessage(topic, msg string, key ...string) {
	topic = strings.TrimSpace(topic)
	msg = strings.TrimSpace(msg)

	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	if len(key) > 0 {
		producerMessage.Key = sarama.StringEncoder(utils.StringUtils(key[0]).MD5())
	} else {
		producerMessage.Key = sarama.StringEncoder(utils.StringUtils(msg).MD5())
	}

	k.AsyncProducerCollector.Input() <- producerMessage
}

func newSyncProducerCollector(brokerList []string) sarama.SyncProducer {

	config := sarama.NewConfig()

	tlsConfig := getKafakaTLSConfigByEnv()
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

	tlsConfig := getKafakaTLSConfigByEnv()
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

func getKafakaTLSConfigByEnv() (t *tls.Config) {
	certFile := utils.GetENV("KAFKA_TLS_CERT")
	keyFile := utils.GetENV("KAFKA_TLS_KEY")
	caFile := utils.GetENV("KAFKA_TLS_CA")
	verifySSL := utils.GetENVToBool("KAFKA_TLS_VERIFYSSL")

	t = utils.CreateTLSConfig(certFile, keyFile, caFile, verifySSL)

	return
}

//MwKafaka Kafa middleware
func (k *Kafka) MwKafaka(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("kafaka", k)
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

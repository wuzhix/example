package queue

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/json-iterator/go"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	once     sync.Once
	producer sarama.AsyncProducer
	brokers  = []string{"127.0.0.1:9092"}
	group    = "test_group"
	topic    = "test_topic"
)

func CreateAsyncProducer(addrs []string) sarama.AsyncProducer {
	// 单例模式，once.Do只执行一次
	once.Do(func() {
		conf := sarama.NewConfig()
		conf.Producer.Compression = sarama.CompressionSnappy   // Compress messages
		conf.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
		producer, _ = sarama.NewAsyncProducer(addrs, conf)
	})
	return producer
}

func ProducerData() {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data := "2019-04-11 15:31:00"
	b, _ := json.Marshal(data)
	kafka := CreateAsyncProducer(brokers)
	for {
		kafka.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(b),
		}
		time.Sleep(1 * time.Second)
	}
}

func ConsumeData() {
	cfg := cluster.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Group.Return.Notifications = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	consumer, err := cluster.NewConsumer(brokers, group, []string{topic}, cfg)
	if err != nil {
		log.Println("cluster.NewConsumer: ", err.Error())
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Println("Error: ", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	// consume messages, watch signals
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
			}
		case <-signals:
			return
		}
	}
}

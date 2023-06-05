package mq

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/nsqio/go-nsq"
	mylog "github.com/wlbwlbwlb/log"
)

var consumers []*nsq.Consumer

var producer *nsq.Producer

var config = nsq.NewConfig()

var custom = Options{}

func Init(opts ...Option) (cancel func(), e error) {
	defer func() {
		if e != nil {
			if len(custom.lookupdAddr) > 0 {
				for _, consumer := range consumers {
					// Gracefully stop the consumer.
					consumer.Stop()
				}
			}
			if len(custom.nsqdAddr) > 0 {
				// Gracefully stop the producer when appropriate (e.g. before shutting down the service)
				producer.Stop()
			}
		}
	}()

	for _, opt := range opts {
		opt(&custom)
	}

	//有订阅需求
	if len(custom.lookupdAddr) > 0 {
		var consumer *nsq.Consumer

		for _, route := range routes {
			consumer, e = nsq.NewConsumer(route.topic, route.channel, config)
			if e != nil {
				return
			}
			consumer.SetLogger(log.New(io.MultiWriter(os.Stderr, mylog.Writer()), "", log.Flags()), nsq.LogLevelInfo)

			// Set the Handler for messages received by this Consumer. Can be called multiple times.
			// See also AddConcurrentHandlers.
			consumer.AddHandler(route.handler)

			// Use nsqlookupd to discover nsqd instances.
			// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
			e = consumer.ConnectToNSQLookupd(custom.lookupdAddr)
			if e != nil {
				return
			}

			consumers = append(consumers, consumer)
		}
	}

	//有发布需求
	if len(custom.nsqdAddr) > 0 {
		producer, e = nsq.NewProducer(custom.nsqdAddr, config)
		if e != nil {
			return
		}
		producer.SetLogger(log.New(io.MultiWriter(os.Stderr, mylog.Writer()), "", log.Flags()), nsq.LogLevelInfo)
	}

	return func() {
		if len(custom.lookupdAddr) > 0 {
			for _, consumer := range consumers {
				// Gracefully stop the consumer.
				consumer.Stop()
			}
		}
		if len(custom.nsqdAddr) > 0 {
			// Gracefully stop the producer when appropriate (e.g. before shutting down the service)
			producer.Stop()
		}
	}, nil
}

func Sub(topic, channel string, handler nsq.Handler) (e error) {
	routes = append(routes, Route{
		topic:   topic,
		channel: channel,
		handler: handler,
	})
	return
}

func Pub(topic string, body []byte) error {
	if nil == producer {
		return nil
	}
	// Synchronously publish a single message to the specified topic.
	// Messages can also be sent asynchronously and/or in batches.
	return producer.Publish(topic, body)
}

func DeferPub(topic string, delay time.Duration, body []byte) error {
	return producer.DeferredPublish(topic, delay, body)
}

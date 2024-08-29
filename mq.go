package mq

import (
	"log"

	"github.com/nsqio/go-nsq"
)

var consumers []*nsq.Consumer

var producer *nsq.Producer

var config = nsq.NewConfig()

var opt = Options{}

func Init(opts ...OptionFunc) (e error) {
	for _, fn := range opts {
		fn(&opt)
	}

	//有订阅需求
	if len(opt.lookupdAddr) > 0 {
		consumer := &nsq.Consumer{
			//
		}
		for _, route := range routes {
			consumer, e = nsq.NewConsumer(route.topic, route.channel, config)
			if e != nil {
				return
			}
			if opt.writer != nil {
				consumer.SetLogger(log.New(opt.writer, "", log.Flags()), nsq.LogLevelInfo)
			}

			// Set the Handler for messages received by this Consumer. Can be called multiple times.
			// See also AddConcurrentHandlers.
			consumer.AddHandler(route.handler)

			// Use nsqlookupd to discover nsqd instances.
			// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
			e = consumer.ConnectToNSQLookupd(opt.lookupdAddr)
			if e != nil {
				return
			}
			consumers = append(consumers, consumer)
		}
	}

	//有发布需求
	if len(opt.nsqdAddr) > 0 {
		producer, e = nsq.NewProducer(opt.nsqdAddr, config)
		if e != nil {
			return
		}
		if opt.writer != nil {
			producer.SetLogger(log.New(opt.writer, "", log.Flags()), nsq.LogLevelInfo)
		}
	}

	return
}

func StopProducer() {
	if len(opt.nsqdAddr) > 0 {
		// Gracefully stop the producer when appropriate (e.g. before shutting down the service)
		producer.Stop()
	}
}

func StopConsumers() {
	if len(opt.lookupdAddr) > 0 {
		for _, consumer := range consumers {
			// Gracefully stop the consumer.
			consumer.Stop()
		}
	}
}

func Sub(topic, channel string, handler nsq.Handler) (e error) {
	routes = append(routes, Route{
		topic:   topic,
		channel: channel,
		handler: handler,
	})
	return
}

//func PubAsync(topic string, body []byte) error {
//	if nil == producer {
//		return errors.New("init first")
//	}
//	return producer.PublishAsync(topic, body, nil)
//}

//func DeferPubAsync(topic string, delay time.Duration, body []byte) error {
//	if nil == producer {
//		return errors.New("init first")
//	}
//	return producer.DeferredPublishAsync(topic, delay, body, nil)
//}

package mq

import "github.com/nsqio/go-nsq"

type Route struct {
	topic   string
	channel string
	handler nsq.Handler
}

var routes []Route

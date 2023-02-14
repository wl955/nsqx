package nsqx

import (
	"errors"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/nsqio/go-nsq"
)

var requests = make(map[int]chan []byte)

func Rpc(from string, requestId int, body []byte) (res []byte, e error) {
	res = make([]byte, 0)

	if _, ok := stubs[from]; !ok {
		return res, errors.New("init stub first")
	}

	ch := make(chan []byte)

	requests[requestId] = ch

	req := struct {
		RequestId int    `json:"request_id"`
		Body      []byte `json:"body"`
	}{
		RequestId: requestId,
		Body:      body,
	}
	b, _ := jsoniter.Marshal(req)

	producer.Publish(from, b)

	select {
	case <-time.After(time.Second):
		e = errors.New("rpc timeout")
	case b = <-ch:
		res = b
	}

	return
}

var stubs = make(map[string]string)

func Stub(to, from string) (e error) {
	routes = append(routes, Route{
		topic:   to,
		channel: from,
		handler: nsq.HandlerFunc(handleMessage),
	})
	stubs[from] = to
	return
}

// HandleMessage implements the Handler interface.
func handleMessage(m *nsq.Message) (e error) {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return
	}

	// do whatever actual message processing is desired
	//err := processMessage(m.Body)
	fmt.Println(string(m.Body))

	req := struct {
		RequestId int    `json:"request_id"`
		Body      []byte `json:"body"`
	}{}
	jsoniter.Unmarshal(m.Body, &req)

	go func() {
		ch, ok := requests[req.RequestId]
		if !ok {
			return
		}

		select {
		case <-time.After(time.Second):
		case ch <- req.Body:
		}
		close(ch)

		delete(requests, req.RequestId)
	}()

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return
}

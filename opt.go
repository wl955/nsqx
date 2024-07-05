package mq

import "io"

type Options struct {
	lookupdAddr string
	nsqdAddr    string
	writer      io.Writer
}

type OptionFunc func(*Options)

func Lookupd(addr string) OptionFunc {
	return func(opts *Options) {
		opts.lookupdAddr = addr
	}
}

func Nsqd(addr string) OptionFunc {
	return func(opts *Options) {
		opts.nsqdAddr = addr
	}
}

func Writer(writer io.Writer) OptionFunc {
	return func(opts *Options) {
		opts.writer = writer
	}
}

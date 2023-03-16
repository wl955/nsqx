package mq

type Options struct {
	lookupdAddr string
	nsqdAddr    string
}

type Option func(*Options)

func Lookupd(addr string) Option {
	return func(opts *Options) {
		opts.lookupdAddr = addr
	}
}

func Nsqd(addr string) Option {
	return func(opts *Options) {
		opts.nsqdAddr = addr
	}
}

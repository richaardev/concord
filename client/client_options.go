package client

import "github.com/disgoorg/disgo/bot"

type ConcordOptions struct {
	disgoOpts []bot.ConfigOpt
}

type ConcordClientOption func(*ConcordOptions)

func WithDisgoOptions(opts ...bot.ConfigOpt) ConcordClientOption {
	return func(o *ConcordOptions) {
		o.disgoOpts = append(o.disgoOpts, opts...)
	}
}

func applyConcordOptions(opts []ConcordClientOption) ConcordOptions {
	options := ConcordOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	return options
}

package client

type Options struct {
	envFn     func() []string
	storePath string
}

type Option func(*Options)

func WithEnvFn(fn func() []string) Option {
	return func(o *Options) {
		o.envFn = fn
	}
}

func WithStorePath(path string) Option {
	return func(o *Options) {
		o.storePath = path
	}
}

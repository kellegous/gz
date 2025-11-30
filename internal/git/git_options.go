package git

type GitOptions struct {
	env []string
}

type GitOption func(*GitOptions)

func WithEnv(env ...string) GitOption {
	return func(o *GitOptions) {
		o.env = env
	}
}

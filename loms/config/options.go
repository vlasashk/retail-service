package config

type options struct {
	configPath string
}

func defaultOptions() options {
	return options{
		configPath: "config.yaml",
	}
}

type Option interface {
	apply(*options)
}

type optionf func(*options)

func (f optionf) apply(o *options) {
	f(o)
}

// WithCustomPath - adds suffix to a topic name.
//
// Use case scenario - if multiple handlers process identical topic,
// suffix will help to track progress separately by saving distinguish names in DB (if it is required)
func WithCustomPath(path string) Option {
	return optionf(func(c *options) {
		c.configPath = path
	})
}

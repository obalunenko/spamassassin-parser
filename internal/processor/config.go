package processor

// Config is a processor instance configuration.
type Config struct {
	Buffer  uint
	Receive struct {
		Response bool
		Errors   bool
	}
}

// NewConfig creates new config filled with sane default values.
func NewConfig() *Config {
	return &Config{
		Buffer: 0,
		Receive: struct {
			Response bool
			Errors   bool
		}{
			Response: true,
			Errors:   false,
		},
	}
}

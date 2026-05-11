package smsgate

type Factory struct {
	config Config
}

func NewFactory(config Config) *Factory {
	return &Factory{
		config: config,
	}
}

func (f *Factory) NewClient(username, password string) *Client {
	return NewClient(f.config, username, password)
}

package smsgate

type Factory struct {
	config  Config
	metrics *Metrics
}

func NewFactory(config Config, metrics *Metrics) *Factory {
	return &Factory{
		config:  config,
		metrics: metrics,
	}
}

func (f *Factory) NewClient(username, password string) *Client {
	return NewClient(f.config, username, password, f.metrics)
}

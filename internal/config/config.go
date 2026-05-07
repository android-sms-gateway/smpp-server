package config

import (
	"fmt"
	"os"

	"github.com/go-core-fx/config"
)

type http struct {
	Address     string   `koanf:"address"`
	ProxyHeader string   `koanf:"proxy_header"`
	Proxies     []string `koanf:"proxies"`

	OpenAPI openAPIConfig `koanf:"openapi"`
}

type openAPIConfig struct {
	Enabled    bool   `koanf:"enabled"`
	PublicHost string `koanf:"public_host"`
	PublicPath string `koanf:"public_path"`
}

type exampleConfig struct {
	Example string `koanf:"example"`
}

type smppConfig struct {
	BindAddress string `koanf:"bind_address"` // SMPP bind address (default: 127.0.0.1:2775)
	TLSCert     string `koanf:"tls_cert"`     // TLS certificate file path
	TLSKey      string `koanf:"tls_key"`      // TLS key file path

	SourceTON uint8 `koanf:"source_ton"` // Source Type of Number
	SourceNPI uint8 `koanf:"source_npi"` // Source Number Plan Indicator
	DestTON   uint8 `koanf:"dest_ton"`   // Destination Type of Number
	DestNPI   uint8 `koanf:"dest_npi"`   // Destination Number Plan Indicator
}

type gatewayConfig struct {
	APIBaseURL     string `koanf:"api_base_url"`     // Gateway API base URL
	WebhookBaseURL string `koanf:"webhook_base_url"` // Webhook base URL for delivery receipts
}

type Config struct {
	HTTP http `koanf:"http"`

	Example exampleConfig `koanf:"example"`

	SMPP    smppConfig    `koanf:"smpp"`
	Gateway gatewayConfig `koanf:"gateway"`
}

func Default() Config {
	//nolint:mnd // default values
	return Config{
		HTTP: http{
			Address:     "127.0.0.1:3000",
			ProxyHeader: "X-Forwarded-For",
			Proxies:     []string{},
			OpenAPI: openAPIConfig{
				Enabled:    true,
				PublicHost: "",
				PublicPath: "",
			},
		},

		Example: exampleConfig{
			Example: "example",
		},

		SMPP: smppConfig{
			BindAddress: "127.0.0.1:2775",
			TLSCert:     "",
			TLSKey:      "",

			SourceTON: 0x01,
			SourceNPI: 0x01,
			DestTON:   0x01,
			DestNPI:   0x01,
		},

		Gateway: gatewayConfig{
			APIBaseURL:     "https://api.sms-gate.app/3rdparty/v1",
			WebhookBaseURL: "",
		},
	}
}

func New() (Config, error) {
	cfg := Default()

	options := []config.Option{}
	if yamlPath := os.Getenv("CONFIG_PATH"); yamlPath != "" {
		options = append(options, config.WithLocalYAML(yamlPath))
	}

	if err := config.Load(&cfg, options...); err != nil {
		return Config{}, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}

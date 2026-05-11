package config

import (
	"github.com/android-sms-gateway/smpp-server/internal/sessions"
	"github.com/android-sms-gateway/smpp-server/internal/smpp"
	"github.com/android-sms-gateway/smpp-server/internal/smsgate"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/openapi"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"config",
		fx.Provide(New, fx.Private),
		fx.Provide(
			func(cfg Config) fiberfx.Config {
				return fiberfx.Config{
					Address:     cfg.HTTP.Address,
					ProxyHeader: cfg.HTTP.ProxyHeader,
					Proxies:     cfg.HTTP.Proxies,
				}
			},
			func(cfg Config) openapi.Config {
				return openapi.Config{
					Enabled:    cfg.HTTP.OpenAPI.Enabled,
					PublicHost: cfg.HTTP.OpenAPI.PublicHost,
					PublicPath: cfg.HTTP.OpenAPI.PublicPath,
				}
			},
		),
		fx.Provide(
			func(cfg Config) smpp.Config {
				return smpp.Config{
					BindAddress: cfg.SMPP.BindAddress,
					TLSCert:     cfg.SMPP.TLSCert,
					TLSKey:      cfg.SMPP.TLSKey,
				}
			},
			func(_ Config) sessions.Config {
				return sessions.Config{}
			},
			func(cfg Config) smsgate.Config {
				return smsgate.Config{
					BaseURL:    cfg.Gateway.BaseURL,
					WebhookURL: cfg.Gateway.WebhookBaseURL,
					Timeout:    cfg.Gateway.Timeout,
				}
			},
		),
	)
}

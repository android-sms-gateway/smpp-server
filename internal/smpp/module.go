package smpp

import (
	"github.com/go-core-fx/fxutil"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"smpp",
		logger.WithNamedLogger("smpp"),
		fx.Provide(NewService),
		fx.Invoke(fxutil.RegisterRunnable[*Service]()),
	)
}

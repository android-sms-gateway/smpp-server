package sessions

import (
	"github.com/go-core-fx/fxutil"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"sessions",
		logger.WithNamedLogger("sessions"),
		fx.Provide(NewMetrics, NewService),
		fx.Invoke(fxutil.RegisterRunnable[*Service]()),
	)
}

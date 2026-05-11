package smpp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/android-sms-gateway/smpp-server/internal/sessions"
	"go.uber.org/zap"
)

type Service struct {
	config Config

	sessionsSvc *sessions.Service

	listener net.Listener

	wg sync.WaitGroup

	logger *zap.Logger
}

func NewService(config Config, sessionsSvc *sessions.Service, logger *zap.Logger) *Service {
	return &Service{
		config: config,

		sessionsSvc: sessionsSvc,

		listener: nil,

		wg: sync.WaitGroup{},

		logger: logger,
	}
}

func (s *Service) Run(ctx context.Context) error {
	if s.config.TLSCert != "" && s.config.TLSKey != "" {
		cert, err := tls.LoadX509KeyPair(s.config.TLSCert, s.config.TLSKey)
		if err != nil {
			return fmt.Errorf("failed to load TLS certificate: %w", err)
		}

		s.listener, err = tls.Listen("tcp", s.config.BindAddress, &tls.Config{
			Certificates: []tls.Certificate{cert},
		})
		if err != nil {
			return fmt.Errorf("failed to create TLS listener: %w", err)
		}

		s.logger.Info("SMPP TLS server listening", zap.String("address", s.config.BindAddress))
	} else {
		var err error
		ln := net.ListenConfig{}
		s.listener, err = ln.Listen(ctx, "tcp", s.config.BindAddress)
		if err != nil {
			return fmt.Errorf("failed to create listener: %w", err)
		}

		s.logger.Info("SMPP server listening", zap.String("address", s.config.BindAddress))
	}

	s.wg.Go(func() {
		<-ctx.Done()
		s.logger.Info("stopping SMPP server")
		if err := s.listener.Close(); err != nil {
			s.logger.Error("failed to close listener", zap.Error(err))
		}
	})
	s.wg.Go(func() {
		s.acceptConnections(ctx)
	})

	s.wg.Wait()

	s.logger.Info("SMPP server stopped")

	return nil
}

func (s *Service) acceptConnections(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			s.logger.Error("failed to accept connection", zap.Error(err))
			continue
		}

		sess := s.sessionsSvc.NewSession(conn)
		s.wg.Go(func() {
			defer s.sessionsSvc.RemoveSession(sess.ID())

			sess.Run(ctx)
		})
	}
}

package sessions

import (
	"context"
	"net"
	"sync"

	"github.com/android-sms-gateway/smpp-server/internal/smsgate"
	"go.uber.org/zap"
)

type Service struct {
	config Config

	clientFactory *smsgate.Factory

	sessions map[string]*Session
	mux      sync.RWMutex

	metrics *Metrics
	logger  *zap.Logger
}

func NewService(config Config, clientFactory *smsgate.Factory, metrics *Metrics, logger *zap.Logger) *Service {
	return &Service{
		config: config,

		clientFactory: clientFactory,

		sessions: map[string]*Session{},
		mux:      sync.RWMutex{},

		metrics: metrics,
		logger:  logger,
	}
}

func (s *Service) Run(ctx context.Context) error {
	<-ctx.Done()
	s.Stop()
	return nil
}

func (s *Service) NewSession(conn net.Conn) *Session {
	s.mux.Lock()
	defer s.mux.Unlock()

	sess := newSession(
		conn,
		s.clientFactory.NewClient,
		s.metrics,
		s.logger,
	)

	s.sessions[sess.ID()] = sess

	return sess
}

func (s *Service) RemoveSession(id string) {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.sessions, id)
}

// Stop closes all active sessions and clears the sessions map.
// Safe to call multiple times; sessions already closed are skipped.
func (s *Service) Stop() {
	s.mux.Lock()
	defer s.mux.Unlock()

	for id, sess := range s.sessions {
		sess.close()
		delete(s.sessions, id)
	}
}

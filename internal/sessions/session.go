package sessions

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/android-sms-gateway/smpp-server/internal/smsgate"
	"github.com/fiorix/go-smpp/v2/smpp/pdu"
	"github.com/fiorix/go-smpp/v2/smpp/pdu/pdufield"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Session struct {
	id   string
	conn net.Conn

	state state

	clientFn func(username, password string) *smsgate.Client
	client   *smsgate.Client

	quit      chan struct{}
	closeOnce sync.Once

	writeMu sync.Mutex

	logger *zap.Logger
}

func newSession(conn net.Conn, clientFn func(username, password string) *smsgate.Client, logger *zap.Logger) *Session {
	id := uuid.New().String()[:8]
	logger = logger.With(zap.String("session", id))

	return &Session{
		id:   id,
		conn: conn,

		state: StateOpen,

		clientFn: clientFn,
		client:   nil,

		quit:      make(chan struct{}),
		closeOnce: sync.Once{},

		writeMu: sync.Mutex{},

		logger: logger,
	}
}

func (s *Session) Run(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
			s.close()
		case <-s.quit:
		}
	}()

	defer s.close()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.quit:
			return
		default:
		}

		if err := s.conn.SetReadDeadline(time.Now().Add(time.Minute)); err != nil {
			s.logger.Warn("Failed to set read deadline", zap.Error(err))
		}

		p, err := pdu.Decode(s.conn)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				s.logger.Debug("Read error", zap.Error(err))
			}
			return
		}

		s.handlePDU(ctx, p)
	}
}

// SendDeliverSM sends a DELIVER_SM PDU to the ESME client (for delivery receipts).
func (s *Session) SendDeliverSM(messageID string, messageState uint8) {
	resp := pdu.NewDeliverSM()
	resp.Header().Status = pdu.Status(ErrNoError)
	_ = resp.Fields().Set(pdufield.MessageID, messageID)
	_ = resp.Fields().Set(pdufield.MessageState, pdufield.Fixed{Data: messageState})
	s.writePDU(resp)
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) handlePDU(ctx context.Context, p pdu.Body) {
	h := p.Header()
	s.logger.Debug("Received PDU", zap.Uint32("id", uint32(h.ID)))

	switch h.ID {
	case pdu.BindTransmitterID:
		s.handleBind(ctx, p, false, false)
	case pdu.BindReceiverID:
		s.handleBind(ctx, p, true, false)
	case pdu.BindTransceiverID:
		s.handleBind(ctx, p, true, true)
	case pdu.SubmitSMID:
		s.handleSubmitSM(ctx, p)
	case pdu.QuerySMID:
		s.handleQuerySM(ctx, p)
	case pdu.UnbindID:
		s.handleUnbind(ctx, p)
	case pdu.EnquireLinkID:
		s.handleEnquireLink(ctx, p)
	// Known-but-unimplemented request PDUs - send error response
	case pdu.DataSMID:
		s.logger.Warn("Unsupported PDU", zap.Uint32("id", uint32(h.ID)))
		resp := pdu.NewDataSMResp()
		resp.Header().Status = pdu.Status(ErrInvalidCmdID)
		s.writeResponse(p, resp)
	case pdu.CancelSMID:
		s.logger.Warn("Unsupported PDU", zap.Uint32("id", uint32(h.ID)))
		resp := pdu.NewCancelSMResp()
		resp.Header().Status = pdu.Status(ErrCancelFail)
		s.writeResponse(p, resp)
	case pdu.ReplaceSMID:
		s.logger.Warn("Unsupported PDU", zap.Uint32("id", uint32(h.ID)))
		resp := pdu.NewReplaceSMResp()
		resp.Header().Status = pdu.Status(ErrReplaceFail)
		s.writeResponse(p, resp)
	case pdu.SubmitMultiID:
		s.logger.Warn("Unsupported PDU", zap.Uint32("id", uint32(h.ID)))
		resp := pdu.NewSubmitMultiResp()
		resp.Header().Status = pdu.Status(ErrSubmitFail)
		s.writeResponse(p, resp)
	// Response PDUs and notifications that don't require a response - log only
	case pdu.AlertNotificationID,
		pdu.DeliverSMID,
		pdu.GenericNACKID,
		pdu.BindReceiverRespID,
		pdu.BindTransmitterRespID,
		pdu.QuerySMRespID,
		pdu.SubmitSMRespID,
		pdu.DeliverSMRespID,
		pdu.UnbindRespID,
		pdu.ReplaceSMRespID,
		pdu.CancelSMRespID,
		pdu.BindTransceiverRespID,
		pdu.OutbindID,
		pdu.EnquireLinkRespID,
		pdu.SubmitMultiRespID,
		pdu.DataSMRespID:
		s.logger.Warn("Unhandled PDU", zap.Uint32("id", uint32(h.ID)))
	default:
		s.logger.Warn("Unhandled PDU", zap.Uint32("id", uint32(h.ID)))
	}
}

func (s *Session) handleBind(ctx context.Context, p pdu.Body, receiver, transceiverMode bool) {
	if s.state != StateOpen {
		s.sendBindResponse(p, receiver, transceiverMode, ErrInvalidBindStatus)
		return
	}

	fields := p.Fields()

	username, ok := fieldString(fields, pdufield.SystemID)
	if !ok || username == "" {
		s.logger.Warn("Missing or empty system_id in bind request")
		s.sendBindResponse(p, receiver, transceiverMode, ErrInvalidSystemID)
		return
	}

	password, ok := fieldString(fields, pdufield.Password)
	if !ok || password == "" {
		s.logger.Warn("Missing or empty password in bind request")
		s.sendBindResponse(p, receiver, transceiverMode, ErrInvalidPassword)
		return
	}

	var status = ErrNoError

	// Authenticate and get token
	client := s.clientFn(username, password)

	if err := client.Ping(ctx); err != nil {
		s.logger.Error("Ping failed", zap.Error(err))
		status = ErrBindFail
	}

	if status == ErrNoError && receiver {
		if err := client.RegisterWebhook(ctx, s.id); err != nil {
			s.logger.Error("Webhook registration failed", zap.Error(err))
			status = ErrBindFail
		}
	}

	if status == ErrNoError {
		s.client = client
		switch {
		case transceiverMode:
			s.state = StateBoundTRX
		case receiver:
			s.state = StateBoundRX
		default:
			s.state = StateBoundTX
		}
		s.logger.Info(
			"Client bound",
			zap.String("username", username),
		)
	} else {
		s.logger.Warn("Bind failed", zap.String("username", username))
	}

	s.sendBindResponse(p, receiver, transceiverMode, status)
}

// sendBindResponse sends the appropriate BIND response based on the bind type.
func (s *Session) sendBindResponse(req pdu.Body, receiver, transceiverMode bool, status uint32) {
	var resp pdu.Body
	switch {
	case transceiverMode:
		resp = pdu.NewBindTransceiverResp()
	case receiver:
		resp = pdu.NewBindReceiverResp()
	default:
		resp = pdu.NewBindTransmitterResp()
	}

	resp.Header().Status = pdu.Status(status)
	_ = resp.Fields().Set(pdufield.SystemID, "SMSGATE")
	s.writeResponse(req, resp)
}

func (s *Session) handleSubmitSM(ctx context.Context, req pdu.Body) {
	if s.state != StateBoundTX && s.state != StateBoundTRX {
		resp := pdu.NewSubmitSMResp()
		resp.Header().Status = pdu.Status(ErrInvalidBindStatus)
		s.writeResponse(req, resp)
		return
	}

	fields := req.Fields()

	destination, ok := fieldString(fields, pdufield.DestinationAddr)
	if !ok {
		s.logger.Warn("Missing destination_addr in submit_sm")
		resp := pdu.NewSubmitSMResp()
		resp.Header().Status = pdu.Status(ErrInvalidDstAddr)
		s.writeResponse(req, resp)
		return
	}

	shortMessage, ok := fieldString(fields, pdufield.ShortMessage)
	if !ok {
		s.logger.Warn("Missing short_message in submit_sm")
		resp := pdu.NewSubmitSMResp()
		resp.Header().Status = pdu.Status(ErrInvalidMsgLen)
		s.writeResponse(req, resp)
		return
	}

	source, ok := fieldString(fields, pdufield.SourceAddr)
	if !ok {
		s.logger.Warn("Missing source_addr in submit_sm")
		resp := pdu.NewSubmitSMResp()
		resp.Header().Status = pdu.Status(ErrInvalidSrcAddr)
		s.writeResponse(req, resp)
		return
	}

	deliveryReport, _ := fieldBool(fields, pdufield.RegisteredDelivery)

	var messageID string
	var status = ErrNoError

	result, err := s.client.SubmitSMS(ctx, smsgate.SubmitRequest{
		Source:         source,
		Destination:    destination,
		Content:        shortMessage,
		DeliveryReport: deliveryReport,
	})
	if err != nil {
		status = ErrSubmitFail
		s.logger.Error("Submit failed", zap.Error(err))
	} else {
		messageID = result.MessageID
	}

	resp := pdu.NewSubmitSMResp()
	resp.Header().Status = pdu.Status(status)
	_ = resp.Fields().Set(pdufield.MessageID, messageID)
	s.writeResponse(req, resp)
}

func (s *Session) handleQuerySM(ctx context.Context, req pdu.Body) {
	if s.state != StateBoundTX && s.state != StateBoundTRX {
		resp := pdu.NewQuerySMResp()
		resp.Header().Status = pdu.Status(ErrInvalidBindStatus)
		s.writeResponse(req, resp)
		return
	}

	fields := req.Fields()

	messageID, ok := fieldString(fields, pdufield.MessageID)
	if !ok {
		s.logger.Warn("Missing message_id in query_sm")
		resp := pdu.NewQuerySMResp()
		resp.Header().Status = pdu.Status(ErrInvalidMsgID)
		s.writeResponse(req, resp)
		return
	}

	var msgState = uint8(smsgate.MessageStateScheduled)
	var status = ErrNoError

	result, err := s.client.QuerySMS(ctx, messageID)
	if err != nil {
		s.logger.Error(
			"Query failed",
			zap.Error(err),
			zap.String("message_id", messageID),
		)
		status = ErrQueryFail
	} else {
		msgState = uint8(result.State)
	}

	s.logger.Debug(
		"Query result",
		zap.String("message_id", messageID),
		zap.Uint8("message_state", msgState),
	)

	resp := pdu.NewQuerySMResp()
	resp.Header().Status = pdu.Status(status)
	_ = resp.Fields().Set(pdufield.MessageState, msgState)
	s.writeResponse(req, resp)
}

func (s *Session) handleUnbind(ctx context.Context, req pdu.Body) {
	s.state = StateOpen
	s.logger.Info("Client unbound", zap.String("session", s.id))

	err := s.client.DeregisterWebhook(ctx)
	if err != nil {
		s.logger.Error("Webhook deregistration failed", zap.Error(err))
	}

	resp := pdu.NewUnbindResp()
	resp.Header().Status = pdu.Status(ErrNoError)
	s.writeResponse(req, resp)
	s.close()
}

func (s *Session) handleEnquireLink(_ context.Context, req pdu.Body) {
	resp := pdu.NewEnquireLinkResp()
	resp.Header().Status = pdu.Status(ErrNoError)
	s.writeResponse(req, resp)
}

func (s *Session) writeResponse(req pdu.Body, resp pdu.Body) {
	resp.Header().Seq = req.Header().Seq
	s.writePDU(resp)
}

func (s *Session) writePDU(p pdu.Body) {
	s.writeMu.Lock()
	err := p.SerializeTo(s.conn)
	s.writeMu.Unlock()
	if err != nil {
		s.logger.Error("Write error", zap.Error(err))
	}
}

func (s *Session) close() {
	s.closeOnce.Do(func() {
		close(s.quit)
		if err := s.conn.Close(); err != nil {
			s.logger.Error("Connection close failed", zap.Error(err))
		}

		// Cleanup on connection loss
		s.cleanup()

		s.logger.Debug("Session closed")
	})
}

// cleanup performs resource cleanup when a session ends unexpectedly.
func (s *Session) cleanup() {
	if s.client != nil {
		err := s.client.Cleanup(context.Background())
		if err != nil {
			s.logger.Error("Client cleanup failed", zap.Error(err))
		}
	}

	s.logger.Info("Session cleanup completed")
}

// fieldString safely extracts a string value from a PDU field map.
// Returns the field value and true if the field exists, or ("", false) otherwise.
func fieldString(fields map[pdufield.Name]pdufield.Body, name pdufield.Name) (string, bool) {
	f, ok := fields[name]
	if !ok {
		return "", false
	}
	return f.String(), true
}

func fieldBool(fields map[pdufield.Name]pdufield.Body, name pdufield.Name) (bool, bool) {
	f, ok := fields[name]
	if !ok {
		return false, false
	}

	if len(f.Bytes()) != 1 {
		return false, false
	}

	return f.Bytes()[0] == 1, true
}

package smsgate

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

type Client struct {
	config Config

	client *smsgateway.Client

	webhookID string
}

func NewClient(config Config, username, password string) *Client {
	return &Client{
		config: config,

		client: smsgateway.NewClient(smsgateway.Config{
			Client:   &http.Client{Timeout: config.Timeout},
			BaseURL:  config.BaseURL,
			User:     username,
			Password: password,
			Token:    "",
		}),

		webhookID: "",
	}
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.client.ListDevices(ctx)
	if err != nil {
		return fmt.Errorf("failed to list devices: %w", err)
	}

	return nil
}

func (c *Client) RegisterWebhook(ctx context.Context, sessionID string) error {
	if c.webhookID != "" {
		if err := c.DeregisterWebhook(ctx); err != nil {
			return fmt.Errorf("failed to replace existing webhook: %w", err)
		}
	}

	webhook := smsgateway.Webhook{
		ID:       "",
		DeviceID: nil,
		URL:      fmt.Sprintf("%s/api/smpp/v1/webhook?session=%s", c.config.WebhookURL, sessionID),
		Event:    smsgateway.WebhookEventSmsReceived,
	}

	res, err := c.client.RegisterWebhook(ctx, webhook)
	if err != nil {
		return fmt.Errorf("failed to register webhook: %w", err)
	}

	c.webhookID = res.ID

	return nil
}

func (c *Client) SubmitSMS(ctx context.Context, request SubmitRequest) (*SubmitResponse, error) {
	msg := smsgateway.Message{
		PhoneNumbers: []string{request.Destination},
		TextMessage: &smsgateway.TextMessage{
			Text: request.Content,
		},
		ID:                 "",
		DeviceID:           "",
		Message:            "",
		DataMessage:        nil,
		IsEncrypted:        false,
		SimNumber:          nil,
		WithDeliveryReport: &request.DeliveryReport,
		Priority:           0,
		TTL:                nil,
		ValidUntil:         nil,
		ScheduleAt:         nil,
	}

	res, err := c.client.Send(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return &SubmitResponse{
		MessageID: res.ID,
	}, nil
}

func (c *Client) QuerySMS(ctx context.Context, messageID string) (*QueryResponse, error) {
	res, err := c.client.GetState(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	state := MessageStateScheduled
	switch res.State {
	case smsgateway.ProcessingStateDelivered:
		state = MessageStateDelivered
	case smsgateway.ProcessingStateFailed:
		state = MessageStateUndeliverable
	case smsgateway.ProcessingStatePending:
		state = MessageStateScheduled
	case smsgateway.ProcessingStateProcessed, smsgateway.ProcessingStateSent:
		state = MessageStateEnroute
	}

	return &QueryResponse{
		MessageID: messageID,
		State:     state,
	}, nil
}

func (c *Client) DeregisterWebhook(ctx context.Context) error {
	if c.webhookID == "" {
		return nil
	}

	if err := c.client.DeleteWebhook(ctx, c.webhookID); err != nil {
		return fmt.Errorf("failed to deregister webhook: %w", err)
	}

	c.webhookID = ""

	return nil
}

func (c *Client) Cleanup(ctx context.Context) error {
	var err error

	err = errors.Join(err, c.DeregisterWebhook(ctx))

	return err
}

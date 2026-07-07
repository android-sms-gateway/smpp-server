package smsgate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

type Client struct {
	config Config

	client *smsgateway.Client

	webhookID string

	metrics *Metrics
}

func NewClient(config Config, username, password string, metrics *Metrics) *Client {
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
		metrics:   metrics,
	}
}

func (c *Client) Ping(ctx context.Context) error {
	defer c.metrics.StartRequest("ping")()

	_, err := c.client.ListDevices(ctx)
	if err != nil {
		c.metrics.IncRequest("ping", false)
		return fmt.Errorf("failed to list devices: %w", err)
	}

	c.metrics.IncRequest("ping", true)
	return nil
}

func (c *Client) RegisterWebhook(ctx context.Context, sessionID string) error {
	defer c.metrics.StartRequest("register_webhook")()

	if c.webhookID != "" {
		if err := c.DeregisterWebhook(ctx); err != nil {
			c.metrics.IncRequest("register_webhook", false)
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
		c.metrics.IncRequest("register_webhook", false)
		return fmt.Errorf("failed to register webhook: %w", err)
	}

	c.webhookID = res.ID

	c.metrics.IncRequest("register_webhook", true)
	return nil
}

func (c *Client) SubmitSMS(ctx context.Context, request SubmitRequest) (*SubmitResponse, error) {
	defer c.metrics.StartRequest("submit_sms")()

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
		c.metrics.IncRequest("submit_sms", false)
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	c.metrics.IncRequest("submit_sms", true)
	return &SubmitResponse{
		MessageID: res.ID,
	}, nil
}

func (c *Client) QuerySMS(ctx context.Context, messageID string) (*QueryResponse, error) {
	defer c.metrics.StartRequest("query_sms")()

	res, err := c.client.GetState(ctx, messageID)
	if err != nil {
		c.metrics.IncRequest("query_sms", false)
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	state := MessageStateScheduled
	switch res.State {
	case smsgateway.ProcessingStateDelivered:
		state = MessageStateDelivered
	case smsgateway.ProcessingStateFailed:
		state = MessageStateUndeliverable
	case smsgateway.ProcessingStatePending, smsgateway.ProcessingStateCancelling:
		state = MessageStateScheduled
	case smsgateway.ProcessingStateProcessed, smsgateway.ProcessingStateSent:
		state = MessageStateEnroute
	case smsgateway.ProcessingStateCancelled:
		state = MessageStateDeleted
	}

	c.metrics.IncRequest("query_sms", true)
	return &QueryResponse{
		MessageID: messageID,
		State:     state,
	}, nil
}

func (c *Client) DeregisterWebhook(ctx context.Context) error {
	if c.webhookID == "" {
		return nil
	}

	defer c.metrics.StartRequest("deregister_webhook")()

	if err := c.client.DeleteWebhook(ctx, c.webhookID); err != nil {
		c.metrics.IncRequest("deregister_webhook", false)
		return fmt.Errorf("failed to deregister webhook: %w", err)
	}

	c.webhookID = ""

	c.metrics.IncRequest("deregister_webhook", true)
	return nil
}

func (c *Client) Cleanup(ctx context.Context) error {
	defer c.metrics.StartRequest("cleanup")()

	err := c.DeregisterWebhook(ctx)
	if err != nil {
		c.metrics.IncRequest("cleanup", false)
		return err
	}

	c.metrics.IncRequest("cleanup", true)
	return nil
}

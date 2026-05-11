package smsgate

import "time"

type Config struct {
	BaseURL    string
	WebhookURL string
	Timeout    time.Duration
}

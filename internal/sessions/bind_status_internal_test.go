package sessions

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/android-sms-gateway/client-go/rest"

	"github.com/android-sms-gateway/smpp-server/internal/smsgate"
)

func TestBindStatusForPingError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want uint32
	}{
		{
			name: "nil error",
			err:  nil,
			want: ErrNoError,
		},
		{
			name: "rest client error",
			err:  fmt.Errorf("failed to list devices: %w", rest.ErrClient),
			want: ErrInvalidPassword,
		},
		{
			// Caveat: all 4xx client errors map to invalid password.
			name: "rest bad request error",
			err:  fmt.Errorf("failed to list devices: %w", rest.ErrBadRequest),
			want: ErrInvalidPassword,
		},
		{
			name: "rest server error",
			err:  fmt.Errorf("failed to list devices: %w", rest.ErrServer),
			want: ErrBindFail,
		},
		{
			name: "network failure",
			err:  errors.New("failed to make request: connection refused"),
			want: ErrBindFail,
		},
		{
			name: "deadline exceeded",
			err:  context.DeadlineExceeded,
			want: ErrBindFail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bindStatusForPingError(tt.err); got != tt.want {
				t.Errorf("bindStatusForPingError(%v) = %d (0x%X), want %d (0x%X)", tt.err, got, got, tt.want, tt.want)
			}
		})
	}
}

// TestBindStatusForPingErrorGatewayUnauthorized proves the real wrap chain
// preserves rest.ErrClient end to end: rest error -> smsgateway.ListDevices
// -> smsgate.Client.Ping, so the mapping relies on %w through the actual
// client library.
func TestBindStatusForPingErrorGatewayUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"unauthorized"}`))
	}))
	defer server.Close()

	client := smsgate.NewClient(smsgate.Config{
		BaseURL: server.URL,
		Timeout: time.Second,
	}, "user", "password", nil)

	err := client.Ping(t.Context())
	if err == nil {
		t.Fatal("expected Ping to fail, got nil")
	}

	if got := bindStatusForPingError(err); got != ErrInvalidPassword {
		t.Errorf("bindStatusForPingError(ping error) = %d (0x%X), want %d (0x%X); wrapped err = %v",
			got, got, ErrInvalidPassword, ErrInvalidPassword, err)
	}
}

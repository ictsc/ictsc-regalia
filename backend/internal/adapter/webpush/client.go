package webpush

import (
	"context"
	"fmt"
	"net/http"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

type Config struct {
	Subject    string
	PublicKey  string
	PrivateKey string
}

type Client struct {
	config Config
}

func New(config Config) *Client {
	return &Client{config: config}
}

func (c *Client) Send(ctx context.Context, subscription core.WebPushSubscription, payload []byte) (bool, error) {
	response, err := webpush.SendNotificationWithContext(ctx, payload, &webpush.Subscription{
		Endpoint: subscription.Endpoint,
		Keys: webpush.Keys{
			P256dh: subscription.P256DH,
			Auth:   subscription.Auth,
		},
	}, &webpush.Options{
		Subscriber:      c.config.Subject,
		VAPIDPublicKey:  c.config.PublicKey,
		VAPIDPrivateKey: c.config.PrivateKey,
		TTL:             300,
	})
	if err != nil {
		return false, fmt.Errorf("send Web Push: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return false, nil
	}
	if response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusGone {
		return true, fmt.Errorf("Web Push subscription expired with status %d", response.StatusCode)
	}
	return false, fmt.Errorf("Web Push service returned status %d", response.StatusCode)
}

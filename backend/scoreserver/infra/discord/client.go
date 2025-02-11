package discord

import "net/http"

type UserClient struct {
	HTTPClient *http.Client
}

package proxy

import (
	"context"
	"net/url"

	httpTransport "github.com/go-kit/kit/transport/http"
)

type Client struct {
	client *httpTransport.Client
}

func NewClient(rawUrl string, enc httpTransport.EncodeRequestFunc, dec httpTransport.DecodeResponseFunc) (*Client, error) {
	url, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	client := httpTransport.NewClient("POST", url, enc, dec)
	return &Client{client: client}, nil
}

func (c *Client) Do(request any) (response any, err error) {
	ep := c.client.Endpoint()
	response, err = ep(context.Background(), request)
	if err != nil {
		return
	}
	return
}

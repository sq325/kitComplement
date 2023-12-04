// Copyright 2023 Sun Quan
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxy

import (
	"context"
	"net/url"

	httpTransport "github.com/go-kit/kit/transport/http"
)

type Client struct {
	*httpTransport.Client
}

func NewClient(rawUrl string, enc httpTransport.EncodeRequestFunc, dec httpTransport.DecodeResponseFunc) (*Client, error) {
	url, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	client := httpTransport.NewClient("POST", url, enc, dec)
	return &Client{client}, nil
}

func (c *Client) Do(request any) (response any, err error) {
	ep := c.Endpoint()
	response, err = ep(context.Background(), request)
	if err != nil {
		return
	}
	return
}

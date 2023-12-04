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

package consul

import (
	"io"
	btlog "log"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/sq325/kitComplement/proxy"
)

var (
	retryMax     = 3
	retryTimeout = 10 * time.Second
)

// factor: url -> endpoint
// instance: ip:port
// path: /search
func FactoryFor(enc kithttp.EncodeRequestFunc, dec kithttp.DecodeResponseFunc, path string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance + path
		} else {
			instance = instance + path
		}
		btlog.Println("instance: ", instance)
		proxyClient, err := proxy.NewClient(instance, enc, dec)
		if err != nil {
			return nil, nil, err
		}
		return proxyClient.Endpoint(), nil, nil
	}
}

func NewEp(consulClient consulsd.Client, logger log.Logger, svcName string, dec kithttp.DecodeResponseFunc, path string) endpoint.Endpoint {
	enc := proxy.EncodeRequest
	factory := FactoryFor(enc, dec, path)
	instancer := consulsd.NewInstancer(consulClient, logger, svcName, nil, true)
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(retryMax, retryTimeout, balancer)
	return retry
}

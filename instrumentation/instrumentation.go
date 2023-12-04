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

package instrumentation

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"

	kitendpoint "github.com/go-kit/kit/endpoint"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	ReqC    metrics.Counter
	ReqErrC metrics.Counter
	ReqL    metrics.Histogram
}

func InstrumentingMiddleware(m Metrics) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request any) (resp any, err error) {
			method, ok := ctx.Value("method").(string)
			if !ok {
				method = "unknow"
			}

			uri, ok := ctx.Value("uri").(string)
			if !ok {
				uri = "unknow"
			}

			defer func() {
				if err != nil {
					m.ReqErrC.With("method", method, "uri", uri).Add(1)
				}
				m.ReqC.With("method", method, "uri", uri).Add(1)
			}()

			defer func(begin time.Time) {
				m.ReqL.With("method", method, "uri", uri).Observe(float64(time.Since(begin).Microseconds()))
			}(time.Now())

			return next(ctx, request)
		}
	}
}

func NewMetrics() Metrics {
	requestCounter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Name: "fileTransfer_request_count",
		Help: "Number of requests received.",
	}, []string{"uri", "method"})
	requestErrCounter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Name: "fileTransfer_request_error_count",
		Help: "Number of error requests received.",
	}, []string{"uri", "method"})
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Name:       "fileTransfer_request_latency_microseconds",
		Help:       "Total duration of requests in microseconds.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"uri", "method"})
	return Metrics{
		ReqC:    requestCounter,
		ReqErrC: requestErrCounter,
		ReqL:    requestLatency,
	}
}

// GinHandlerFunc gin handler
func GinHandlerFunc(method,
	uri string,
	ep kitendpoint.Endpoint,
	dec httptransport.DecodeRequestFunc,
	enc httptransport.EncodeResponseFunc,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		opt := httptransport.ServerBefore(
			func(ctx context.Context, _ *http.Request) context.Context {
				ctx = context.WithValue(ctx, "method", method)
				ctx = context.WithValue(ctx, "uri", uri)
				ctx = context.WithValue(ctx, "ClientIp", c.ClientIP())
				return ctx
			})

		h := httptransport.NewServer(
			ep,
			dec,
			enc,
			opt,
		)
		h.ServeHTTP(c.Writer, c.Request)
	}
}

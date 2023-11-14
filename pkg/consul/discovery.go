package consul

import (
	"io"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/sq325/kitComplement/pkg/proxy"
)

// factor: url -> endpoint
func FactorFor(url string, enc kithttp.EncodeRequestFunc, dec kithttp.DecodeResponseFunc) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance
		}

		proxyClient, err := proxy.NewClient(instance, enc, dec)
		if err != nil {
			return nil, nil, err
		}
		return proxyClient.Endpoint(), nil, nil
	}
}


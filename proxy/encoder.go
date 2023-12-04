package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// 把request -> http.r.body
func EncodeRequest(ctx context.Context, r *http.Request, req any) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req); err != nil {
		return err
	}
	r.Body = io.NopCloser(&buf)
	return nil
}

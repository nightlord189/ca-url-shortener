package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nightlord189/ca-url-shortener/pkg/log"
	"io"
	"net/http"
)

func responseJSON(ctx context.Context, w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	marshaled, err := json.Marshal(data)
	if err != nil {
		log.Ctx(ctx).Errorf("marshal response error: %v", err.Error())
	}
	if _, err = w.Write(marshaled); err != nil {
		log.Ctx(ctx).Errorf("write response error: %v", err.Error())
	}
}

func parseBodyJSON(r *http.Request, out interface{}) error {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read request body error: %w", err)
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Ctx(r.Context()).Errorf("close request body error: %v", err.Error())
		}
	}()

	if err := json.Unmarshal(rawBody, out); err != nil {
		return fmt.Errorf("unmarshal json error: %w", err)
	}

	return nil
}

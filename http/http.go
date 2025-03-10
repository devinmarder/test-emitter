package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
)

func Publisher(ctx context.Context, url string, headers http.Header, msgs <-chan string, log zerolog.Logger) error {
	client := &http.Client{}
	for msg := range msgs {
		log.Debug().Str("url", url).Str("msg", msg).Msg("publishing message")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(msg))
		if err != nil {
			return err
		}
		req.Header = headers
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			log.Warn().Int("status_code", resp.StatusCode).Msg("failed to publish message")
		}

		resp.Body.Close()
	}
	return nil
}

func ParseHeaders(headers string) (http.Header, error) {
	h := http.Header{}
	for _, header := range strings.Split(headers, ";") {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header: %s", header)
		}
		h.Add(parts[0], parts[1])
	}
	return h, nil
}

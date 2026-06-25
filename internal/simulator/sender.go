package simulator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func SendLog(ctx context.Context, url string, log TrafficLogPayload) error {
	body, err := json.Marshal(log)
	if err != nil {
		return errors.New("could not marshal traffic log err: " + err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return errors.New("could not create http put request err: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.New("could not send http put request err: " + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		if len(respBody) == 0 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		return errors.New("unexpected response: " + string(respBody))
	}

	return nil
}

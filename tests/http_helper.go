package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const requestTimeout = 45 * time.Second

type apiClient struct {
	baseURL    string
	httpClient *http.Client
}

type apiResponse struct {
	statusCode int
	body       []byte
}

func newAPIClient(baseURL string) *apiClient {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableKeepAlives = true

	return &apiClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout:   requestTimeout,
			Transport: transport,
		},
	}
}

func (c *apiClient) doJSON(t *testing.T, method, path string, body any, headers map[string]string) apiResponse {
	t.Helper()

	var rawBody []byte
	var err error
	if body != nil {
		rawBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal %s %s request: %v", method, path, err)
		}
	}

	requestURL := c.baseURL + path
	req, err := http.NewRequest(method, requestURL, bytes.NewReader(rawBody))
	if err != nil {
		t.Fatalf("create %s %s request: %v", method, requestURL, err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	for name, value := range headers {
		req.Header.Set(name, value)
	}

	t.Logf("request: %s %s headers=%s body=%s", method, requestURL, safeHeaders(req.Header), redactJSON(rawBody))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		t.Fatalf("send %s %s request: %v", method, requestURL, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		t.Fatalf("read %s %s response: %v", method, requestURL, err)
	}

	t.Logf("response: %s %s status=%d headers=%s body=%s", method, requestURL, resp.StatusCode, safeHeaders(resp.Header), redactJSON(responseBody))
	return apiResponse{statusCode: resp.StatusCode, body: responseBody}
}

func decodeJSON[T any](t *testing.T, response apiResponse) T {
	t.Helper()

	var result T
	if err := json.Unmarshal(response.body, &result); err != nil {
		t.Fatalf("decode JSON response: %v; body=%s", err, response.body)
	}

	return result
}

func requireStatus(t *testing.T, response apiResponse, expected int) {
	t.Helper()
	if response.statusCode != expected {
		t.Fatalf("HTTP status=%d, want %d; body=%s", response.statusCode, expected, response.body)
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func safeHeaders(headers http.Header) string {
	safe := make(http.Header)
	for _, name := range []string{"Accept", "Content-Type", "Authorization"} {
		if values := headers.Values(name); len(values) != 0 {
			for _, value := range values {
				safe.Add(name, value)
			}
		}
	}

	if safe.Get("Authorization") != "" {
		safe.Set("Authorization", "[redacted]")
	}

	return fmt.Sprintf("%v", safe)
}

func redactJSON(body []byte) string {
	if len(body) == 0 {
		return "<empty>"
	}

	var value any
	if err := json.Unmarshal(body, &value); err != nil {
		return string(body)
	}

	redactValue(value)
	redacted, err := json.Marshal(value)
	if err != nil {
		return "<unavailable>"
	}
	return string(redacted)
}

func redactValue(value any) {
	switch typed := value.(type) {
	case map[string]any:
		for key, item := range typed {
			switch strings.ToLower(key) {
			case "password", "access_token", "refresh_token":
				typed[key] = "[redacted]"
			default:
				redactValue(item)
			}
		}
	case []any:
		for _, item := range typed {
			redactValue(item)
		}
	}
}

package setup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	ApiScheme             = "http"
	ApiHost               = "localhost:8000"
	DefaultRequestTimeout = 5 * time.Second
)

type HttpReqBuilder struct {
	Method  string
	Url     string
	User    string
	Body    any
	Timeout time.Duration
}

func NewHttpReq(t *testing.T, params HttpReqBuilder) *http.Request {
	t.Helper()

	if params.Timeout == 0 {
		params.Timeout = DefaultRequestTimeout
	}

	//nolint:contextcheck,govet
	ctx, _ := context.WithTimeout(context.Background(), params.Timeout)

	var bodyReader io.Reader
	if params.Body != nil {
		jsonContent, err := json.Marshal(params.Body)
		if err != nil {
			t.Errorf("Failed to marshal request body: %v", err)
			t.FailNow()
			return nil
		}

		bodyReader = bytes.NewReader(jsonContent)
	}

	reqUrl, err := url.Parse(params.Url)
	if err != nil {
		t.Errorf("Failed to parse URL: %v", err)
		t.FailNow()
		return nil
	}

	if !reqUrl.IsAbs() {
		reqUrl.Scheme = ApiScheme
		reqUrl.Host = ApiHost
	}

	req, err := http.NewRequestWithContext(ctx, params.Method, reqUrl.String(), bodyReader)
	if err != nil {
		t.Errorf("Failed to create request: %v", err)
		t.FailNow()
		return nil
	}

	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if params.User != "" {
		req.Header.Set("x-user", params.User)
	}

	return req
}

func AssertStatusCode(t *testing.T, resp *http.Response, expected int) {
	t.Helper()

	if resp.StatusCode != expected {
		status := fmt.Sprintf("Expected: %d\nActual: %d", expected, resp.StatusCode)
		url := fmt.Sprintf("[%s] %s", resp.Request.Method, resp.Request.URL.String())

		respContent := ""
		if resp.Body != nil {
			rawRes, err := io.ReadAll(resp.Body)
			if err == nil {
				respContent = string(rawRes)
			}
		}
		resp := fmt.Sprintf("Response:\n%s", respContent)

		t.Errorf("%s\n%s\n%s\n%s", "HTTP request failed", status, url, resp)
		t.FailNow()
	}

	require.Equal(t, expected, resp.StatusCode, "Unexpected status code")
}

func ReadResponseBody[T any](resp *http.Response) (T, error) {
	var obj T

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return obj, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(body, &obj)
	if err != nil {
		return obj, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return obj, nil
}

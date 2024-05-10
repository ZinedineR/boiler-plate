package httputils

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

// DoHttpRequest is a function that sends an HTTP request.
// It accepts the HTTP method, target URL, headers, and data as parameters.
// It creates a new HTTP client and a new HTTP request with the provided parameters.
// If there is an error while creating the request, it returns nil and the error.
// It then sets the headers for the request.
// It sends the request and returns the response and any error that occurred.
// If there is an error while sending the request, it returns nil and the error.
func DoHttpRequest(ctx context.Context, method string, target string, headers map[string]any, data []byte) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, target, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, fmt.Sprint(value))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

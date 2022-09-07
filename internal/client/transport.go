package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type transport struct{}

type contextKey struct {
	name string
}

var ContextKeyRequestStart = &contextKey{"RequestStart"}

// RoundTrip is the core part of this module and implements http.RoundTripper.
// Executes HTTP request with request/response logging.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := context.WithValue(req.Context(), ContextKeyRequestStart, time.Now())
	req = req.WithContext(ctx)

	t.logRequest(req)

	resp, err := t.transport().RoundTrip(req)
	if err != nil {
		return resp, err
	}

	t.logResponse(resp)

	return resp, err
}

func (t *transport) logRequest(req *http.Request) {
	var body string = formatRequestBody(req)
	tflog.Trace(
		req.Context(),
		fmt.Sprintf("[%s] %s\n%s\n%s\n",
			req.Method,
			req.URL,
			formatHeaders(req.Header),
			body))
}

func (t *transport) logResponse(resp *http.Response) {
	var body string = formatResponseBody(resp)
	tflog.Trace(
		resp.Request.Context(),
		fmt.Sprintf("[%d] %s\n%s\n%s\n",
			resp.StatusCode,
			resp.Request.URL,
			formatHeaders(resp.Header), body))
}

func (t *transport) transport() http.RoundTripper {
	return http.DefaultTransport
}

func formatRequestBody(req *http.Request) string {
	if req.Body == nil {
		return ""
	}

	b, err := io.ReadAll(req.Body)
	if err != nil {
		return ""
	}

	req.Body = io.NopCloser(bytes.NewReader(b))

	return formatBody(b, req.Header.Get("Content-Type") == "application/json")
}

func formatResponseBody(res *http.Response) string {
	if res.Body == nil {
		return ""
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return ""
	}

	res.Body = io.NopCloser(bytes.NewReader(b))

	return formatBody(b, res.Header.Get("Content-Type") == "application/json")
}

func formatBody(body []byte, isJson bool) string {
	if !isJson {
		return string(body)
	}

	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "\t")
	if err == nil {
		return string(prettyJSON.Bytes())
	} else {
		return string(body)
	}
}

func formatHeaders(header http.Header) string {
	var strHeaders string = ""
	for k, v := range header {
		if k == "Authorization" {
			continue
		}
		strHeaders += fmt.Sprintf("%s: %s\n", k, strings.Join(v, ", "))
	}
	return strHeaders
}

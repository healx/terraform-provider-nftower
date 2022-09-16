package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type TowerClient struct {
	userAgent string
	apiKey    string
	apiUrl    *url.URL
	orgId     int64
	http      *retryablehttp.Client
}

func NewTowerClient(ctx context.Context, userAgent string, apiKey string, apiUrl string, org string) (*TowerClient, error) {
	u, _ := url.Parse(apiUrl)

	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil
	httpClient.RequestLogHook = (func(_ retryablehttp.Logger, req *http.Request, attempt int) {
		var body string = formatRequestBody(req)
		tflog.Trace(
			req.Context(),
			fmt.Sprintf("[%s] %s\n%s\n%s\n",
				req.Method,
				req.URL,
				formatHeaders(req.Header),
				body))
	})
	httpClient.ResponseLogHook = (func(_ retryablehttp.Logger, resp *http.Response) {
		var body string = formatResponseBody(resp)
		tflog.Trace(
			resp.Request.Context(),
			fmt.Sprintf("[%d] %s\n%s\n%s\n",
				resp.StatusCode,
				resp.Request.URL,
				formatHeaders(resp.Header), body))
	})

	c := &TowerClient{
		userAgent: userAgent,
		apiKey:    apiKey,
		apiUrl:    u.ResolveReference(&url.URL{Path: "/"}),
		http:      httpClient,
	}

	orgId, err := c.getOrgIdFromName(ctx, org)

	if err != nil {
		return nil, err
	}

	c.orgId = orgId

	return c, nil
}

func (c *TowerClient) getOrgIdFromName(ctx context.Context, orgName string) (int64, error) {
	tflog.Trace(ctx, fmt.Sprintf("Getting orgId from name for %s", orgName))
	res, err := c.request(ctx, "GET", "/orgs", nil, nil)

	if err != nil {
		return -1, err
	}

	if orgs, ok := res.(map[string]interface{}); ok {
		for _, org := range orgs["organizations"].([]interface{}) {
			o, _ := org.(map[string]interface{})
			if o["name"].(string) == orgName {
				return int64(o["orgId"].(float64)), nil
			}
		}
	}

	return -1, fmt.Errorf("Could not find an organization with the name %s", orgName)
}

func (c *TowerClient) request(ctx context.Context, method string, path string, query map[string]string, payload interface{}) (interface{}, error) {
	var r io.Reader
	if payload != nil {
		buf := &bytes.Buffer{}
		r = buf
		err := json.NewEncoder(buf).Encode(payload)
		if err != nil {
			return nil, err
		}
	}

	var querystring string = ""
	if query != nil {
		for k, v := range query {
			querystring += fmt.Sprintf("%s=%s", k, url.QueryEscape(v))
		}
	}

	req, err := retryablehttp.NewRequestWithContext(
		ctx,
		method,
		c.apiUrl.ResolveReference(&url.URL{Path: path, RawQuery: querystring}).String(),
		r)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("User-Agent", c.userAgent)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	httpResp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	var resp interface{}
	body, err := io.ReadAll(httpResp.Body)

	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode > 399 {
		return nil, fmt.Errorf("Tower API returned status: %s %s %s", httpResp.Status, httpResp.Request.URL, string(body))
	}

	if body == nil || len(body) == 0 {
		return body, nil
	}

	err = json.Unmarshal(body, &resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
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

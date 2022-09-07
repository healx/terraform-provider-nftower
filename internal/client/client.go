package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type TowerClient struct {
	userAgent string
	apiKey    string
	apiUrl    *url.URL
	orgId     int64
	http      *http.Client
}

func NewTowerClient(ctx context.Context, userAgent string, apiKey string, apiUrl string, org string) (*TowerClient, error) {
	u, _ := url.Parse(apiUrl)
	c := &TowerClient{
		userAgent: userAgent,
		apiKey:    apiKey,
		apiUrl:    u.ResolveReference(&url.URL{Path: "/"}),
		http: &http.Client{
			Timeout:   30 * time.Second,
			Transport: &transport{},
		},
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
		for k,v := range query {
			querystring += fmt.Sprintf("%s=%s", k, url.QueryEscape(v))
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.apiUrl.ResolveReference(&url.URL{Path: path, RawQuery: querystring}).String(), r)
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

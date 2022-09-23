package client

import (
	"context"
	"fmt"
)

func (c *TowerClient) CreateWorkspace(ctx context.Context, name string, fullName string, description string, visibility string) (int64, error) {

	payload := map[string]interface{}{
		"workspace": map[string]string{
			"name":        name,
			"fullName":    fullName,
			"description": description,
			"visibility":  visibility,
		},
	}

	res, err := c.requestWithJsonPayload(ctx, "POST", fmt.Sprintf("/orgs/%d/workspaces", c.orgId), nil, payload)

	if err != nil {
		return -1, err
	}

	if res == nil {
		return -1, fmt.Errorf("Empty response from server")
	}

	workspaceObj := res.(map[string]interface{})
	workspace := workspaceObj["workspace"].(map[string]interface{})

	return int64(workspace["id"].(float64)), nil
}

func (c *TowerClient) GetWorkspace(ctx context.Context, id int64) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", fmt.Sprintf("/orgs/%d/workspaces/%d", c.orgId, id), nil)

	if err != nil {
		return nil, err
	}

	workspace := res.(map[string]interface{})

	return workspace["workspace"].(map[string]interface{}), nil
}

func (c *TowerClient) GetWorkspaceByName(ctx context.Context, name string) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", fmt.Sprintf("/orgs/%d/workspaces", c.orgId), nil)

	if err != nil {
		return nil, err
	}

	if workspaces, ok := res.(map[string]interface{}); ok {
		for _, workspace := range workspaces["workspaces"].([]interface{}) {
			o, _ := workspace.(map[string]interface{})
			if o["name"].(string) == name {
				return c.GetWorkspace(ctx, int64(o["id"].(float64)))
			}
		}
	}

	return nil, fmt.Errorf("Could not find a workspace with the name '%s'", name)
}

func (c *TowerClient) DeleteWorkspace(ctx context.Context, id int64) error {
	_, err := c.requestWithoutPayload(ctx, "DELETE", fmt.Sprintf("/orgs/%d/workspaces/%d", c.orgId, id), nil)
	return err
}

func (c *TowerClient) UpdateWorkspace(ctx context.Context, id int64, fullName string, description string, visibility string) error {

	payload := map[string]interface{}{
		"fullName":    fullName,
		"description": description,
		"visibility":  visibility,
	}

	_, err := c.requestWithJsonPayload(ctx, "PUT", fmt.Sprintf("/orgs/%d/workspaces/%d", c.orgId, id), nil, payload)
	return err
}

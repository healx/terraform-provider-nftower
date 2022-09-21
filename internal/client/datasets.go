package client

import (
	"context"
	"fmt"
)

func (c *TowerClient) CreateDataset(ctx context.Context, workspaceId string, name string, description string) (string, error) {
	payload := map[string]interface{}{
		"name":        name,
		"description": description,
	}

	res, err := c.requestWithJsonPayload(ctx, "POST", fmt.Sprintf("/workspaces/%s/datasets", workspaceId), nil, payload)

	if err != nil {
		return "", err
	}

	datasetObj := res.(map[string]interface{})
	dataset := datasetObj["dataset"].(map[string]interface{})
	datasetId := dataset["id"].(string)

	return datasetId, nil
}

func (c *TowerClient) GetDataset(ctx context.Context, workspaceId string, id string) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", fmt.Sprintf("/workspaces/%s/datasets/%s/metadata", workspaceId, id), nil)

	if err != nil {
		return nil, err
	}

	datasetObj := res.(map[string]interface{})
	return datasetObj["dataset"].(map[string]interface{}), nil
}

func (c *TowerClient) UpdateDataset(ctx context.Context, workspaceId string, id string, name string, description string) error {
	payload := map[string]interface{}{
		"name":        name,
		"description": description,
	}

	_, err := c.requestWithJsonPayload(ctx, "PUT", fmt.Sprintf("/workspaces/%s/datasets/%s", workspaceId, id), nil, payload)

	return err
}

func (c *TowerClient) DeleteDataset(ctx context.Context, workspaceId string, id string) error {
	_, err := c.requestWithoutPayload(ctx, "DELETE", fmt.Sprintf("/workspaces/%s/datasets/%s", workspaceId, id), nil)
	return err
}

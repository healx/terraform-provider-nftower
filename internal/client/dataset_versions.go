package client

import (
	"context"
	"strings"
	"fmt"
)

func (c *TowerClient) CreateDatasetVersion(ctx context.Context, workspaceId string, datasetId string, fileContents string, filename string, hasHeader bool) (int, error) {
	body, contentType, err := c.prepareFilePayload(strings.NewReader(fileContents), filename)

	if err != nil {
		return -1, err
	}

	res, err := c.request(ctx, "POST", fmt.Sprintf("/workspaces/%s/datasets/%s/upload", workspaceId, datasetId), map[string]string{ "header": fmt.Sprintf("%t", hasHeader) }, body, contentType)

	if err != nil {
		return -1, err
	}

	versionObj := res.(map[string]interface{})
	version := versionObj["version"].(map[string]interface{})
	return int(version["version"].(float64)), nil
}

func (c *TowerClient) GetDatasetVersion(ctx context.Context, workspaceId string, datasetId string, versionId int) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", fmt.Sprintf("/workspaces/%s/datasets/%s/versions", workspaceId, datasetId), nil)

	if err != nil {
		return nil, err
	}

	versionsObj := res.(map[string]interface{})
	for _, v := range versionsObj["versions"].([]interface{}) {
		version := v.(map[string]interface{})
		if int(version["version"].(float64)) == versionId {
			contents, err := c.getDatasetContent(ctx, workspaceId, datasetId, versionId, version["fileName"].(string))
			if err != nil {
				return nil, err
			}
			version["contents"] = contents
			return version, nil
		}
	}

	return nil, fmt.Errorf("Could not find version %d for dataset %s", versionId, datasetId)
}

func (c *TowerClient) getDatasetContent(ctx context.Context, workspaceId string, datasetId string, versionId int, filename string) (string, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", fmt.Sprintf("/workspaces/%s/datasets/%s/v/%d/n/%s", workspaceId, datasetId, versionId, filename), nil)

	if err != nil {
		return "", err
	}

	return string(res.([]byte)), nil
}
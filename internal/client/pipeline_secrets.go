package client

import (
	"context"
	"fmt"
)

func (c *TowerClient) CreatePipelineSecrets(
	ctx context.Context,
	workspaceId string,
	name string,
	value string) (string, error) {

	payload := map[string]interface{}{
		"name":  name,
		"value": value,
	}

	res, err := c.requestWithJsonPayload(ctx, "POST", "/pipeline-secrets", map[string]string{"workspaceId": workspaceId}, payload)

	if err != nil {
		return "", err
	}

	if res == nil {
		return "", fmt.Errorf("empty response from server")
	}

	secrets := res.(map[string]interface{})

	secretId := fmt.Sprintf("%.0f", secrets["secretId"].(float64))

	return secretId, nil
}

func (c *TowerClient) GetPipelineSecretByName(ctx context.Context, workspaceId string, name string) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", "/pipeline-secrets", map[string]string{"workspaceId": workspaceId})

	if err != nil {
		return nil, err
	}

	if piplineSecretsList, ok := res.(map[string]interface{}); ok {
		for _, pipelineSecrets := range piplineSecretsList["pipelineSecrets"].([]interface{}) {
			o, _ := pipelineSecrets.(map[string]interface{})
			if o["name"].(string) == name {
				return c.GetPipelineSecret(ctx, workspaceId, o["id"].(string))
			}
		}
	}

	return nil, fmt.Errorf("could not find pipeline-secrets with the name '%s'", name)
}

func (c *TowerClient) GetPipelineSecret(ctx context.Context, workspaceId string, id string) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", fmt.Sprintf("/pipeline-secrets/%s", id), map[string]string{"workspaceId": workspaceId})

	if err != nil {
		if v, ok := err.(towerError); ok {
			if v.statusCode == 403 {
				// when the remote pipeline-secrets have been deleted,
				// tower returns a 403 instead of a 404 :(
				return nil, nil
			}
		}
		return nil, err
	}

	pipelineSecretsObj := res.(map[string]interface{})

	return pipelineSecretsObj["pipelineSecret"].(map[string]interface{}), nil
}

func (c *TowerClient) DeletePipelineSecrets(ctx context.Context, workspaceId string, id string) error {
	_, err := c.requestWithoutPayload(ctx, "DELETE", fmt.Sprintf("/pipeline-secrets/%s", id), map[string]string{"workspaceId": workspaceId})
	return err
}

func (c *TowerClient) UpdatePipelineSecrets(
	ctx context.Context,
	id string,
	workspaceId string,
	name string,
	value string,
) error {

	payload := map[string]interface{}{
		"name":  name,
		"value": value,
	}

	_, err := c.requestWithJsonPayload(ctx, "PUT", fmt.Sprintf("/pipeline-secrets/%s", id), map[string]string{"workspaceId": workspaceId}, payload)
	return err
}

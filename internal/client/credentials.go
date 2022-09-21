package client

import (
	"context"
	"fmt"
)

func (c *TowerClient) CreateCredentialsAWS(
	ctx context.Context,
	workspaceId string,
	name string,
	description string,
	accessKey string,
	secretKey string,
	assumeRoleArn string) (string, error) {

	payload := map[string]interface{}{
		"credentials": map[string]interface{}{
			"name":        name,
			"description": description,
			"provider":    "aws",
			"keys": map[string]interface{}{
				"accessKey":     accessKey,
				"secretKey":     secretKey,
				"assumeRoleArn": assumeRoleArn,
			},
		},
	}

	return c.createCredentials(ctx, workspaceId, payload)
}

func (c *TowerClient) CreateCredentialsGithub(
	ctx context.Context,
	workspaceId string,
	name string,
	description string,
	baseUrl string,
	username string,
	accessToken string) (string, error) {

	payload := map[string]interface{}{
		"credentials": map[string]interface{}{
			"name":        name,
			"description": description,
			"provider":    "github",
			"baseUrl":     baseUrl,
			"keys": map[string]interface{}{
				"username": username,
				"password": accessToken,
			},
		},
	}

	return c.createCredentials(ctx, workspaceId, payload)
}

func (c *TowerClient) createCredentials(ctx context.Context, workspaceId string, payload map[string]interface{}) (string, error) {
	res, err := c.requestWithJsonPayload(ctx, "POST", "/credentials", map[string]string{"workspaceId": workspaceId}, payload)

	if err != nil {
		return "", err
	}

	if res == nil {
		return "", fmt.Errorf("Empty response from server")
	}

	credentials := res.(map[string]interface{})

	return credentials["credentialsId"].(string), nil
}

func (c *TowerClient) GetCredentialsByName(ctx context.Context, workspaceId string, name string) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", "/credentials", map[string]string{"workspaceId": workspaceId})

	if err != nil {
		return nil, err
	}

	if credentialsList, ok := res.(map[string]interface{}); ok {
		for _, credentials := range credentialsList["credentials"].([]interface{}) {
			o, _ := credentials.(map[string]interface{})
			if o["name"].(string) == name {
				return c.GetCredentials(ctx, workspaceId, o["id"].(string))
			}
		}
	}

	return nil, fmt.Errorf("Could not find credentials with the name '%s'", name)
}

func (c *TowerClient) GetCredentials(ctx context.Context, workspaceId string, id string) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", fmt.Sprintf("/credentials/%s", id), map[string]string{"workspaceId": workspaceId})

	if err != nil {
		return nil, err
	}

	credentialsObj := res.(map[string]interface{})

	return credentialsObj["credentials"].(map[string]interface{}), nil
}

func (c *TowerClient) DeleteCredentials(ctx context.Context, workspaceId string, id string) error {
	_, err := c.requestWithoutPayload(ctx, "DELETE", fmt.Sprintf("/credentials/%s", id), map[string]string{"workspaceId": workspaceId})
	return err
}

func (c *TowerClient) UpdateCredentialsAWS(
	ctx context.Context,
	id string,
	workspaceId string,
	description string,
	accessKey string,
	secretKey string,
	assumeRoleArn string) error {

	payload := map[string]interface{}{
		"credentials": map[string]interface{}{
			"id":          id,
			"description": description,
			"provider":    "aws",
			"keys": map[string]interface{}{
				"accessKey":     accessKey,
				"secretKey":     secretKey,
				"assumeRoleArn": assumeRoleArn,
			},
		},
	}

	return c.updateCredentials(ctx, id, workspaceId, payload)
}

func (c *TowerClient) UpdateCredentialsGithub(
	ctx context.Context,
	id string,
	workspaceId string,
	description string,
	baseUrl string,
	username string,
	accessToken string) error {

	payload := map[string]interface{}{
		"credentials": map[string]interface{}{
			"id":          id,
			"description": description,
			"provider":    "github",
			"baseUrl":     baseUrl,
			"keys": map[string]interface{}{
				"username": username,
				"password": accessToken,
			},
		},
	}

	return c.updateCredentials(ctx, id, workspaceId, payload)
}

func (c *TowerClient) updateCredentials(ctx context.Context, id string, workspaceId string, payload map[string]interface{}) error {
	_, err := c.requestWithJsonPayload(ctx, "PUT", fmt.Sprintf("/credentials/%s", id), map[string]string{"workspaceId": workspaceId}, payload)
	return err
}

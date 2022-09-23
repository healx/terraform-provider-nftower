package client

import (
	"context"
	"fmt"
	"strconv"
)

func (c *TowerClient) CreateToken(ctx context.Context, name string) (string, string, error) {
	payload := map[string]interface{}{
		"name": name,
	}

	res, err := c.requestWithJsonPayload(ctx, "POST", "/tokens", nil, payload)

	if err != nil {
		return "", "", err
	}

	tokenObj := res.(map[string]interface{})
	token := tokenObj["token"].(map[string]interface{})

	return fmt.Sprintf("%d", int64(token["id"].(float64))), tokenObj["accessKey"].(string), nil
}

func (c *TowerClient) GetToken(ctx context.Context, id string) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", "/tokens", nil)

	if err != nil {
		return nil, err
	}

	tokensObj := res.(map[string]interface{})
	tokenId, _ := strconv.ParseInt(id, 10, 64)

	for _, v := range tokensObj["tokens"].([]interface{}) {
		token := v.(map[string]interface{})

		if int64(token["id"].(float64)) == tokenId {
			return token, nil
		}
	}

	return nil, fmt.Errorf("Could not find token with id %s", id)
}

func (c *TowerClient) DeleteToken(ctx context.Context, id string) error {
	_, err := c.requestWithoutPayload(ctx, "DELETE", fmt.Sprintf("/tokens/%s", id), nil)
	return err
}
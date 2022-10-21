package client

import (
	"context"
)

func (c *TowerClient) createLabels(ctx context.Context, workspaceId string, labels []string) ([]int64, error) {
	for _, v := range labels {
		payload := map[string]interface{}{
			"name": v,
		}

		_, err := c.requestWithJsonPayload(ctx, "POST", "/labels", map[string]string{"workspaceId": workspaceId}, payload)

		if err != nil {
			if v, ok := err.(towerError); ok {
				if v.statusCode == 409 {
					// label already exists
					continue
				}
			}
			return nil, err
		}
	}

	labelObjs, err := c.getLabels(ctx, workspaceId, labels)

	if err != nil {
		return nil, nil
	}

	labelIds := []int64{}

	for _, l := range labelObjs {
		lbl := l.(map[string]interface{})
		labelIds = append(labelIds, int64(lbl["id"].(float64)))
	}

	return labelIds, nil
}

func (c *TowerClient) getLabels(ctx context.Context, workspaceId string, labels []string) ([]interface{}, error) {
	// list all labels
	res, err := c.requestWithoutPayload(ctx, "GET", "/labels", map[string]string{"workspaceId": workspaceId})

	if err != nil {
		return nil, err
	}

	remoteLabels := res.(map[string]interface{})

	labelsToReturn := []interface{}{}

	for _, l := range labels {
		for _, v := range remoteLabels["labels"].([]interface{}) {
			rl := v.(map[string]interface{})
			if l == rl["name"] {
				labelsToReturn = append(labelsToReturn, v)
			}
		}
	}

	return labelsToReturn, nil
}

package client

import (
	"context"
	"fmt"
)

func (c *TowerClient) CreateWorkspaceParticipant(ctx context.Context, workspaceId string, memberId int64, role string) (int64, string, error) {

	payload := map[string]interface{}{
		"memberId": memberId,
	}

	res, err := c.requestWithJsonPayload(ctx, "PUT", fmt.Sprintf("/orgs/%d/workspaces/%s/participants/add", c.orgId, workspaceId), nil, payload)

	if err != nil {
		return -1, "", err
	}

	if res == nil {
		return -1, "", fmt.Errorf("Empty response from server")
	}

	participantObj := res.(map[string]interface{})
	participant := participantObj["participant"].(map[string]interface{})

	participantId := int64(participant["participantId"].(float64))

	err = c.UpdateWorkspaceParticipantRole(ctx, workspaceId, participantId, role)

	return participantId, participant["email"].(string), err
}

func (c *TowerClient) UpdateWorkspaceParticipantRole(ctx context.Context, workspaceId string, id int64, role string) error {
	payload := map[string]interface{}{
		"role": role,
	}

	_, err := c.requestWithJsonPayload(ctx, "PUT", fmt.Sprintf("/orgs/%d/workspaces/%s/participants/%d/role", c.orgId, workspaceId, id), nil, payload)
	return err
}

func (c *TowerClient) GetWorkspaceParticipant(ctx context.Context, workspaceId string, email string) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", fmt.Sprintf("/orgs/%d/workspaces/%s/participants", c.orgId, workspaceId), map[string]string{"search": email})

	if err != nil {
		return nil, err
	}

	participants := res.(map[string]interface{})

	if int64(participants["totalSize"].(float64)) == 0 {
		return nil, nil
	}

	participant := participants["participants"].([]interface{})

	return participant[0].(map[string]interface{}), nil
}

func (c *TowerClient) DeleteWorkspaceParticipant(ctx context.Context, workspaceId string, id int64) error {
	_, err := c.requestWithoutPayload(ctx, "DELETE", fmt.Sprintf("/orgs/%d/workspaces/%s/participants/%d", c.orgId, workspaceId, id), nil)
	return err
}

package client

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (c *TowerClient) CreateWorkspaceParticipant(ctx context.Context, workspaceId string, memberId int64, role string) (int64, string, error) {

	payload := map[string]interface{}{
		"memberId": memberId,
	}

	res, err := c.requestWithJsonPayload(ctx, "PUT", fmt.Sprintf("/orgs/%d/workspaces/%s/participants/add", c.orgId, workspaceId), nil, payload)

	tflog.Debug(ctx, "-----")
	tflog.Debug(ctx, fmt.Sprintf("workspaceId: %s, memberId: %d", workspaceId, memberId))
	tflog.Debug(ctx, fmt.Sprintf("%v", res))
	tflog.Debug(ctx, fmt.Sprintf("%v", err))
	tflog.Debug(ctx, "-----")

	participantExists := err != nil && err.Error() == fmt.Sprintf("Tower API returned status: 409 Conflict https://api.tower.nf/orgs/%d/workspaces/%s/participants/add {\"message\":\"Already a participant\"}", c.orgId, workspaceId)

	if !participantExists {
		return -1, "", err
	}

	if res == nil && !participantExists {
		return -1, "", fmt.Errorf("Empty response from server")
	}

	var participantObj map[string]interface{}

	if participantExists {
		ctx = tflog.SetField(ctx, "organizationId", c.orgId)
		ctx = tflog.SetField(ctx, "workspaceId", workspaceId)
		ctx = tflog.SetField(ctx, "memberId", memberId)
		tflog.Debug(ctx, "Member already exists, updating current state and role")

		participants, err := c.GetWorkspaceParticipants(ctx, workspaceId)

		if err != nil {
			return -1, "", err
		}

		if participants == nil {
			return -1, "", fmt.Errorf("Empty response from server")
		}

		var participant map[string]interface{}
		for _, value := range participants {
			p := value.(map[string]interface{})
			ctx = tflog.SetField(ctx, "participant", fmt.Sprintf("%v", p))
			ctx = tflog.SetField(ctx, "memberId", fmt.Sprintf("%d", memberId))
			tflog.Debug(ctx, "Checking if participant matches ID")
			if int64(p["memberId"].(float64)) == memberId {
				participant = p
				break
			}
		}

		if participant == nil {
			return -1, "", fmt.Errorf("No participant found with member ID: %d", memberId)
		}

		participantObj = map[string]interface{}{"participant": participant}
	}


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

func (c *TowerClient) GetWorkspaceParticipants(ctx context.Context, workspaceId string) ([]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", fmt.Sprintf("/orgs/%d/workspaces/%s/participants", c.orgId, workspaceId), make(map[string]string))

	if err != nil {
		return nil, err
	}

	participants := res.(map[string]interface{})

	if int64(participants["totalSize"].(float64)) == 0 {
		return nil, nil
	}

	return participants["participants"].([]interface{}), nil
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

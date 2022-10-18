package client

import (
	"context"
	"fmt"
)

func (c *TowerClient) CreateOrganizationMember(ctx context.Context, email string, role string) (int64, error) {

	payload := map[string]interface{}{
		"user": email,
	}

	res, err := c.requestWithJsonPayload(ctx, "PUT", fmt.Sprintf("/orgs/%d/members/add", c.orgId), nil, payload)

	if err != nil {
		return -1, err
	}

	if res == nil {
		return -1, fmt.Errorf("Empty response from server")
	}

	memberObj := res.(map[string]interface{})
	member := memberObj["member"].(map[string]interface{})

	memberId := int64(member["memberId"].(float64))

	err = c.UpdateOrganizationMemberRole(ctx, memberId, role)

	return int64(member["memberId"].(float64)), err
}

func (c *TowerClient) UpdateOrganizationMemberRole(ctx context.Context, id int64, role string) error {
	payload := map[string]interface{}{
		"role": role,
	}

	_, err := c.requestWithJsonPayload(ctx, "PUT", fmt.Sprintf("/orgs/%d/members/%d/role", c.orgId, id), nil, payload)
	return err
}

func (c *TowerClient) GetOrganizationMember(ctx context.Context, email string) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", fmt.Sprintf("/orgs/%d/members", c.orgId), map[string]string{"search": email})

	if err != nil {
		return nil, err
	}

	members := res.(map[string]interface{})

	if int64(members["totalSize"].(float64)) == 0 {
		return nil, nil
	}

	member := members["members"].([]interface{})

	return member[0].(map[string]interface{}), nil
}

func (c *TowerClient) DeleteOrganizationMember(ctx context.Context, id int64) error {
	_, err := c.requestWithoutPayload(ctx, "DELETE", fmt.Sprintf("/orgs/%d/members/%d", c.orgId, id), nil)
	return err
}

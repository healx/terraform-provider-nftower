package client

import (
	"context"
	"fmt"
)

func (c *TowerClient) CreatePipeline(
	ctx context.Context,
	workspaceId string,
	name string,
	description string,
	computeEnvironmentId string,
	pipeline string,
	workDir string,
	revision string,
	preRunScript string,
	postRunScript string,
	configProfiles []interface{},
	pipelineParameters string,
	nextflowConfig string,
	towerConfig string,
	mainScript string,
	workflowEntryName string,
	schemaName string,
	workspaceSecrets []interface{}) (int64, error) {

	launchPayload := map[string]interface{}{
		"computeEnvId": computeEnvironmentId,
		"pipeline":     pipeline,
		"workDir":      workDir,
	}

	payload := map[string]interface{}{
		"name":        name,
		"description": description,
		"launch": setOptionalPipelineFields(
			launchPayload,
			revision,
			preRunScript,
			postRunScript,
			configProfiles,
			pipelineParameters,
			nextflowConfig,
			towerConfig,
			mainScript,
			workflowEntryName,
			schemaName,
			workspaceSecrets),
	}

	res, err := c.requestWithJsonPayload(ctx, "POST", "/pipelines", map[string]string{"workspaceId": workspaceId}, payload)

	if err != nil {
		return -1, err
	}

	pipelineObj := res.(map[string]interface{})
	pl := pipelineObj["pipeline"].(map[string]interface{})

	return int64(pl["pipelineId"].(float64)), nil
}

func (c *TowerClient) GetPipeline(ctx context.Context, workspaceId string, id string) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", fmt.Sprintf("/pipelines/%s", id), map[string]string{"workspaceId": workspaceId})

	if err != nil {
		if v, ok := err.(towerError); ok {
			if v.statusCode == 403 {
				// when the remote pipeline has been deleted,
				// tower returns a 403 instead of a 404 :(
				return nil, nil
			}
		}
		return nil, err
	}

	pipelineObj := res.(map[string]interface{})
	pipeline := pipelineObj["pipeline"].(map[string]interface{})

	launch, err := c.getPipelineLaunchInfo(ctx, workspaceId, id)

	if err != nil {
		return nil, err
	}

	// merge the maps
	for k, v := range launch {
		pipeline[k] = v
	}

	return pipeline, nil
}

func (c *TowerClient) GetPipelineByName(ctx context.Context, workspaceId string, name string) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", "/pipelines", map[string]string{"workspaceId": workspaceId, "search": name})

	if err != nil {
		return nil, err
	}

	pipelines := res.(map[string]interface{})

	if int64(pipelines["totalSize"].(float64)) == 0 {
		return nil, nil
	}

	pipeline := pipelines["pipelines"].([]interface{})
	p := pipeline[0].(map[string]interface{})
	id := int64(p["pipelineId"].(float64))

	return c.GetPipeline(ctx, workspaceId, fmt.Sprintf("%d", id))
}

func (c *TowerClient) getPipelineLaunchInfo(ctx context.Context, workspaceId string, id string) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", fmt.Sprintf("/pipelines/%s/launch", id), map[string]string{"workspaceId": workspaceId})

	if err != nil {
		return nil, err
	}

	launchObj := res.(map[string]interface{})

	return launchObj["launch"].(map[string]interface{}), nil
}

func (c *TowerClient) DeletePipeline(ctx context.Context, workspaceId string, id string) error {
	_, err := c.requestWithoutPayload(ctx, "DELETE", fmt.Sprintf("/pipelines/%s", id), map[string]string{"workspaceId": workspaceId})
	return err
}

func (c *TowerClient) UpdatePipeline(
	ctx context.Context,
	workspaceId string,
	id string,
	description string,
	computeEnvironmentId string,
	pipeline string,
	workDir string,
	revision string,
	preRunScript string,
	postRunScript string,
	configProfiles []interface{},
	pipelineParameters string,
	nextflowConfig string,
	towerConfig string,
	mainScript string,
	workflowEntryName string,
	schemaName string,
	workspaceSecrets []interface{}) error {

	launchPayload := map[string]interface{}{
		"computeEnvId": computeEnvironmentId,
		"pipeline":     pipeline,
		"workDir":      workDir,
	}

	payload := map[string]interface{}{
		"description": description,
		"launch": setOptionalPipelineFields(
			launchPayload,
			revision,
			preRunScript,
			postRunScript,
			configProfiles,
			pipelineParameters,
			nextflowConfig,
			towerConfig,
			mainScript,
			workflowEntryName,
			schemaName,
			workspaceSecrets),
	}

	_, err := c.requestWithJsonPayload(ctx, "PUT", fmt.Sprintf("/pipelines/%s", id), map[string]string{"workspaceId": workspaceId}, payload)
	return err
}

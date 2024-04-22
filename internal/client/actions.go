package client

import (
	"context"
	"fmt"
)

func (c *TowerClient) CreateAction(
	ctx context.Context,
	workspaceId string,
	name string,
	source string,
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
	workspaceSecrets []interface{},
	labels []string) (string, error) {

	labelIds, err := c.createLabels(ctx, workspaceId, labels)

	if err != nil {
		return "", err
	}

	launchPayload := map[string]interface{}{
		"computeEnvId": computeEnvironmentId,
		"pipeline":     pipeline,
		"workDir":      workDir,
		"labelsIds":    labelIds,
	}

	payload := map[string]interface{}{
		"name":   name,
		"source": source,
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

	res, err := c.requestWithJsonPayload(ctx, "POST", "/actions", map[string]string{"workspaceId": workspaceId}, payload)

	if err != nil {
		return "", err
	}

	actionObj := res.(map[string]interface{})
	return actionObj["actionId"].(string), nil
}

func (c *TowerClient) GetAction(ctx context.Context, workspaceId string, id string) (map[string]interface{}, error) {
	res, err := c.requestWithoutPayload(ctx, "GET", fmt.Sprintf("/actions/%s", id), map[string]string{"workspaceId": workspaceId, "attributes": "labels"})

	if err != nil {
		return nil, err
	}

	actionObj := res.(map[string]interface{})

	return actionObj["action"].(map[string]interface{}), nil
}

func (c *TowerClient) UpdateAction(
	ctx context.Context,
	workspaceId string,
	id string,
	pipeline string,
	launchId string,
	computeEnvironmentId string,
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
	workspaceSecrets []interface{},
	labels []string) error {

	labelIds, err := c.createLabels(ctx, workspaceId, labels)

	if err != nil {
		return err
	}

	launchPayload := map[string]interface{}{
		"id":           launchId,
		"computeEnvId": computeEnvironmentId,
		"pipeline":     pipeline,
		"workDir":      workDir,
		"labelsIds":    labelIds,
	}

	payload := map[string]interface{}{
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

	_, err = c.requestWithJsonPayload(ctx, "PUT", fmt.Sprintf("/actions/%s", id), map[string]string{"workspaceId": workspaceId}, payload)
	return err
}

func (c *TowerClient) DeleteAction(ctx context.Context, workspaceId string, id string) error {
	_, err := c.requestWithoutPayload(ctx, "DELETE", fmt.Sprintf("/actions/%s", id), map[string]string{"workspaceId": workspaceId})
	return err
}

func setOptionalPipelineFields(
	payload map[string]interface{},
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
	workspaceSecrets []interface{}) map[string]interface{} {
	if revision != "" {
		payload["revision"] = revision
	}

	if configProfiles != nil {
		profiles := []string{}
		for _, v := range configProfiles {
			profiles = append(profiles, v.(string))
		}
		payload["configProfiles"] = profiles
	}

	if towerConfig != "" {
		payload["towerConfig"] = towerConfig
	}

	if nextflowConfig != "" {
		payload["configText"] = nextflowConfig
	}

	if mainScript != "" {
		payload["mainScript"] = mainScript
	}

	if preRunScript != "" {
		payload["preRunScript"] = preRunScript
	}

	if postRunScript != "" {
		payload["postRunScript"] = postRunScript
	}

	if workflowEntryName != "" {
		payload["entryName"] = workflowEntryName
	}

	if schemaName != "" {
		payload["schemaName"] = schemaName
	}

	if pipelineParameters != "" {
		payload["paramsText"] = pipelineParameters
	}

	if workspaceSecrets != nil {
		secrets := []string{}
		for _, v := range workspaceSecrets {
			secrets = append(secrets, v.(string))
		}
		payload["workspaceSecrets"] = secrets
	}

	return payload
}

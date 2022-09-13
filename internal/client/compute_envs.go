package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type ComputeEnvConfigEnvVar struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Head    bool   `json:"head"`
	Compute bool   `json:"compute"`
}

type ComputeEnvAWSBatchConfig struct {
	Region       string `json:"region"`
	ComputeQueue string `json:"computeQueue"`
	HeadQueue    string `json:"headQueue"`
	CliPath      string `json:"cliPath"`
	WorkDir      string `json:"workDir"`

	ExecutionRole  string `json:"executionRole,omitempty"`
	HeadJobRole    string `json:"headJobRole,omitempty"`
	ComputeJobRole string `json:"computeJobRole,omitempty"`

	PreRunScript    string                    `json:"preRunScript,omitempty"`
	PostRunScript   string                    `json:"postRunScript,omitempty"`
	HeadJobCpus     int                       `json:"headJobCpus"`
	HeadJobMemoryMb int                       `json:"headJobMemoryMb"`
	Environment     []*ComputeEnvConfigEnvVar `json:"environment,omitempty"`
}

func (c *TowerClient) CreateAWSBatchComputeEnv(
	ctx context.Context,
	workspaceId string,
	name string,
	description string,
	credentialsId string,
	config *ComputeEnvAWSBatchConfig) (string, error) {

	payload := map[string]interface{}{
		"computeEnv": map[string]interface{}{
			"name":          name,
			"description":   description,
			"platform":      "aws-batch",
			"credentialsId": credentialsId,
			"config":        marshalComputeEnvAWSBatchConfig(config),
		},
	}

	return c.createComputeEnv(ctx, workspaceId, payload)
}

func (c *TowerClient) createComputeEnv(ctx context.Context, workspaceId string, payload map[string]interface{}) (string, error) {
	res, err := c.request(ctx, "POST", "/compute-envs", map[string]string{"workspaceId": workspaceId}, payload)

	if err != nil {
		return "", err
	}

	if res == nil {
		return "", fmt.Errorf("Empty response from server")
	}

	computeEnv := res.(map[string]interface{})

	return computeEnv["computeEnvId"].(string), nil
}

func (c *TowerClient) GetComputeEnv(ctx context.Context, workspaceId string, id string) (map[string]interface{}, error) {
	res, err := c.request(ctx, "GET", fmt.Sprintf("/compute-envs/%s", id), map[string]string{"workspaceId": workspaceId}, nil)

	if err != nil {
		return nil, err
	}

	computeEnvObj := res.(map[string]interface{})
	computeEnv := computeEnvObj["computeEnv"].(map[string]interface{})

	switch computeEnv["platform"].(string) {
	case "aws-batch":
		config, err := unmarshalComputeEnvAWSBatchConfig(computeEnv["config"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		computeEnv["config"] = *config
	default:
		return nil, fmt.Errorf("unsupported platform: %s", computeEnv["platform"])
	}

	return computeEnv, nil
}

func (c *TowerClient) GetComputeEnvByName(ctx context.Context, workspaceId string, name string) (map[string]interface{}, error) {
	res, err := c.request(ctx, "GET", "/compute-envs", map[string]string{"workspaceId": workspaceId}, nil)

	if err != nil {
		return nil, err
	}

	if computeEnvs, ok := res.(map[string]interface{}); ok {
		for _, computeEnv := range computeEnvs["computeEnvs"].([]interface{}) {
			o, _ := computeEnv.(map[string]interface{})
			if o["name"].(string) == name {
				return c.GetComputeEnv(ctx, workspaceId, o["id"].(string))
			}
		}
	}

	return nil, fmt.Errorf("Could not find a computeEnv with the name '%s'", name)
}

func (c *TowerClient) DeleteComputeEnv(ctx context.Context, workspaceId string, id string) error {
	_, err := c.request(ctx, "DELETE", fmt.Sprintf("/compute-envs/%s", id), map[string]string{"workspaceId": workspaceId}, nil)
	return err
}

func unmarshalComputeEnvAWSBatchConfig(payload map[string]interface{}) (*ComputeEnvAWSBatchConfig, error) {
	var output ComputeEnvAWSBatchConfig

	b, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &output)

	if err != nil {
		return nil, err
	}

	return &output, nil
}

func marshalComputeEnvAWSBatchConfig(config *ComputeEnvAWSBatchConfig) map[string]interface{} {
	payload := map[string]interface{}{
		"region":       config.Region,
		"computeQueue": config.ComputeQueue,
		"headQueue":    config.HeadQueue,
		"workDir":      config.WorkDir,
		"cliPath":      config.CliPath,
		"environment":  config.Environment,
	}

	if config.ComputeJobRole != "" {
		payload["computeJobRole"] = config.ComputeJobRole
	}

	if config.HeadJobRole != "" {
		payload["headJobRole"] = config.HeadJobRole
	}

	if config.ExecutionRole != "" {
		payload["executionRole"] = config.ExecutionRole
	}

	if config.HeadJobCpus != 0 {
		payload["headJobCpus"] = config.HeadJobCpus
	}

	if config.HeadJobMemoryMb != 0 {
		payload["headJobMemoryMb"] = config.HeadJobMemoryMb
	}

	return payload
}

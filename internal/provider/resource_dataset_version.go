package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func resourceDatasetVersion() *schema.Resource {
	return &schema.Resource{
		Description: "A workspace inside a tower organization.",

		CreateContext: resourceDatasetVersionCreate,
		ReadContext:   resourceDatasetVersionRead,
		DeleteContext: schema.NoopContext,

		Schema: map[string]*schema.Schema{
			"dataset_id": {
				Description: "The id of the dataset to upload to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"workspace_id": {
				Description: "The id of the workspace in which the dataset lives.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"file_name": {
				Description: "The name of the file",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"contents": {
				Description: "The contents of the dataset. Must be CSV or TSV.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"has_header": {
				Description: "Whether the first row contains field headers.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"media_type": {
				Description: "The computed mime-type of the dataset file",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_updated": {
				Description: "The last updated datetime of the dataset.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"version": {
				Description: "The version number of the dataset.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"url": {
				Description: "The url of the dataset file.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceDatasetVersionCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.TowerClient)

	id, err := c.CreateDatasetVersion(
		ctx,
		d.Get("workspace_id").(string),
		d.Get("dataset_id").(string),
		d.Get("contents").(string),
		d.Get("file_name").(string),
		d.Get("has_header").(bool))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%d", d.Get("dataset_id").(string), id))

	return resourceDatasetVersionRead(ctx, d, meta)
}

func resourceDatasetVersionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.TowerClient)

	datasetId, versionId, err := resourceDatasetVersionParseId(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	version, err := c.GetDatasetVersion(
		ctx,
		d.Get("workspace_id").(string),
		datasetId,
		versionId)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("contents", version["contents"].(string))
	d.Set("has_header", version["hasHeader"].(bool))
	d.Set("version", versionId)
	d.Set("last_updated", version["lastUpdated"].(string))
	d.Set("media_type", version["mediaType"].(string))
	d.Set("url", version["url"].(string))

	return nil
}

func resourceDatasetVersionParseId(id string) (string, int, error) {
	parts := strings.Split(id, ":")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", -1, fmt.Errorf("Expected identifier with format: dataset_id:version_id. Got: %v", id)
	}

	versionId, err := strconv.ParseInt(parts[1], 10, 32)

	if err != nil {
		return "", -1, fmt.Errorf("Expected versionId to be an integer, got %v", parts[1])
	}

	return parts[0], int(versionId), nil
}

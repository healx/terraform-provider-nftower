package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenLabels(labels []interface{}) []interface{} {
	res := []interface{}{}
	for _, label := range labels {
		l, _ := label.(map[string]interface{})
		res = append(res, l["name"].(string))
	}

	return res
}

func expandLabels(labels *schema.Set) []string {
	res := []string{}

	for _, l := range labels.List() {
		res = append(res, l.(string))
	}

	return res
}

package lib

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ExpandStringSlice(input []interface{}) *[]string {
	result := make([]string, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, item.(string))
		} else {
			result = append(result, "")
		}
	}
	return &result
}

func ExpandStringMap(m map[string]interface{}) *map[string]string {
	output := make(map[string]string, len(m))
	for i, v := range m {
		output[i] = v.(string)
	}
	return &output
}

func FlattenStringSlice(input *[]string) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}

func FlattenStringMap(m *map[string]string) map[string]interface{} {
	if m == nil {
		return map[string]interface{}{}
	}
	output := make(map[string]interface{}, len(*m))
	for i, v := range *m {
		output[i] = v
	}
	return output
}

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func SetFromStringSlice(slice []string) *schema.Set {
	set := &schema.Set{F: schema.HashString}
	for _, v := range slice {
		set.Add(v)
	}
	return set
}

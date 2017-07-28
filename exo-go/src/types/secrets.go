package types

import "fmt"

type Secrets map[string]string

// converts map to .tfvars string:
// {a:b, c:d} ->
// a="b"
// c="d"
func (s Secrets) ToTfString() TFString {
	tfvars := ""
	for key, value := range s {
		tfvars += fmt.Sprintf("%s=\"%s\"\n", key, value)
	}
	return TFString(tfvars)
}

func (map1 Secrets) HasConflictingKey(map2 Secrets) bool {
	for k := range map1 {
		if _, hasKey := map2[k]; hasKey {
			return true
		}
	}
	return false
}

func (map1 Secrets) MergeSecrets(map2 Secrets) Secrets {
	for k, v := range map1 {
		map2[k] = v
	}
	return map2
}

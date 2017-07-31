package types

import (
	"fmt"
	"sort"
	"strings"
)

// Secrets map contains maps from secret keys to values
type Secrets map[string]string

// NewSecrets creats a Secrets map from a .tfvars string
// a="b"
// c="d" ->
// {a:b, c:d}
func NewSecrets(str string) Secrets {
	secretsMap := Secrets{}
	secretPairs := strings.Split(str, "\n")
	if len(secretPairs) > 1 {
		for _, secret := range secretPairs {
			s := strings.Split(secret, "=")
			secretsMap[s[0]] = strings.Trim(s[1], "\"")
		}
	}
	return secretsMap
}

// TfString converts map to .tfvars string:
// {a:b, c:d} ->
// a="b"
// c="d"
func (s Secrets) TfString() string {
	tfvars := []string{}

	for key, value := range s {
		tfvars = append(tfvars, fmt.Sprintf("%s=\"%s\"", key, value))
	}

	sort.Strings(tfvars)
	return strings.Join(tfvars, "\n")
}

// Merge makes sures two maps. If there are conflicting keys, the value in newSecrets is used
func (s Secrets) Merge(newSecrets Secrets) Secrets {
	for k, v := range newSecrets {
		s[k] = v
	}
	return s
}

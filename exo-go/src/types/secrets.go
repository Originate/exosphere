package types

import (
	"fmt"
	"sort"
	"strings"
)

// Secrets map contains maps from secret keys to values
type Secrets map[string]string

// NewSecrets creats a Secrets map from a .tfvars string
// Input: `a="b"\nc="d"`
// Output: Secrets{"a": "b", "c": "d"}
func NewSecrets(str string) Secrets {
	secretsMap := Secrets{}
	secretPairs := strings.Split(str, "\n")
	for _, secret := range secretPairs {
		s := strings.Split(secret, "=")
		if len(s) > 1 {
			secretsMap[s[0]] = strings.Trim(s[1], "\"")
		}
	}
	return secretsMap
}

// TfString converts map to .tfvars string:
// Input:  Secrets{"a": "b", "c": "d"}
// Output: `a="b"\nc="d"`
func (s Secrets) TfString() string {
	tfvars := []string{}

	for key, value := range s {
		tfvars = append(tfvars, fmt.Sprintf("%s=\"%s\"", key, value))
	}

	sort.Strings(tfvars)
	return strings.Join(tfvars, "\n")
}

// Merge merges two secret maps. It overwrites existing keys of s if conflicting keys exist
func (s Secrets) Merge(newSecrets Secrets) Secrets {
	for k, v := range newSecrets {
		s[k] = v
	}
	return s
}

// DeleteSecrets deletes secrets from s. Ignores them if they do not exist
func (s Secrets) DeleteSecrets(toDelete []string) Secrets {
	for _, key := range toDelete {
		delete(s, key)
	}
	return s
}

// Keys returns all the keys for a secrets map
func (s Secrets) Keys() []string {
	keys := []string{}
	for k := range s {
		keys = append(keys, k)
	}
	return keys
}

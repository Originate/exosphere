package types

import (
	"errors"
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

// HasConflictingKey checks that two secrets map doesn't have clonficting keys
func (s Secrets) HasConflictingKey(map2 Secrets) bool {
	for k := range s {
		if _, hasKey := map2[k]; hasKey {
			return true
		}
	}
	return false
}

// ValidateAndMerge makes sures two maps do not have conflicting keys and merges them
func (s Secrets) ValidateAndMerge(newSecrets Secrets) (Secrets, error) {
	if s.HasConflictingKey(newSecrets) {
		return nil, errors.New("new secrets have key(s) that conflict with existing secrets. Use 'exo configure update' to update existing keys")
	}
	return s.mergeSecrets(newSecrets), nil
}

func (s Secrets) mergeSecrets(map2 Secrets) Secrets {
	for k, v := range s {
		map2[k] = v
	}
	return map2
}

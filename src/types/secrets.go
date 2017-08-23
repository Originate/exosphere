package types

// Secrets map contains maps from secret keys to values
type Secrets map[string]string

// Merge merges two secret maps. It overwrites existing keys of s if conflicting keys exist
func (s Secrets) Merge(newSecrets Secrets) Secrets {
	for k, v := range newSecrets {
		s[k] = v
	}
	return s
}

// Delete deletes secrets from s. Ignores them if they do not exist
func (s Secrets) Delete(toDelete []string) Secrets {
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

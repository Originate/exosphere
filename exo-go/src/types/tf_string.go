package types

import "strings"

// TFString is an abstract of a string representing a .tfvars file
type TFString string

// ToSecrets converts .tfvars string to .tfvars map:
// a="b"
// c="d" ->
// {a:b, c:d}
func (str TFString) ToSecrets() Secrets {
	secretsMap := Secrets{}
	secretPairs := strings.Split(string(str), "\n")
	for _, secret := range secretPairs {
		s := strings.Split(secret, "=")
		secretsMap[s[0]] = strings.Trim(s[1], "\"")
	}
	return secretsMap
}

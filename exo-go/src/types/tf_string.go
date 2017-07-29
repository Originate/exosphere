package types

import "strings"

// TFString is an abstract of a string representing a .tfvars file
type TFString string

// ToMap converts .tfvars string to .tfvars map:
// a="b"
// c="d" ->
// {a:b, c:d}
func (str TFString) ToMap() Secrets {
	secretsMap := Secrets{}
	secretPairs := strings.Split(string(str), "\n")
	for _, secret := range secretPairs {
		s := strings.Split(secret, "=")
		secretsMap[s[0]] = strings.Trim(s[1], "\"")
	}
	return secretsMap
}

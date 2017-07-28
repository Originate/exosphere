package types

import "strings"

type TFString string

// converts .tfvars string to .tfvars map:
// a="b"
// c="d" ->
// {a:b, c:d}
func (str TFString) ToMap() Secrets {
	secretsMap := Secrets{}
	secretPairs := strings.Split(string(str), "\n")
	secretPairs = secretPairs[:len(secretPairs)-1] //remove trailing empty elem
	for _, secret := range secretPairs {
		s := strings.Split(secret, "=")
		secretsMap[s[0]] = strings.Trim(s[1], "\"")
	}
	return secretsMap
}

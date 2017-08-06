package ostools

import "os/user"

// GetHomeDirectory returns the path to the user's home directory
func GetHomeDirectory() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

package main

import "errors"

// @todo rename file and split functions out

func ValidateMageRoot(path string) (bool, error) {
	if path == "" {
		return false, errors.New("please provide the full path to your magento root using --mage-root")
	}

	return true, nil
}

func ValidateDBCredentials(user string, password string, dbName string, host string) (bool, error) {
	requiredDbArgs := []string{user, password, dbName, host}

	for _, arg := range requiredDbArgs {
		if arg == "" {
			return false, errors.New("please provide database credentials see --help")
		}
	}

	return true, nil
}

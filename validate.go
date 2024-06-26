package main

import "errors"

func ValidateMageRoot(path string) (bool, error) {
	if path == "" {
		return false, errors.New("please provide the full path to your magento root using --mage-root")
	}

	return true, nil
}

func ValidateDBCredentials(user string, dbName string, host string) (bool, error) {
	requiredDbArgs := []string{user, dbName, host}

	for _, arg := range requiredDbArgs {
		if arg == "" {
			return false, errors.New("please provide database credentials see --help")
		}
	}

	return true, nil
}

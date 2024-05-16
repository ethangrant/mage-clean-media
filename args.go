package main

import (
	"errors"
	"github.com/manifoldco/promptui"
	"github.com/fatih/color"
)

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

func FullExecutionPrompt(isDryrun bool) (bool) {
	prompt := promptui.Prompt{
		Label:    "Warning: this is not a dry run. If you would like to continue type 'yes'",
	}

	result, err := prompt.Run()
	if err != nil {
		return false
	}

	if result == "yes" {
		return true
	}

	return false
}

func DeleteMessage(isDryRun bool) (string) {
	var deleteMessage string = "DRY-RUN: "

	if !isDryRun {
		deleteMessage = "REMOVING: "
	}

	deleteMessage = color.YellowString(deleteMessage)

	return deleteMessage
}
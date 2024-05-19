package main

import "github.com/manifoldco/promptui"

func FullExecutionPrompt(isDryrun bool) bool {
	prompt := promptui.Prompt{
		Label: "Warning: this is not a dry run. If you would like to continue type 'yes'",
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

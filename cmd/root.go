package cmd

import (
	"fmt"
	"fvm/script"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fvm",
	Short: "fvm provides an interface to manage fec-builder versions.",
	Run: func(cmd *cobra.Command, args []string) {
		currentVersion, _ := script.GetCurrentVersion()

		question := []*survey.Question{
			{
				Name: "Version",
				Prompt: &survey.Select{
					Message: "Choose a version:",
					Options: script.GetLocalVersionList(),
					Default: currentVersion,
				},
			},
		}

		answers := struct {
			Version string
		}{}

		err := survey.Ask(question, &answers)
		if err != nil {
			return
		}

		script.SwitchVersion(answers.Version)
	},
}

func init() {
	rootCmd.AddCommand(createVersionCmd())
	rootCmd.AddCommand(createInstallCmd())
	rootCmd.AddCommand(createUseCmd())
	rootCmd.AddCommand(createListCmd())

	// 隐藏 completion 功能，不要问原因... 没细看是啥作用，所以隐藏了
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

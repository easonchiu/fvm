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
		list := script.GetLocalVersionList()
		currentVersion, _ := script.GetCurrentVersion()

		if len(list) == 0 {
			fmt.Println("")
			fmt.Println("You do not have any version installed\nPlease use \"fvm install [version]\" to install")
			fmt.Println("")
			return
		}

		// 没有安装过时，当前默认指向 list 的第一项
		if len(currentVersion) == 0 {
			currentVersion = list[0]
		}

		question := []*survey.Question{
			{
				Name: "Version",
				Prompt: &survey.Select{
					Message: "Choose a version:",
					Options: list,
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

		err = script.SwitchVersion(answers.Version)
		if err != nil {
			fmt.Println("")
			fmt.Println(err)
			fmt.Println("")
			return
		}

		fmt.Println("")
		fmt.Printf("🎉 Switched to version %v\n", answers.Version)
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(createVersionCmd())
	rootCmd.AddCommand(createInstallCmd())
	rootCmd.AddCommand(createUseCmd())
	rootCmd.AddCommand(createListCmd())
	rootCmd.AddCommand(createRemoveCmd())

	// 隐藏 completion 功能，不要问原因... 没细看是啥作用，所以隐藏了
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

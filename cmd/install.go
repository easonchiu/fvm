package cmd

import (
	"errors"
	"fmt"
	"fvm/script"

	"github.com/spf13/cobra"
)

func createInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install [VERSION]",
		Short:   "Install fec-builder, example: fvm install 2.7.1",
		Example: "fvm install 2.7.1",
		Run: func(cmd *cobra.Command, args []string) {
			err := script.InstallVersion(args[0])

			if err != nil {
				fmt.Println("")
				fmt.Println(err)
				fmt.Println("")
				return
			}

			fmt.Println("")
			fmt.Printf("🎉 Version %v has been installed successfully\n", args[0])
			fmt.Printf("Using [fvm use %v] to switch version %v\n", args[0], args[0])
			fmt.Println("")
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("Version required")
			}
			return nil
		},
	}

	return cmd
}

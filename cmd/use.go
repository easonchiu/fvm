package cmd

import (
	"errors"
	"fmt"
	"fvm/script"

	"github.com/spf13/cobra"
)

func createUseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use [VERSION]",
		Short: "Switch the version of fec-builder, example: fvm use 2.7.1",
		Run: func(cmd *cobra.Command, args []string) {
			err := script.SwitchVersion(args[0])

			if err != nil {
				fmt.Println("")
				fmt.Println(err)
				fmt.Println("")
				return
			}

			fmt.Println("")
			fmt.Printf("🎉 Switched to version %v\n", args[0])
			fmt.Println("")
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("Version required")
			}

			return nil
		},
	}
}

package cmd

import (
	"errors"
	"fmt"
	"fvm/script"

	"github.com/spf13/cobra"
)

func createRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [VERSION]",
		Short: "Remove the version of fec-builder, example: fvm remove 2.7.1",
		Run: func(cmd *cobra.Command, args []string) {
			err := script.RemoveVersion(args[0])

			if err != nil {
				fmt.Println("")
				fmt.Println(err)
				fmt.Println("")
				return
			}

			fmt.Println("")
			fmt.Printf("The version %v has been removed\n", args[0])
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

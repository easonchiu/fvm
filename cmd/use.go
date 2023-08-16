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
		Short: "xxxxxxx",
		Run: func(cmd *cobra.Command, args []string) {
			err := script.SwitchVersion(args[0])

			if err != nil {
				fmt.Println("")
				fmt.Println(err)
				fmt.Println("")
				return
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("version required")
			}

			return nil
		},
	}
}

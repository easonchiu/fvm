package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func createVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Output version of fvm",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("1.0.0")
		},
	}
}

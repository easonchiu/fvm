package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

func createVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Output version of fvm",
		Run: func(cmd *cobra.Command, args []string) {
			bytes, err := os.ReadFile("./package.json")
			if err != nil {
				fmt.Println("1.0.0")
				return
			}
			version := gjson.GetBytes(bytes, "version")
			fmt.Println(version.String())
		},
	}
}

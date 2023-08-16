package cmd

import (
	"fmt"
	"fvm/script"

	"github.com/spf13/cobra"
)

func createListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Output downloaded versions",
		Run: func(cmd *cobra.Command, args []string) {
			list := script.GetLocalVersionList()
			cur, _ := script.GetCurrentVersion()

			// 如果有 cur 信息，但 list 中没有，可能是因为从来没用过该工具，安装一下
			exists := false
			if len(cur) != 0 {
				for _, v := range list {
					if v == cur {
						exists = true
					}
				}
				if !exists {
					_ = script.InstallVersion(cur)
					list = script.GetLocalVersionList()
				}
			}

			// 打印版本列表
			fmt.Println("")
			for _, v := range list {
				if cur == v {
					fmt.Println(">", v)
				} else {
					fmt.Println(" ", v)
				}
			}
			fmt.Println("")
		},
	}
}

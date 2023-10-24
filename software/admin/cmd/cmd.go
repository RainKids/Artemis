package cmd

import (
	"admin/cmd/server"
	"admin/cmd/version"
	"admin/global"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:          "Artemis",
	Long:         "Artemis admin Server",
	SilenceUsage: true,
	Short:        `Artemis admin Server is starting`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New("至少需要一个参数")
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func tip() {
	usageStr := `欢迎使用 ` + `Artemis admin Server` + global.Version + ` 可以使用 ` + `-h` + ` 查看命令`
	fmt.Printf("%s\n", usageStr)
}
func init() {
	rootCmd.AddCommand(version.StartCmd)
	rootCmd.AddCommand(server.StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

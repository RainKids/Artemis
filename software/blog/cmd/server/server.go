package server

import (
	"blog/pkg/application"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	configFile string
	app        *application.Application
	err        error
	StartCmd   = &cobra.Command{
		Use:          "server",
		Short:        "Artemis blog server",
		Example:      "Artemis blog server -c config/setting.yml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			tip()
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func tip() {
	usageStr := `starting MicroCool API server`
	fmt.Printf("%s\n", usageStr)
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config/settings.yaml", "Start server with provided configuration file")
}

func setup() {
	//wrappers.InitFlowQPS()
	//wrappers.InitCircuitBreaker()
	app, err = CreateApp(configFile)
	if err != nil {
		panic(err)
	}
}

func run() error {
	if err = app.Start(); err != nil {
		panic(err)
	}
	app.AwaitSignal()
	return err
}

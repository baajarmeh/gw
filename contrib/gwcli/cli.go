package gwcli

import (
	"github.com/oceanho/gw"
	"github.com/urfave/cli/v2"
)

func App() *cli.App {
	appName := "gw-cli"
	app := &cli.App{
		Name:    appName,
		Usage:   "The gw framework command tools.",
		Version: gw.Version,
		Commands: []*cli.Command{
			{
				Name:      "createapp",
				HelpName:  appName + " createapp <App's Name>",
				Usage:     "create a gw application scaffold.",
				ArgsUsage: "arguments.",
				Action: func(context *cli.Context) error {
					return nil
				},
			},
		},
	}
	return app
}

func startApp() {

}

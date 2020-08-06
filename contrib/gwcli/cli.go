package gwcli

import (
	"github.com/oceanho/gw"
	"github.com/urfave/cli/v2"
)

func App() *cli.App {
	appName := "gw-cli"
	app := &cli.App{
		Name:     appName,
		Usage:    "The gw framework command tools.",
		HelpName: appName,
		Version:  gw.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dir",
				Usage:       "Specifies base directory",
				DefaultText: ".",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "newproject",
				HelpName:  appName + " newproject ",
				Usage:     "create a gw project scaffold",
				ArgsUsage: "<project name>",
				Action: func(context *cli.Context) error {
					return nil
				},
			},
			{
				Name:      "createapp",
				HelpName:  appName + " createapp ",
				Usage:     "create a gw module app scaffold",
				ArgsUsage: "<app name>",
				Action: func(context *cli.Context) error {
					return nil
				},
			},
		},
	}
	return app
}

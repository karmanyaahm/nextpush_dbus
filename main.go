package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"unifiedpush.org/go/nextpush_dbus/cmd"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "lang",
			Value: "english",
			Usage: "Language for the greeting",
		},
		&cli.StringFlag{
			Name:  "instance",
			Value: "",
			Usage: "For multiple instances of the command",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:        "login",
			Usage:       "Log in to Nextcloud",
			Action:      cmd.Login,
			Description: "The first argument is the domain name of the nextcloud server being logged into, for example 'nextcloud.example.com'",
		},
		{
			Name:   "logout",
			Usage:  "logout",
			Action: cmd.Logout,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

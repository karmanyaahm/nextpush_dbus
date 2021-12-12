package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"
	ncauth "k.malhotra.cc/go/nextcloud_authv2/auth"
	"unifiedpush.org/go/nextpush_dbus/api"
	"unifiedpush.org/go/nextpush_dbus/auth"
	"unifiedpush.org/go/nextpush_dbus/config"
	"unifiedpush.org/go/nextpush_dbus/dbusutil"
)

func Login(c *cli.Context) error {
	store, err := auth.NewNCPasswordStorage(config.CliCtx.String("instance"))
	if err != nil {
		return err
	}

	if _, _, _, err := store.RetreivePwd(); err == nil {
		return cli.Exit("Already logged in, log out before logging in again", 0)
	}

	if c.NArg() != 1 {
		cli.ShowSubcommandHelp(c)
		return cli.Exit("Incorrect Arguments", 1)
	}
	server, uname, passwd, err := ncauth.Authenticate(context.TODO(), c.Args().First(), api.UserAgent, c.App.Writer, c.App.Reader)
	if err != nil {
		return err
	}

	err = store.SavePwd(uname, server, passwd)
	fmt.Sprintln(c.App.Writer, "Successfully saved the password")
	return err
}

func Logout(c *cli.Context) error {
	store, err := auth.NewNCPasswordStorage(config.CliCtx.String("instance"))
	if err != nil {
		return err
	}

	//TODO unregister the app password from nextcloud server
	conn, err := dbusutil.GetDBusConn()
	if err != nil {
		return err
	}
	return store.DeletePwd(conn)
}

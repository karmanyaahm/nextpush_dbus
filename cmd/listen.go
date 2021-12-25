package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"zgo.at/errors"

	"github.com/urfave/cli/v2"
	"unifiedpush.org/go/nextpush_dbus/api"
	"unifiedpush.org/go/nextpush_dbus/auth"
	"unifiedpush.org/go/nextpush_dbus/config"
	"unifiedpush.org/go/nextpush_dbus/store"
	"unifiedpush.org/go/np2p_dbus/distributor"
	"unifiedpush.org/go/np2p_dbus/utils"
)

func Listen(ctx *cli.Context) error {
	config.CliCtx = ctx

	store, err := store.InitStorage(utils.StoragePath("nextpush.db"))
	if err != nil {
		return err
	}

	server, _, _, err := auth.GetCreds()
	if err != nil {
		return errors.Wrap(err, "Are you logged in?")
	}
	fmt.Println("Connecting to", server)

	dbus := distributor.NewDBus("org.unifiedpush.Distributor.NextPush")
	dbus.StartHandling(handler{dbus: dbus, store: store})

	//fmt.Println(server + api.NCAppPathPrefix + "/push/" + token)

	id := store.GetDeviceID()
	if id == nil || *id == "" {
		id = createDevice(ctx, store)
	}

	messages := make(chan api.Message)

	go sendMsgtoApp(store, messages, dbus)

	for {
		utils.Log.Infoln("Syncing...")
		fmt.Println(*id)
		err = api.Sync(*id, messages)
		if err != nil {
			fmt.Println("SYNCERR " + err.Error())
			time.Sleep(10 * time.Second)
		}
	}

	//err = api.DeleteDevice(*id)
	//if err != nil {
	//	log.Fatal(err)
	//}
	return nil
}

func sendMsgtoApp(st store.Storage, messages chan api.Message, dbus *distributor.DBus) {
	for {
		m := <-messages
		utils.Log.Debugln(m.Token, m.Message)
		conn := st.GetConnectionbyPublic(m.Token)

		_ = dbus.NewConnector(conn.AppID).Message(conn.AppToken, string(m.Message), "")

	}
}

func createDevice(ctx *cli.Context, store store.Storage) (id *string) {
	idnew, err := api.CreateDevice("mydev")
	utils.Log.Debugln("Registering new device")
	if err != nil {
		log.Fatal(err)
	}

	err = store.SetDeviceID(idnew)
	if err != nil {
		log.Fatal(err)
	}
	return &idnew
}

func checkLoggedIn(ctx *cli.Context) {
	_, _, _, err := auth.GetCreds()
	if err != nil {
		log.Fatal(err)
	}
}

type handler struct {
	store store.Storage
	dbus  *distributor.DBus
}

func (h handler) Register(appName, token string) (endpoint, refuseReason string, err error) {
	utils.Log.Debugln(appName, "registration request")

	existing := h.store.GetConnectionbyPrivate(token)
	if existing != nil {
		return api.GetEndpointFromApp(existing.PublicToken), "", nil
	}
	devID := h.store.GetDeviceID()
	if devID == nil {
		return "", "", errors.New("NextPush distributor Not logged in")
	}

	pubtoken, err := api.CreateApp(*devID, appName)
	if err != nil {
		ret := errors.Wrap(err, "Nextcloud error")
		utils.Log.Infoln(ret)
		return "", "", rawErr(ret)
	}

	settings := api.GetEndpointFromApp("<token>")
	conn := h.store.NewConnectionFull(appName, token, pubtoken, settings)
	utils.Log.Debugln("registered new", conn.AppID, conn.Settings)
	if conn != nil {
		return api.GetEndpointFromApp(conn.PublicToken), "", nil
	}

	return "", "", errors.New("Unknown error with NextPush")
}
func (h handler) Unregister(token string) {
	deletedConn, err := h.store.DeleteConnection(token)
	if err != nil {
		//log
	}
	err = api.DeleteApp(deletedConn.PublicToken)
	if err != nil {
		//log
	}
	utils.Log.Debugln("deleted", deletedConn)

	_ = h.dbus.NewConnector(deletedConn.AppID).Unregistered(deletedConn.AppToken)
}

func rawErr(e error) error {
	ans := strings.SplitN(e.Error(), "\n", 2)
	return fmt.Errorf(ans[0])
}

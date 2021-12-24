package auth

import (
	"errors"

	"github.com/godbus/dbus/v5"
	keyring "github.com/ppacher/go-dbus-keyring"
	"unifiedpush.org/go/nextpush_dbus/config"
	"unifiedpush.org/go/nextpush_dbus/dbusutil"
)

var (
	errCorruptedSecret = errors.New("Corrupted Password in Keyring")
)

//caller needs to close session
func initConn(bus *dbus.Conn) (session keyring.Session, collection keyring.Collection, err error) {
	secrets, err := keyring.GetSecretService(bus)
	if err != nil {
		return
	}
	session, err = secrets.OpenSession()
	if err != nil {
		return
	}
	collection, err = secrets.GetDefaultCollection()

	return
}

func NewNCPasswordStorage(instance string) (storage ncPasswordStorage, err error) {
	conn, err := dbusutil.GetDBusConn()
	if err != nil {
		return
	}
	storage = ncPasswordStorage{keyringPasswordLabel: "Nextcloud-UP-" + instance, dbusConn: conn}
	return
}

type ncPasswordStorage struct {
	keyringPasswordLabel string
	dbusConn             *dbus.Conn
}

func (store ncPasswordStorage) SavePwd(user, server, passwd string) (err error) {
	session, collection, err := initConn(store.dbusConn)
	if err != nil {
		return
	}
	defer session.Close()

	attrs := map[string]string{"server": "Nextcloud", "user": user, "serverURL": server}
	_, err = collection.CreateItem(session.Path(), store.keyringPasswordLabel, attrs, []byte(passwd), "text/plain", true)

	return
}

func (store ncPasswordStorage) RetreivePwd() (server, user, passwd string, err error) {
	session, collection, err := initConn(store.dbusConn)
	if err != nil {
		return
	}
	defer session.Close()

	item, err := collection.GetItem(store.keyringPasswordLabel)
	if err != nil {
		return
	}
	attrs, err := item.GetAttributes()
	if err != nil {
		return
	}
	secret, err := item.GetSecret(session.Path())
	if err != nil {
		return
	}

	user, ok := attrs["user"]
	if !ok {
		err = errCorruptedSecret
		return
	}
	server, ok = attrs["serverURL"]
	if !ok {
		err = errCorruptedSecret
		return
	}

	passwd = string(secret.Value)
	return
}

func (store ncPasswordStorage) DeletePwd(bus *dbus.Conn) (err error) {
	session, collection, err := initConn(bus)
	if err != nil {
		return
	}
	defer session.Close()

	item, err := collection.GetItem(store.keyringPasswordLabel)
	if err != nil {
		return
	}

	err = item.Delete()

	return
}
func GetCreds() (server, username, password string, err error) {
	//TODO cache credentials so each time it doesn't call to dbus

	store, err := NewNCPasswordStorage(config.CliCtx.String("instance"))
	if err != nil {
		return
	}

	return store.RetreivePwd()
}

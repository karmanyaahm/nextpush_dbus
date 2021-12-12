package dbusutil

import "github.com/godbus/dbus/v5"

var conn *dbus.Conn

func GetDBusConn() (*dbus.Conn, error) {
	if conn != nil {
		return conn, nil
	}

	conn, err := dbus.SessionBus()
	return conn, err
}

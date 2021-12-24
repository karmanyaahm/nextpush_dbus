module unifiedpush.org/go/nextpush_dbus

go 1.17

require (
	github.com/godbus/dbus/v5 v5.0.6
	github.com/ppacher/go-dbus-keyring v1.0.1
	github.com/urfave/cli/v2 v2.3.0
	k.malhotra.cc/go/nextcloud_authv2/auth v0.0.0-20211212060846-d87709108551
)

require (
	astuart.co/go-sse v1.0.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/mattn/go-sqlite3 v1.14.8 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	gorm.io/driver/sqlite v1.1.5 // indirect
	gorm.io/gorm v1.21.15 // indirect
	unifiedpush.org/go/np2p_dbus v0.0.0-20210917013344-e7eac6892e24 // indirect
	zgo.at/errors v1.1.0 // indirect
)

replace unifiedpush.org/go/np2p_dbus => ../np2p_linux

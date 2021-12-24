package store

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"unifiedpush.org/go/np2p_dbus/storage"
	"unifiedpush.org/go/np2p_dbus/utils"
)

func InitStorage(filepath string) (Storage, error) {
	s, e := storage.InitStorage(filepath)

	s.DB().AutoMigrate(&KVStore{})
	//s.DB().Config.Logger = logger.Default.LogMode(4)
	return Storage{*s}, e
}

type Storage struct {
	storage.Storage
}

type KVStore struct {
	Key   string `gorm:"primaryKey"`
	Value string
}

func (kv *KVStore) Get(db *gorm.DB) error {
	return db.First(kv).Error
}

func (kv KVStore) Set(db *gorm.DB) error {
	return db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&kv).Error
}

func (s Storage) GetDeviceID() (ans *string) {
	answer := KVStore{Key: "device-id"}
	if err := answer.Get(s.DB()); err != nil {
		utils.Log.Debugln(err)
		return nil
	}
	return &answer.Value
}

func (s Storage) SetDeviceID(id string) error {
	return KVStore{"device-id", id}.Set(s.DB())
}

func (s Storage) GetConnectionbyPrivate(privateToken string) *storage.Connection {
	c := storage.Connection{AppToken: privateToken}
	return s.GetFirst(c)
}

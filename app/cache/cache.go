package cache

import (
	"bytes"
	"encoding/gob"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/config"
)

var conn *memcache.Client

func getMc() *memcache.Client {
	if conn == nil {
		conf := config.GetMemcacheConfig()
		conn = memcache.New(conf.GetConnectionString())
	}
	return conn
}

func IsMissError(err error) bool {
	return jerr.HasError(err, "memcache: cache miss")
}

func SetItem(name string, value interface{}) error {
	return SetItemWithExpiration(name, value, 0)
}

func SetItemWithExpiration(name string, value interface{}, expireSeconds int32) error {
	writer := new(bytes.Buffer)
	encoder := gob.NewEncoder(writer)
	encoder.Encode(value)
	mc := getMc()
	err := mc.Set(&memcache.Item{
		Key:        name,
		Value:      writer.Bytes(),
		Expiration: expireSeconds,
	})
	if err != nil {
		return jerr.Get("error writing memcache item", err)
	}
	return nil
}

func GetItem(name string, value interface{}) error {
	mc := getMc()
	it, err := mc.Get(name)
	if err != nil {
		return jerr.Get("error getting memcache item", err)
	}
	reader := bytes.NewReader(it.Value)
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(value)
	if err != nil {
		return jerr.Get("error decoding value", err)
	}
	return nil
}

func DeleteItem(name string) error {
	mc := getMc()
	err := mc.Delete(name)
	if err != nil {
		return jerr.Get("error deleting memcache item", err)
	}
	return nil
}

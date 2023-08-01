package config

import (
	"sync"
)

type CacheInterface interface {
	Get(string)
	GetAll()
	Set(string, interface{})
	Delete(string)
}

type Cache struct {
	syncMap sync.Map
	CacheInterface
}

func (c *Cache) Get(key string) (interface{}, bool) {
	rawVal, ok := c.syncMap.Load(key)
	if !ok {
		return nil, false
	}

	return rawVal, true
}

func (c *Cache) GetAll() (interface{}) {
	var values = make(map[string]interface{})

	c.syncMap.Range(func(rawKey, rawVal interface{}) bool {
		k, ok := rawKey.(string)

		if !ok {
			return false
		}

		values[k] = rawVal
		return true
	})

	return values
}


func (c *Cache) Set(key string, value interface{}) (bool) {
	c.syncMap.Store(key, value)

	return true
}

func (c *Cache) Delete(key string) (bool) {
	c.syncMap.Delete(key)

	return true
}

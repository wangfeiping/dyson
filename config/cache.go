package config

import (
	"sync"
)

var cmux sync.RWMutex
var cache *Cache

func init() {
	cache = &Cache{
		mapper: make(map[string]string)}
}

type Cache struct {
	mapper map[string]string
}

func GetCache() *Cache {
	cmux.Lock()
	defer cmux.Unlock()

	return cache
}

func (c *Cache) Get(name string) string {
	return c.mapper[name]
}

func (c *Cache) Put(name string, value string) {
	c.mapper[name] = value
}

func (c *Cache) Clear() {
	cmux.Lock()
	defer cmux.Unlock()

	c.mapper = make(map[string]string)
}

package cache

import "time"

type item struct {
	value      string
	expiration time.Time
}

func (i *item) isExpired() bool {
	if i.expiration.IsZero() {
		return false
	}
	now := time.Now()
	return i.expiration.Before(now)
}

type Cache struct {
	items map[string]item
}

func NewCache() Cache {
	return Cache{items: map[string]item{}}
}

func (c Cache) Get(key string) (string, bool) {
	i, ok := c.items[key]
	if !ok {
		return "", false
	}
	if i.isExpired() {
		delete(c.items, key)
		return "", false
	}
	return i.value, true
}

func (c Cache) Put(key, value string) {
	i := item{value: value}
	c.items[key] = i
}

func (c Cache) Keys() []string {
	result := make([]string, 0, len(c.items))

	for k := range c.items {
		_, ok := c.Get(k)
		if ok {
			result = append(result, k)
		}
	}

	return result
}

func (c Cache) PutTill(key, value string, deadline time.Time) {
	i := item{value: value, expiration: deadline}
	c.items[key] = i
}

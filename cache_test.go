package cache

import (
	"testing"
	"time"
)

type kvPair struct {
	key          string
	value        string
	expValue     string
	deadline     time.Time
	shouldExpire bool
}

func TestCache(t *testing.T) {
	cases := map[string]struct {
		kvPairs []kvPair
		expKeys []string
		waitFor time.Duration
	}{
		"empty": {
			kvPairs: []kvPair{},
			expKeys: []string{},
		},
		"one value": {
			kvPairs: []kvPair{{
				key:      "key1",
				value:    "value1",
				expValue: "value1",
			}},
			expKeys: []string{"key1"},
		},
		"several values": {
			kvPairs: []kvPair{
				{
					key:      "key1",
					value:    "value1",
					expValue: "value1",
				},
				{
					key:      "key2",
					value:    "value2",
					expValue: "value2",
				},
				{
					key:      "key3",
					value:    "value3",
					expValue: "value3",
				},
			},
			expKeys: []string{"key1", "key2", "key3"},
		},
		"overwrite a value": {
			kvPairs: []kvPair{
				{
					key:      "key1",
					value:    "value1",
					expValue: "anotherValue1",
				},
				{
					key:      "key2",
					value:    "value2",
					expValue: "value2",
				},
				{
					key:      "key3",
					value:    "value3",
					expValue: "value3",
				},
				{
					key:      "key1",
					value:    "anotherValue1",
					expValue: "anotherValue1",
				},
			},
			expKeys: []string{"key1", "key2", "key3"},
		},
		"expired values": {
			kvPairs: []kvPair{
				{
					key:      "key1",
					value:    "value1",
					expValue: "value1",
				},
				{
					key:          "key2",
					value:        "value2",
					expValue:     "value2",
					shouldExpire: true,
					deadline:     time.Now().Add(time.Second * 2),
				},
				{
					key:      "key3",
					value:    "value3",
					expValue: "value3",
					deadline: time.Now().Add(time.Minute * 2),
				},
			},
			expKeys: []string{"key1", "key3"},
			waitFor: time.Second * 3,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			c := NewCache()
			for _, p := range tt.kvPairs {
				if p.deadline.IsZero() {
					c.Put(p.key, p.value)
				} else {
					c.PutTill(p.key, p.value, p.deadline)
				}
			}
			time.Sleep(tt.waitFor)
			for _, p := range tt.kvPairs {
				v, ok := c.Get(p.key)
				if !p.shouldExpire {
					if !ok {
						t.Errorf("Get: value is not present, while it should")
					}

					if v == "" {
						t.Errorf("Get: returned value is empty, while should be set")
					}

					if v != p.expValue {
						t.Errorf("Get: returned value incorrect: want \"%s\", got \"%s\"", p.expValue, v)
					}
				} else {
					if ok {
						t.Errorf("Get: and expired value is present in the cache(returned ok==true), while it should't")
					}

					if v != "" {
						t.Errorf("Get: and expired value is not an empty value: \"%s\"", v)
					}
				}
			}

			v, ok := c.Get("notExistingKey")
			if ok {
				t.Errorf("Get: random key is present in the cache(returned ok==true), while it should't")
			}

			if v != "" {
				t.Errorf("Get: random key is not an empty value: \"%s\"", v)
			}

			keys := c.Keys()
			if len(tt.expKeys) != len(keys) {
				t.Errorf("Keys: number of returned keys is incorrect: exp: %d, got %d", len(tt.expKeys), len(keys))
			}
			for _, expKey := range tt.expKeys {
				exists := false
				for _, key := range keys {
					if expKey == key {
						exists = true
					}
				}
				if !exists {
					t.Errorf("Keys: a key \"%s\" is not present in the Keys() method output", expKey)
				}
			}
		})
	}
}

func TestCache_OverwriteValue(t *testing.T) {
	c := NewCache()

	c.Put("key1", "value1")
	c.Put("key1", "value2")

	v, ok := c.Get("key1")
	if !ok {
		t.Errorf("Get: value is not present, while it should")
	}

	if v == "" {
		t.Errorf("Get: returned value is empty, while should be set")
	}

	if v != "value2" {
		t.Errorf("Get: returned value incorrect: want \"%s\", got \"%s\"", "value2", v)
	}

}

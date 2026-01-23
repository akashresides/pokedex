package pokecache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second

	cases := []struct {
		inputKey   string
		inputValue []byte
	}{
		{
			inputKey:   "https://example.com",
			inputValue: []byte("testdata"),
		},
		{
			inputKey:   "https://example.com/path",
			inputValue: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		cache := NewCache(interval)
		cache.Add(c.inputKey, c.inputValue)

		val, ok := cache.Get(c.inputKey)
		if !ok {
			t.Errorf("expected to find key %s", c.inputKey)
			continue
		}

		if string(val) != string(c.inputValue) {
			t.Errorf("case %d: expected value %s, got %s", i, string(c.inputValue), string(val))
		}

		if len(cache.entries) != 1 {
			t.Errorf("case %d: expected 1 entry in cache, got %d", i, len(cache.entries))
		}
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond

	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key after reaping")
	}
}

func TestReapLoopNotExpired(t *testing.T) {
	const baseTime = 10 * time.Millisecond
	const waitTime = baseTime / 2

	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	time.Sleep(waitTime)

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key before reaping")
	}
}

func TestAddOverwrite(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)

	key := "https://example.com"
	cache.Add(key, []byte("original"))
	cache.Add(key, []byte("updated"))

	val, ok := cache.Get(key)
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	if string(val) != "updated" {
		t.Errorf("expected value 'updated', got '%s'", string(val))
	}
}

func TestGetNonExistent(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)

	_, ok := cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find non-existent key")
	}
}

package main

import (
	"sync"
)

type escapeMetricNameCache struct {
	lock    sync.RWMutex
	size    int
	entries map[string]string
}

func newEscapeMetricNameCache(size int) *escapeMetricNameCache {
	if size == 0 {
		return nil
	}

	return &escapeMetricNameCache{
		size:    size,
		entries: make(map[string]string, size+1),
	}
}

func (e *escapeMetricNameCache) get(key string) (string, bool) {
	if e == nil {
		return "", false
	}

	e.lock.RLock()
	v, ok := e.entries[key]
	e.lock.RUnlock()

	escapeCacheGets.Inc()
	if ok {
		escapeCacheHits.Inc()
	}

	return v, ok
}

func (e *escapeMetricNameCache) set(key string, value string) {
	if e == nil {
		return
	}

	var count int

	e.lock.Lock()

	e.entries[key] = value

	// evict if needed
	if len(e.entries) > e.size {
		for k := range e.entries {
			delete(e.entries, k)
			break
		}
	}

	count = len(e.entries)

	e.lock.Unlock()

	escapeCacheItems.Set(float64(count))
}

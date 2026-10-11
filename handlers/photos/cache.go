package photos

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

const (
	versionTTL   = 5 * time.Second // how stale a cached answer can be after a sync, across instances
	entryTTL     = 10 * time.Minute
	maxCacheSize = 500
)

type cacheEntry struct {
	body    []byte
	expires time.Time
}

// responseCache keeps serialized list responses in memory. Entries are keyed by the index's data_version, so a
// completed sync invalidates every instance without any shared store: this is the job Redis would otherwise do.
type responseCache struct {
	version func(ctx context.Context) (int64, error)
	now     func() time.Time

	hits, builds, coalesced atomic.Int64

	mu      sync.Mutex
	entries map[string]cacheEntry
	ver     int64
	verAt   time.Time
	sf      singleflight.Group
}

func newResponseCache(version func(ctx context.Context) (int64, error)) *responseCache {
	return &responseCache{version: version, now: time.Now, entries: map[string]cacheEntry{}}
}

func (c *responseCache) currentVersion(ctx context.Context) (int64, error) {
	c.mu.Lock()
	if !c.verAt.IsZero() && c.now().Sub(c.verAt) < versionTTL {
		v := c.ver
		c.mu.Unlock()
		return v, nil
	}
	c.mu.Unlock()
	v, err, _ := c.sf.Do("version", func() (any, error) { return c.version(ctx) })
	if err != nil {
		return 0, err
	}
	c.mu.Lock()
	c.ver, c.verAt = v.(int64), c.now()
	c.mu.Unlock()
	return v.(int64), nil
}

// get returns the cached body for key, building it once even under concurrent identical requests.
func (c *responseCache) get(ctx context.Context, key string, build func() ([]byte, error)) ([]byte, error) {
	ver, err := c.currentVersion(ctx)
	if err != nil {
		return nil, err
	}
	full := versionedKey(ver, key)

	c.mu.Lock()
	if e, ok := c.entries[full]; ok && c.now().Before(e.expires) {
		c.mu.Unlock()
		c.hits.Add(1)
		return e.body, nil
	}
	c.mu.Unlock()

	built := false
	body, err, _ := c.sf.Do(full, func() (any, error) {
		built = true
		c.builds.Add(1)
		b, err := build()
		if err != nil {
			return nil, err
		}
		c.put(full, b)
		return b, nil
	})
	if err != nil {
		return nil, err
	}
	if !built {
		c.coalesced.Add(1)
	}
	return body.([]byte), nil
}

// CacheStats counts how requests were answered: from memory, by building from the database, or by waiting for a
// build another request had already started.
type CacheStats struct {
	Hits      int64   `json:"hits"`
	Builds    int64   `json:"builds"`
	Coalesced int64   `json:"coalesced"`
	Entries   int     `json:"entries"`
	HitRatio  float64 `json:"hit_ratio"`
}

func (c *responseCache) Stats() CacheStats {
	c.mu.Lock()
	entries := len(c.entries)
	c.mu.Unlock()
	st := CacheStats{Hits: c.hits.Load(), Builds: c.builds.Load(), Coalesced: c.coalesced.Load(), Entries: entries}
	if total := st.Hits + st.Builds + st.Coalesced; total > 0 {
		st.HitRatio = float64(st.Hits+st.Coalesced) / float64(total)
	}
	return st
}

func (c *responseCache) put(key string, body []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) >= maxCacheSize {
		now := c.now()
		for k, e := range c.entries {
			if now.After(e.expires) {
				delete(c.entries, k)
			}
		}
		if len(c.entries) >= maxCacheSize {
			c.entries = map[string]cacheEntry{} // older versions are unreachable anyway
		}
	}
	c.entries[key] = cacheEntry{body: body, expires: c.now().Add(entryTTL)}
}

func versionedKey(version int64, key string) string {
	return strconv.FormatInt(version, 10) + "|" + key
}

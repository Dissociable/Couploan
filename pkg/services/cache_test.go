package services

import (
	"context"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheClient(t *testing.T) {
	type cacheTest struct {
		Value string
	}
	// Cache some data
	data := cacheTest{Value: "abcdef"}
	group := "testgroup"
	key := "testkey"
	err := c.Cache.
		Set().
		Group(group).
		Key(key).
		Data(data).
		Save(context.Background())
	require.NoError(t, err)

	// Get the data
	fromCache, err := c.Cache.
		Get().
		Group(group).
		Key(key).
		Type(new(cacheTest)).
		Fetch(context.Background())
	require.NoError(t, err)
	cast, ok := fromCache.(*cacheTest)
	require.True(t, ok)
	assert.Equal(t, data, *cast)

	// The same key with the wrong group should fail
	_, err = c.Cache.
		Get().
		Key(key).
		Type(new(cacheTest)).
		Fetch(context.Background())
	assert.Error(t, err)

	// Flush the data
	err = c.Cache.
		Flush().
		Group(group).
		Key(key).
		Execute(context.Background())
	require.NoError(t, err)

	// The data should be gone
	assertFlushed := func() {
		// The data should be gone
		_, err = c.Cache.
			Get().
			Group(group).
			Key(key).
			Type(new(cacheTest)).
			Fetch(context.Background())
		assert.Equal(t, store.NotFoundWithCause(redis.Nil), err)
	}
	assertFlushed()

	// Set with tags
	err = c.Cache.
		Set().
		Group(group).
		Key(key).
		Data(data).
		Tags("tag1").
		Save(context.Background())
	require.NoError(t, err)

	// Flush the tag
	err = c.Cache.
		Flush().
		Tags("tag1").
		Execute(context.Background())
	require.NoError(t, err)

	// The data should be gone
	assertFlushed()

	// Set with expiration
	err = c.Cache.
		Set().
		Group(group).
		Key(key).
		Data(data).
		Expiration(time.Millisecond).
		Save(context.Background())
	require.NoError(t, err)

	// Wait for expiration
	time.Sleep(time.Millisecond * 2)

	// The data should be gone
	assertFlushed()
}

package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(1)

		key := Key("keyOne")
		val := 1
		isExistInCache := c.Set(key, val)
		require.False(t, isExistInCache)

		valueFromCache, isExistInCache := c.Get(key)
		require.True(t, isExistInCache)
		require.Equal(t, val, valueFromCache)

		c.Clear()
		valueFromCache, isExistInCache = c.Get(key)
		require.False(t, isExistInCache)
		require.Nil(t, valueFromCache)
	})

	t.Run("Write a key twice and read it", func(t *testing.T) {
		c := NewCache(25)

		isExistInCache := c.Set("key1", 1)
		require.False(t, isExistInCache)

		isExistInCache = c.Set("key1", 2)
		require.True(t, isExistInCache)

		valueFromCache, isExistInCache := c.Get("key1")
		require.True(t, isExistInCache)
		require.Equal(t, 2, valueFromCache)
	})

	t.Run("knockout first element", func(t *testing.T) {
		c := NewCache(3) // cap 3

		isExistInCache := c.Set("key1", 1)
		require.False(t, isExistInCache)

		isExistInCache = c.Set("key2", 2)
		require.False(t, isExistInCache)

		isExistInCache = c.Set("key3", 3)
		require.False(t, isExistInCache)

		isExistInCache = c.Set("key4", 4)
		require.False(t, isExistInCache)

		valueFromCache, isExistInCache := c.Get("key1")
		require.False(t, isExistInCache)
		require.Nil(t, valueFromCache)

		valueFromCache, isExistInCache = c.Get("key2")
		require.True(t, isExistInCache)
		require.Equal(t, 2, valueFromCache)

		valueFromCache, isExistInCache = c.Get("key3")
		require.True(t, isExistInCache)
		require.Equal(t, 3, valueFromCache)

		valueFromCache, isExistInCache = c.Get("key4")
		require.True(t, isExistInCache)
		require.Equal(t, 4, valueFromCache)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

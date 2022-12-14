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
		c := NewCache[string, int](10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache[string, int](5)

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
		require.Empty(t, val)
	})

	t.Run("cache repeat key", func(t *testing.T) {
		capacity := 3

		c := NewCache[string, int](capacity)
		cs := c.(*lruCache[string, int])

		c.Set("aaa", 100)
		c.Set("aaa", 200)
		c.Set("bbb", 300)
		c.Set("ccc", 400)

		require.Equal(t, capacity, cs.queue.Len())
		require.Equal(t, capacity, len(cs.items))
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache[string, int](2)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Empty(t, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 300, val)
	})

	t.Run("purge logic considering permutation", func(t *testing.T) {
		c := NewCache[string, int](2)

		c.Set("aaa", 100)
		c.Set("bbb", 200)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)
		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		c.Set("ccc", 300)
		val, ok = c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 300, val)
		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)
		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Empty(t, val)
	})

	t.Run("clear", func(t *testing.T) {
		c := NewCache[string, int](2)
		c.Set("aaa", 100)
		c.Clear()

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Empty(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	// t.Skip()

	c := NewCache[string, int](10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(strconv.Itoa(i), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(strconv.Itoa(rand.Intn(1_000_000)))
		}
	}()

	wg.Wait()
}

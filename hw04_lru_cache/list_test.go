package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("check nil", func(t *testing.T) {
		list := NewList()

		list.PushFront(10)
		list.PushBack(20)

		list.Remove(nil)
		list.MoveToFront(nil)
	})

	t.Run("MoveToFront", func(t *testing.T) {
		list := NewList()

		item1 := list.PushFront(10)
		item2 := list.PushBack(20)

		list.MoveToFront(item1)
		require.Equal(t, item1, list.Front())

		list.MoveToFront(item2)
		require.Equal(t, item2, list.Front())
	})

	t.Run("Remove", func(t *testing.T) {
		list := NewList()

		item1 := list.PushFront(10)
		item2 := list.PushBack(20)
		item3 := list.PushBack(30)

		list.Remove(item2)
		list.Remove(item2)
		require.Equal(t, item1, list.Front())
		require.Equal(t, item3, list.Back())
		require.Equal(t, 2, list.Len())

		list.Remove(item1)
		list.Remove(item1)
		require.Equal(t, item3, list.Front())
		require.Equal(t, item3, list.Back())
		require.Equal(t, 1, list.Len())
	})

	t.Run("PushBack", func(t *testing.T) {
		list := NewList()

		a, b := 10, 20

		list.PushBack(a)
		require.Equal(t, a, list.Front().Value)
		require.Equal(t, a, list.Back().Value)
		require.Equal(t, 1, list.Len())

		list.PushBack(b)
		require.Equal(t, a, list.Front().Value)
		require.Equal(t, b, list.Back().Value)
		require.Equal(t, 2, list.Len())
	})

	t.Run("PushFront", func(t *testing.T) {
		list := NewList()

		a, b := 10, 20

		list.PushFront(a)
		require.Equal(t, a, list.Front().Value)
		require.Equal(t, a, list.Back().Value)
		require.Equal(t, 1, list.Len())

		list.PushFront(b)
		require.Equal(t, b, list.Front().Value)
		require.Equal(t, a, list.Back().Value)
		require.Equal(t, 2, list.Len())
	})
}

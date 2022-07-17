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

	t.Run("Start push from front", func(t *testing.T) {
		l := NewList()

		l.PushFront(1) // [1]
		l.PushFront(2) // [2, 1]
		l.PushBack(3)  // [2, 1, 3]
		l.PushBack(4)  // [2, 1, 3, 4]

		require.Equal(t, 4, l.Len())
		require.Equal(t, 2, l.Front().Value)
		require.Equal(t, 1, l.Front().Next.Value)
		require.Equal(t, 3, l.Front().Next.Next.Value)
		require.Equal(t, 4, l.Front().Next.Next.Next.Value)
		require.Equal(t, 4, l.Back().Value)
		require.Equal(t, 3, l.Back().Prev.Value)
		require.Equal(t, 1, l.Back().Prev.Prev.Value)
		require.Equal(t, 2, l.Back().Prev.Prev.Prev.Value)
	})

	t.Run("Start push from back", func(t *testing.T) {
		l := NewList()

		l.PushBack(1)  // [1]
		l.PushBack(2)  // [1, 2]
		l.PushFront(3) // [3, 1, 2]
		l.PushFront(4) // [4, 3, 1, 2]

		require.Equal(t, 4, l.Len())
		require.Equal(t, 4, l.Front().Value)
		require.Equal(t, 3, l.Front().Next.Value)
		require.Equal(t, 1, l.Front().Next.Next.Value)
		require.Equal(t, 2, l.Front().Next.Next.Next.Value)
		require.Equal(t, 2, l.Back().Value)
		require.Equal(t, 1, l.Back().Prev.Value)
		require.Equal(t, 3, l.Back().Prev.Prev.Value)
		require.Equal(t, 4, l.Back().Prev.Prev.Prev.Value)
	})

	t.Run("Remove exist item", func(t *testing.T) {
		l := NewList()

		l.PushFront(1) // [1]

		require.Equal(t, 1, l.Len())

		require.True(t, l.Remove(l.Front()))
		require.Nil(t, l.Front())
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("Remove not exist item", func(t *testing.T) {
		l := NewList()
		require.False(t, l.Remove(l.Front()))
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
}

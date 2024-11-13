package hw04lrucache

import (
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

//nolint:funlen
func TestList(t *testing.T) {
	t.Run("push front", func(t *testing.T) {
		l := NewList()
		l.PushFront(1)

		require.NotNil(t, l.Front())
		require.NotNil(t, l.Back())

		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Front().Next)

		require.Nil(t, l.Back().Prev)
		require.Nil(t, l.Back().Next)

		require.Equal(t, 1, l.Len())
		require.Equal(t, l.Front(), l.Back())
		require.Equal(t, l.Front().Value, 1)
	})

	t.Run("push back", func(t *testing.T) {
		l := NewList()
		l.PushBack(1)

		require.NotNil(t, l.Front())
		require.NotNil(t, l.Back())

		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Front().Next)

		require.Nil(t, l.Back().Prev)
		require.Nil(t, l.Back().Next)

		require.Equal(t, 1, l.Len())
		require.Equal(t, l.Front(), l.Back())
		require.Equal(t, l.Front().Value, 1)
	})

	t.Run("push front back", func(t *testing.T) {
		l := NewList()
		l.PushFront(1)
		l.PushBack(2)

		require.Equal(t, 2, l.Len())
		require.Equal(t, l.Front().Next, l.Back())
		require.Equal(t, l.Front(), l.Back().Prev)
	})

	t.Run("remove front", func(t *testing.T) {
		l := NewList()
		l.PushFront(1)
		l.PushFront(2)
		l.PushFront(3)

		require.Equal(t, 3, l.Len())
		require.Equal(t, l.Front().Value, 3)

		l.Remove(l.Front())
		require.Equal(t, 2, l.Len())
		require.Equal(t, l.Front().Value, 2)

		l.Remove(l.Front())
		require.Equal(t, 1, l.Len())
		require.Equal(t, l.Front().Value, 1)

		l.Remove(l.Front())
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("remove back", func(t *testing.T) {
		l := NewList()
		l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)

		require.Equal(t, 3, l.Len())
		require.Equal(t, l.Back().Value, 3)

		l.Remove(l.Back())
		require.Equal(t, 2, l.Len())
		require.Equal(t, l.Back().Value, 2)

		l.Remove(l.Back())
		require.Equal(t, 1, l.Len())
		require.Equal(t, l.Back().Value, 1)

		l.Remove(l.Back())
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("remove middle", func(t *testing.T) {
		l := NewList()
		l.PushFront(3)
		l.PushFront(2)
		l.PushFront(1)

		require.Equal(t, 3, l.Len())
		require.Equal(t, l.Back().Value, 3)
		require.Equal(t, l.Front().Value, 1)

		l.Remove(l.Front().Next)
		require.Equal(t, 2, l.Len())
		require.Equal(t, l.Front().Value, 1)
		require.Equal(t, l.Back().Value, 3)
		require.Equal(t, l.Front().Next, l.Back())
		require.Equal(t, l.Front(), l.Back().Prev)
	})

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
}

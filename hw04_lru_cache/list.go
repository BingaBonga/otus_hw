package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.len++
	newFront := &ListItem{Value: v}
	currentFront := l.front
	l.front = newFront

	if currentFront != nil {
		newFront.Next = currentFront
		currentFront.Prev = newFront
	}

	if l.back == nil {
		l.back = l.front
	}

	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.len++
	newBack := &ListItem{Value: v}
	currentBack := l.back
	l.back = newBack

	if currentBack != nil {
		newBack.Prev = currentBack
		currentBack.Next = newBack
	}

	if l.front == nil {
		l.front = newBack
	}

	return l.back
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	l.len--
	prevItem := i.Prev
	nextItem := i.Next

	if prevItem == nextItem {
		l.front = nil
		l.back = nil
		return
	}

	if prevItem != nil {
		prevItem.Next = nextItem
	}

	if nextItem != nil {
		nextItem.Prev = prevItem
	}

	if i == l.front {
		l.front = nextItem
	}

	if i == l.back {
		l.back = prevItem
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.front {
		return
	}

	l.Remove(i)
	l.PushFront(i.Value)
}

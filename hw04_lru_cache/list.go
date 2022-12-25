package hw04lrucache

type List[T any] interface {
	Len() int
	Front() *ListItem[T]
	Back() *ListItem[T]
	PushFront(v T) *ListItem[T]
	PushBack(v T) *ListItem[T]
	Remove(i *ListItem[T])
	MoveToFront(i *ListItem[T])
}

type ListItem[T any] struct {
	Value T
	Next  *ListItem[T]
	Prev  *ListItem[T]
	List  *list[T]
}

type list[T any] struct {
	len  int
	head *ListItem[T]
	tail *ListItem[T]
}

func NewList[T any]() List[T] {
	return new(list[T])
}

func (l *list[T]) Len() int {
	return l.len
}

func (l *list[T]) Front() *ListItem[T] {
	return l.head
}

func (l *list[T]) Back() *ListItem[T] {
	return l.tail
}

func (l *list[T]) PushFront(v T) *ListItem[T] {
	item := &ListItem[T]{
		Value: v,
		Next:  l.head,
		Prev:  nil,
		List:  l,
	}

	if l.head == nil {
		l.tail = item
	} else {
		l.head.Prev = item
	}

	l.head = item
	l.len++

	return item
}

func (l *list[T]) PushBack(v T) *ListItem[T] {
	item := &ListItem[T]{
		Value: v,
		Next:  nil,
		Prev:  l.tail,
		List:  l,
	}

	if l.tail == nil {
		l.head = item
	} else {
		l.tail.Next = item
	}

	l.tail = item
	l.len++

	return item
}

func (l *list[T]) Remove(item *ListItem[T]) {
	if item == nil || l.len == 0 {
		return
	}
	if item.List == nil || item.List != l {
		return
	}

	if item == l.head {
		l.head = item.Next
	}
	if item == l.tail {
		l.tail = item.Prev
	}

	if item.Next != nil {
		item.Next.Prev = item.Prev
	}
	if item.Prev != nil {
		item.Prev.Next = item.Next
	}

	item.List = nil
	l.len--
}

func (l *list[T]) MoveToFront(item *ListItem[T]) {
	if item == nil || item == l.head {
		return
	}
	if item.List == nil || item.List != l {
		return
	}

	l.Remove(item)

	if l.head != nil {
		l.head.Prev = item
		item.Next = l.head
	}

	l.head = item

	if l.tail == nil {
		l.tail = item
	}

	item.List = l
	l.len++
}

var _ List[string] = (*list[string])(nil)

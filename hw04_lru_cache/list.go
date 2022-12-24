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

func (l *list[T]) Remove(i *ListItem[T]) {
	if i == nil || l.len == 0 {
		return
	}
	if i.List == nil || i.List != l {
		return
	}

	if i == l.head {
		l.head = i.Next
	}
	if i == l.tail {
		l.tail = i.Prev
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	i.List = nil
	l.len--
}

func (l *list[T]) MoveToFront(i *ListItem[T]) {
	if i == nil || i == l.head {
		return
	}
	if i.List == nil || i.List != l {
		return
	}

	l.Remove(i)

	if l.head != nil {
		l.head.Prev = i
		i.Next = l.head
	}

	l.head = i

	if l.tail == nil {
		l.tail = i
	}

	i.List = l
	l.len++
}

var _ List[string] = (*list[string])(nil)

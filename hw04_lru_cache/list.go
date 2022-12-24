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
	List  *list
}

type list struct {
	len  int
	head *ListItem
	tail *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
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

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
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

func (l *list) Remove(i *ListItem) {
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

func (l *list) MoveToFront(i *ListItem) {
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

var _ List = (*list)(nil)

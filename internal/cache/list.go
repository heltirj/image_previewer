package cache

type addType int

const (
	addTypeFront addType = 1
	addTypeBack  addType = 2
)

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
	length int
	front  *ListItem
	back   *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	return l.addItem(&ListItem{
		Value: v,
	}, addTypeFront)
}

func (l *list) PushBack(v interface{}) *ListItem {
	return l.addItem(&ListItem{
		Value: v,
	}, addTypeBack)
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	if i == l.Back() {
		l.back = i.Prev
	}

	if i == l.Front() {
		l.front = i.Next
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.addItem(i, addTypeFront)
}

func (l *list) addItem(i *ListItem, addType addType) *ListItem {
	if l.length == 0 {
		l.front = i
		l.back = i
		l.length = 1
		return i
	}

	switch addType {
	case addTypeFront:
		i.Next = l.Front()
		l.front.Prev = i
		l.front = i

	case addTypeBack:
		i.Prev = l.Back()
		l.back.Next = i
		l.back = i
	}

	l.length++
	return i
}

package hw04lrucache

import (
	"fmt"
)

type List interface {
	Len() int                  // count of list
	Front() *ListItem          // first element of list
	Back() *ListItem           // last element of list
	PushFront(v any) *ListItem // add value to end position
	PushBack(v any) *ListItem  // add value to start position
	Remove(i *ListItem) bool   // remove element from list, bool return
	MoveToFront(i *ListItem)   // move element to start position
}

type ListItem struct {
	Value any
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	first  *ListItem
	last   *ListItem
	length int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v any) *ListItem {
	item := makeListItem(v)

	if l.first != nil {
		item.Next = l.first
		l.first.Prev = item
	} else {
		l.last = item
	}
	l.first = item
	l.length++
	return item
}

func (l *list) PushBack(v any) *ListItem {
	item := makeListItem(v)

	if l.last != nil {
		item.Prev = l.last
		l.last.Next = item
	} else {
		l.first = item
	}
	l.last = item
	l.length++
	return item
}

func (l *list) Remove(item *ListItem) bool {
	if item == nil {
		fmt.Println("Can't remove item with nil value")
		return false
	}

	switch {
	case item.Next != nil && item.Prev != nil:
		item.Next.Prev, item.Prev.Next = item.Prev, item.Next
	case item.Next != nil && item.Prev == nil:
		item.Next.Prev = item.Prev
		l.first = item.Next
	case item.Next == nil && item.Prev != nil:
		item.Prev.Next = item.Next
		l.last = item.Prev
	case item.Next == nil && item.Prev == nil:
		l.first = item.Next
		l.last = item.Prev
	}

	item.Prev = nil
	item.Next = nil
	l.length--
	return true
}

func (l *list) MoveToFront(item *ListItem) {
	if l.length <= 1 {
		return
	}
	l.Remove(item)
	l.PushFront(item.Value)
}

func makeListItem(val interface{}) *ListItem {
	return &ListItem{
		Value: val,
	}
}

func NewList() List {
	return new(list)
}

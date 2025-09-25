package chdyn

import (
	"fmt"
	"sync"
)

type Node[V any] struct {
	v    V
	next *Node[V]
}

func (n *Node[V]) reset() {
	var v V
	n.v = v

	n.next = nil
}

func (n *Node[V]) Value() V {
	return n.v
}

type List[V any] struct {
	pool sync.Pool
	head *Node[V]
	tail *Node[V]
	size uint32
}

func NewList[V any]() *List[V] {
	return &List[V]{
		pool: sync.Pool{
			New: func() any {
				return &Node[V]{}
			},
		},
		head: nil,
		tail: nil,
		size: 0,
	}
}

func (l *List[V]) getNode() *Node[V] {
	nAny := l.pool.Get()
	n, ok := nAny.(*Node[V])
	if !ok {
		panic(fmt.Errorf("failed to cast node, pool corrupt, nAny=%#v", nAny))
	}
	return n
}

func (l *List[V]) returnNode(n *Node[V]) {
	n.reset()
	l.pool.Put(n)
}

func (l *List[V]) Size() uint32 {
	return l.size
}

func (l *List[V]) Push(v V) {
	n := l.getNode()
	n.v = v
	n.next = nil

	if l.tail == nil {
		l.head = n
		l.tail = n
		l.size = 1
		return
	}

	l.tail.next = n
	l.tail = n
	l.size += 1
}

func (l *List[V]) Head() *Node[V] {
	return l.head
}

func (l *List[V]) Remove() {
	if l.head == nil {
		return
	}

	n := l.head
	defer l.returnNode(n)

	if n.next == nil {
		l.head = nil
		l.tail = nil
		l.size = 0
	} else {
		l.head = n.next
		l.size -= 1
	}
}

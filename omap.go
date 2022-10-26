package omap

import "fmt"

type node[K comparable, V any] struct {
	key   K
	value V
	prev  *node[K, V]
	next  *node[K, V]
}

type OrderedMap[K comparable, V any] struct {
	root  *node[K, V]
	nodes map[K]*node[K, V]
}

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

func newNode[K comparable, V any](key K, value V, prev *node[K, V], next *node[K, V]) *node[K, V] {
	return &node[K, V]{key: key, value: value, prev: prev, next: next}
}

func (n *node[K, V]) entry() *Entry[K, V] {
	return &Entry[K, V]{Key: n.key, Value: n.value}
}

// New constructs an OrderedMap.
func New[K comparable, V any]() *OrderedMap[K, V] {
	root := &node[K, V]{}
	root.prev = root
	root.next = root
	om := &OrderedMap[K, V]{
		root:  root,
		nodes: make(map[K]*node[K, V]),
	}
	return om
}

// Contains checks if a key is in the map.
func (om *OrderedMap[K, V]) Contains(key K) (ok bool) {
	_, ok = om.nodes[key]
	return ok
}

// Get returns the value for a key.
func (om *OrderedMap[K, V]) Get(key K) (value V, ok bool) {
	if n, ok := om.nodes[key]; ok {
		return n.value, true
	}
	return
}

// GetAt returns the key and value at a given index.
func (om *OrderedMap[K, V]) GetAt(index int) (e *Entry[K, V], ok bool) {
	length := om.Len()
	if length == 0 {
		return nil, false
	}

	root := om.root
	idx := index % length
	if idx < 0 {
		idx += length
	}
	curr := root.next
	for i := 0; i < idx; i++ {
		curr = curr.next
		if curr == root {
			curr = curr.next
		}
	}

	return curr.entry(), true
}

// Set sets the value for a key.
func (om *OrderedMap[K, V]) Set(key K, value V) {
	if n, ok := om.nodes[key]; ok {
		n.value = value
		return
	}

	root := om.root
	tail := root.prev
	n := newNode(key, value, tail, root)
	tail.next = n
	root.prev = n
	om.nodes[key] = n
}

// Insert inserts a key and value at a given index.
func (om *OrderedMap[K, V]) Insert(index int, key K, value V) {
	e, ok := om.GetAt(index)
	if !ok || e.Key == key {
		om.Set(key, value)
		return
	}
	om.Delete(key)

	neibor := om.nodes[e.Key]
	n := newNode(key, value, nil, nil)
	if index >= 0 {
		n.prev = neibor.prev
		n.next = neibor
		neibor.prev.next = n
		neibor.prev = n
	} else {
		n.prev = neibor
		n.next = neibor.next
		neibor.next.prev = n
		neibor.next = n
	}
	om.nodes[key] = n
}

// Clear clears the map.
func (om *OrderedMap[K, V]) Clear() {
	for key := range om.nodes {
		delete(om.nodes, key)
	}
}

// Delete deletes the value for a key.
func (om *OrderedMap[K, V]) Delete(key K) (ok bool) {
	n, ok := om.nodes[key]
	if ok {
		n.prev.next = n.next
		n.next.prev = n.prev
		delete(om.nodes, key)
	}
	return ok
}

// DeleteAt deletes the key and value at a given index.
func (om *OrderedMap[K, V]) DeleteAt(index int) (ok bool) {
	e, ok := om.GetAt(index)
	if !ok {
		return false
	}
	return om.Delete(e.Key)
}

// Pop deletes and returns the key and value for a key.
func (om *OrderedMap[K, V]) Pop(key K) (e *Entry[K, V], ok bool) {
	if n, ok := om.nodes[key]; ok {
		return n.entry(), om.Delete(key)
	}
	return
}

// PopAt deletes and returns the key and value at a given index.
func (om *OrderedMap[K, V]) PopAt(index int) (e *Entry[K, V], ok bool) {
	if e, ok := om.GetAt(index); ok {
		return e, om.Delete(e.Key)
	}
	return
}

// Len returns the number of items in the map.
func (om *OrderedMap[K, V]) Len() int {
	return len(om.nodes)
}

func (om *OrderedMap[K, V]) String() string {
	builder := make([]string, len(om.nodes))
	for i, item := range om.Items() {
		builder[i] = fmt.Sprintf("%v:%v", item.Key, item.Value)
	}
	return fmt.Sprintf("OrderedMap%v", builder)
}

// Keys returns the keys in the map.
func (om *OrderedMap[K, V]) Keys() []K {
	keys := make([]K, 0, om.Len())
	for n := om.root.next; n != om.root; n = n.next {
		keys = append(keys, n.key)
	}
	return keys
}

// Values returns the values in the map.
func (om *OrderedMap[K, V]) Values() []V {
	values := make([]V, 0, om.Len())
	for n := om.root.next; n != om.root; n = n.next {
		values = append(values, n.value)
	}
	return values
}

// Items returns the items in the map.
func (om *OrderedMap[K, V]) Items() []*Entry[K, V] {
	items := make([]*Entry[K, V], 0, om.Len())
	for n := om.root.next; n != om.root; n = n.next {
		items = append(items, n.entry())
	}
	return items
}

// IterFunc returns an iterator function over contained items.
func (om *OrderedMap[K, V]) IterFunc() func() (*Entry[K, V], bool) {
	root := om.root
	curr := root.next
	return func() (e *Entry[K, V], ok bool) {
		if curr != root {
			e = curr.entry()
			curr = curr.next
			return e, true
		}
		return
	}
}

// ToHead moves the item to the head of the map.
func (om *OrderedMap[K, V]) ToHead(key K) {
	n, ok := om.nodes[key]
	if !ok {
		return
	}

	n.prev.next = n.next
	n.next.prev = n.prev

	head := om.root.next
	n.prev = om.root
	om.root.next = n
	n.next = head
	head.prev = n
}

// ToTail moves the item to the tail of the map.
func (om *OrderedMap[K, V]) ToTail(key K) {
	n, ok := om.nodes[key]
	if !ok {
		return
	}

	n.prev.next = n.next
	n.next.prev = n.prev

	tail := om.root.prev
	n.prev = tail
	tail.next = n
	n.next = om.root
	om.root.prev = n
}

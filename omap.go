package omap

import "fmt"

type node[K comparable] struct {
	key  K
	prev *node[K]
	next *node[K]
}

type OrderedMap[K comparable, V any] struct {
	root   *node[K]
	nodes  map[K]*node[K]
	values map[K]V
}

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

func New[K comparable, V any]() *OrderedMap[K, V] {
	root := &node[K]{}
	root.prev = root
	root.next = root
	om := &OrderedMap[K, V]{
		root:   root,
		nodes:  make(map[K]*node[K]),
		values: make(map[K]V),
	}
	return om
}

func (om *OrderedMap[K, V]) Contains(key K) bool {
	_, ok := om.values[key]
	return ok
}

func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	v, ok := om.values[key]
	return v, ok
}

func (om *OrderedMap[K, V]) GetAt(index int) (*Entry[K, V], bool) {
	length := om.Len()
	if length == 0 {
		return nil, false
	}

	var curr *node[K]
	root := om.root
	idx := index % length
	if idx < 0 {
		idx += length
	}
	curr = root.next
	for i := 0; i < idx; i++ {
		curr = curr.next
		if curr == root {
			curr = curr.next
		}
	}

	return newEntry(curr.key, om.values[curr.key]), true
}

func (om *OrderedMap[K, V]) Set(key K, value V) {
	if om.Contains(key) {
		om.values[key] = value
		return
	}

	root := om.root
	tail := root.prev
	n := &node[K]{
		key:  key,
		prev: tail,
		next: root,
	}
	tail.next = n
	root.prev = n
	om.nodes[key] = n
	om.values[key] = value
}

func (om *OrderedMap[K, V]) Insert(index int, key K, value V) {
	e, ok := om.GetAt(index)
	if !ok || e.Key == key {
		om.Set(key, value)
		return
	}
	om.Delete(key)

	neibor := om.nodes[e.Key]
	n := &node[K]{key: key}
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
	om.values[key] = value
}

func (om *OrderedMap[K, V]) Clear() {
	for key := range om.nodes {
		delete(om.nodes, key)
		delete(om.values, key)
	}
}

func (om *OrderedMap[K, V]) Delete(key K) bool {
	n, ok := om.nodes[key]
	if ok {
		n.prev.next = n.next
		n.next.prev = n.prev
		delete(om.nodes, key)
		delete(om.values, key)
	}
	return ok
}

func (om *OrderedMap[K, V]) DeleteAt(index int) bool {
	e, ok := om.GetAt(index)
	if !ok {
		return false
	}

	return om.Delete(e.Key)
}

func (om *OrderedMap[K, V]) Pop(key K) (*Entry[K, V], bool) {
	v, ok := om.Get(key)
	if !ok {
		return nil, false
	}

	return newEntry(key, v), om.Delete(key)
}

func (om *OrderedMap[K, V]) PopAt(index int) (*Entry[K, V], bool) {
	e, ok := om.GetAt(index)
	if !ok {
		return nil, false
	}

	return e, om.Delete(e.Key)
}

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

func newEntry[K comparable, V any](key K, value V) *Entry[K, V] {
	return &Entry[K, V]{
		Key:   key,
		Value: value,
	}
}

func (om *OrderedMap[K, V]) Keys() []K {
	keys := make([]K, 0, om.Len())
	for n := om.root.next; n != om.root; n = n.next {
		keys = append(keys, n.key)
	}
	return keys
}

func (om *OrderedMap[K, V]) Values() []V {
	values := make([]V, 0, om.Len())
	for n := om.root.next; n != om.root; n = n.next {
		values = append(values, om.values[n.key])
	}
	return values
}

func (om *OrderedMap[K, V]) Items() []*Entry[K, V] {
	items := make([]*Entry[K, V], 0, om.Len())
	for n := om.root.next; n != om.root; n = n.next {
		items = append(items, newEntry(n.key, om.values[n.key]))
	}
	return items
}

func (om *OrderedMap[K, V]) IterFunc() func() (*Entry[K, V], bool) {
	var curr *node[K]
	root := om.root
	curr = root.next
	return func() (*Entry[K, V], bool) {
		if curr == root {
			return nil, false
		}
		item := newEntry(curr.key, om.values[curr.key])
		curr = curr.next
		return item, true
	}
}

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

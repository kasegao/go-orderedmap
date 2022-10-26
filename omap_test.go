package omap

import (
	"strconv"
	"testing"
)

func makeOMap(length int) *OrderedMap[int, string] {
	om := New[int, string]()
	for i := 0; i < length; i++ {
		om.Set(i, strconv.Itoa(i))
	}
	return om
}

func TestGetSet(t *testing.T) {
	length := 5
	om := makeOMap(length)

	for i := 0; i < length; i++ {
		v, ok := om.Get(i)
		if !ok {
			t.Errorf("Get(%d) failed", i)
		}
		if v != strconv.Itoa(i) {
			t.Errorf("Get(%d) returned %s, want %d", i, v, i)
		}
	}

	if got, ok := om.Get(length); ok {
		t.Errorf("om.Get(%d) = %v, %v; want _, false", length, got, ok)
	}
}

func TestDelete(t *testing.T) {
	length := 5
	om := makeOMap(length)

	index := 3
	om.Delete(index)
	for i := 0; i < length; i++ {
		v, ok := om.Get(i)
		if i == index {
			if ok {
				t.Errorf("Get(%d) succeeded, want failure", i)
			}
			continue
		}

		if !ok {
			t.Errorf("Get(%d) failed", i)
		}
		if v != strconv.Itoa(i) {
			t.Errorf("Get(%d) returned %s, want %d", i, v, i)
		}
	}
}

func TestItems(t *testing.T) {
	length := 5
	om := makeOMap(length)

	items := om.Items()
	for i := 0; i < length; i++ {
		if items[i].Key != i {
			t.Errorf("items[%d].Key = %d, want %d", i, items[i].Key, i)
		}
		if items[i].Value != strconv.Itoa(i) {
			t.Errorf("items[%d].Value = %s, want %d", i, items[i].Value, i)
		}
	}
}

func TestIterFunc(t *testing.T) {
	length := 5
	om := makeOMap(length)

	iter := om.IterFunc()
	for i := 0; i < length; i++ {
		item, ok := iter()
		if !ok {
			t.Errorf("iter() failed")
		}
		if item.Key != i {
			t.Errorf("item.Key = %d, want %d", item.Key, i)
		}
		if item.Value != strconv.Itoa(i) {
			t.Errorf("item.Value = %s, want %d", item.Value, i)
		}
	}
}

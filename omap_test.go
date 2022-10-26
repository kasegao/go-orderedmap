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

func TestClear(t *testing.T) {
	length := 5
	om := makeOMap(length)

	om.Clear()
	if om.Len() != 0 {
		t.Errorf("om.Len() = %d, want 0", om.Len())
	}
}

func TestDeleteAt(t *testing.T) {
	length := 5
	om := makeOMap(length)

	om.DeleteAt(2)
	om.DeleteAt(-1)

	keys := om.Keys()
	if keys[0] != 0 || keys[1] != 1 || keys[2] != 3 {
		t.Errorf("keys = %v, want [0, 1, 3]", keys)
	}
}

func TestInsert(t *testing.T) {
	length := 5
	om := makeOMap(length)

	key := 10
	value := "10"
	om.Insert(1, key, value)

	items := om.Items()
	if i0 := items[0]; i0.Key != 0 || i0.Value != "0" {
		t.Errorf("items[0] = %v, want {0, 0}", i0)
	}
	if i1 := items[1]; i1.Key != key || i1.Value != value {
		t.Errorf("items[1] = %v, want {%d, %s}", i1, key, value)
	}
	if i2 := items[2]; i2.Key != 1 || i2.Value != "1" {
		t.Errorf("items[2] = %v, want {1, 1}", i2)
	}
}

func TestToHeadAndTail(t *testing.T) {
	length := 5
	om := makeOMap(length)

	om.ToHead(3)
	om.ToTail(0)

	keys := om.Keys()
	if keys[0] != 3 || keys[1] != 1 || keys[2] != 2 || keys[3] != 4 || keys[4] != 0 {
		t.Errorf("keys = %v, want [3, 1, 2, 4, 0]", keys)
	}
}

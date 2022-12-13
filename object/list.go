package object

import (
	"bytes"
	"fmt"
	"strings"
)

// List of objects
type List struct {
	// Items holds the list of objects
	Items []Object
}

// Type returns the type of this object.
func (ls *List) Type() Type {
	return LIST
}

// Inspect returns a string-representation of the given object.
func (ls *List) Inspect() string {
	var out bytes.Buffer
	elements := make([]string, 0)
	for _, e := range ls.Items {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

func (ls *List) InvokeMethod(method string, args ...Object) Object {
	switch method {
	case "append":
		if len(args) != 1 {
			return NewError("type error: array.append() expects one argument")
		}
		ls.Append(args[0])
		return ls
	case "clear":
		if len(args) != 0 {
			return NewError("type error: array.clear() expects zero arguments")
		}
		ls.Clear()
		return ls
	case "copy":
		if len(args) != 0 {
			return NewError("type error: array.copy() expects zero arguments")
		}
		return ls.Copy()
	case "count":
		if len(args) != 1 {
			return NewError("type error: array.count() expects one argument")
		}
		return NewInt(ls.Count(args[0]))
	case "extend":
		if len(args) != 1 {
			return NewError("type error: array.extend() expects one argument")
		}
		other, err := AsList(args[0])
		if err != nil {
			return err
		}
		ls.Extend(other)
		return ls
	case "index":
		if len(args) != 1 {
			return NewError("type error: array.index() expects one argument")
		}
		return NewInt(ls.Index(args[0]))
	case "insert":
		if len(args) != 2 {
			return NewError("type error: array.insert() expects two arguments")
		}
		idx, err := AsInteger(args[0])
		if err != nil {
			return err
		}
		ls.Insert(idx, args[1])
		return ls
	case "pop":
		if len(args) != 1 {
			return NewError("type error: array.pop() expects one argument")
		}
		idx, err := AsInteger(args[0])
		if err != nil {
			return err
		}
		return ls.Pop(idx)
	}
	return NewError("type error: %s object has no method %s", ls.Type(), method)
}

// Append adds an element at the end of the list.
func (ls *List) Append(obj Object) {
	ls.Items = append(ls.Items, obj)
}

// Clear removes all the elements from the list.
func (ls *List) Clear() {
	ls.Items = make([]Object, 0)
}

// Copy returns a copy of the list.
func (ls *List) Copy() *List {
	result := &List{Items: make([]Object, 0, len(ls.Items))}
	copy(result.Items, ls.Items)
	return result
}

// Count returns the number of elements with the specified value.
func (ls *List) Count(obj Object) int64 {
	count := int64(0)
	for _, item := range ls.Items {
		if item == obj { // TODO: equality check
			count++
		}
	}
	return count
}

// Extend adds the elements of a list (or any iterable), to the end of the current list.
func (ls *List) Extend(other *List) {
	ls.Items = append(ls.Items, other.Items...)
}

// Index returns the index of the first element with the specified value.
func (ls *List) Index(obj Object) int64 {
	for i, item := range ls.Items {
		if item == obj { // TODO: equality check
			return int64(i)
		}
	}
	return int64(-1)
}

// Insert adds an element at the specified position.
func (ls *List) Insert(index int64, obj Object) {
	ls.Items = append(ls.Items, nil)
	copy(ls.Items[index+1:], ls.Items[index:])
	ls.Items[index] = obj
}

// Pop removes the element at the specified position.
func (ls *List) Pop(index int64) Object {
	if index < 0 || index >= int64(len(ls.Items)) {
		return NewError("index out of range")
	}
	result := ls.Items[index]
	ls.Items = append(ls.Items[:index], ls.Items[index+1:]...)
	return result
}

// Remove removes the first item with the specified value.
func (ls *List) Remove(obj Object) {
	index := ls.Index(obj)
	if index == -1 {
		return
	}
	ls.Items = append(ls.Items[:index], ls.Items[index+1:]...)
}

// Reverse reverses the order of the list.
func (ls *List) Reverse() {
	for i, j := 0, len(ls.Items)-1; i < j; i, j = i+1, j-1 {
		ls.Items[i], ls.Items[j] = ls.Items[j], ls.Items[i]
	}
}

// Sort sorts the list.
func (ls *List) Sort() {
	// TODO
}

func (ls *List) ToInterface() interface{} {
	return "<ARRAY>"
}

func (ls *List) String() string {
	items := make([]string, 0, len(ls.Items))
	for _, item := range ls.Items {
		items = append(items, fmt.Sprintf("%s", item))
	}
	return fmt.Sprintf("List([%s])", strings.Join(items, ", "))
}

func (ls *List) Compare(other Object) (int, error) {
	typeComp := CompareTypes(ls, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherArr := other.(*List)
	if len(ls.Items) > len(otherArr.Items) {
		return 1, nil
	} else if len(ls.Items) < len(otherArr.Items) {
		return -1, nil
	}
	for i := 0; i < len(ls.Items); i++ {
		comparable, ok := ls.Items[i].(Comparable)
		if !ok {
			return 0, fmt.Errorf("type error: %s object is not comparable",
				ls.Items[i].Type())
		}
		comp, err := comparable.Compare(otherArr.Items[i])
		if err != nil {
			return 0, err
		}
		if comp != 0 {
			return comp, nil
		}
	}
	return 0, nil
}

func (ls *List) Reversed() *List {
	result := &List{Items: make([]Object, 0, len(ls.Items))}
	size := len(ls.Items)
	for i := 0; i < size; i++ {
		result.Items = append(result.Items, ls.Items[size-1-i])
	}
	return result
}

func NewList(items []Object) *List {
	return &List{Items: items}
}

func NewStringList(s []string) *List {
	array := &List{Items: make([]Object, 0, len(s))}
	for _, item := range s {
		array.Items = append(array.Items, &String{Value: item})
	}
	return array
}

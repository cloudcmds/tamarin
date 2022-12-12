package object

import (
	"fmt"
)

// Bool wraps bool and implements Object and Hashable interface.
type Bool struct {
	// Value holds the boolean value we wrap.
	Value bool
}

func (b *Bool) Type() Type {
	return BOOL
}

func (b *Bool) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Bool) HashKey() Key {
	var value int64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return Key{Type: b.Type(), IntValue: value}
}

func (b *Bool) InvokeMethod(method string, args ...Object) Object {
	return NewError("type error: %s object has no method %s", b.Type(), method)
}

func (b *Bool) ToInterface() interface{} {
	return b.Value
}

func (b *Bool) String() string {
	return fmt.Sprintf("Bool(%v)", b.Value)
}

func (b *Bool) Compare(other Object) (int, error) {
	typeComp := CompareTypes(b, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherBool := other.(*Bool)
	if b.Value == otherBool.Value {
		return 0, nil
	}
	if b.Value {
		return 1, nil
	}
	return -1, nil
}

func NewBoolean(value bool) *Bool {
	if value {
		return True
	}
	return False
}

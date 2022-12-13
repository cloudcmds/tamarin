package object

import (
	"bytes"
	"fmt"
)

type Module struct {
	Name  string
	Scope interface{}
}

func (m *Module) Type() Type {
	return MODULE
}

func (m *Module) Inspect() string {
	var out bytes.Buffer
	out.WriteString("/")
	out.WriteString(m.Name)
	out.WriteString("/")
	return out.String()
}

func (m *Module) InvokeMethod(method string, args ...Object) Object {
	return NewError("type error: %s object has no method %s", m.Type(), method)
}

func (m *Module) ToInterface() interface{} {
	return "<MODULE>"
}

func (m *Module) String() string {
	return fmt.Sprintf("Module(%s)", m.Name)
}

func (m *Module) Compare(other Object) (int, error) {
	typeComp := CompareTypes(m, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherStr := other.(*Module)
	if m.Name == otherStr.Name {
		return 0, nil
	}
	if m.Name > otherStr.Name {
		return 1, nil
	}
	return -1, nil
}

func (m *Module) Equals(other Object) Object {
	return NewBoolean(other.Type() == MODULE && m.Name == other.(*Module).Name)
}

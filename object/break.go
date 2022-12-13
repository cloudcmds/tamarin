package object

// BreakValue is an implementation detail used to handle break statements
type BreakValue struct{}

func (rv *BreakValue) Type() Type {
	return BREAK_VALUE
}

func (rv *BreakValue) Inspect() string {
	return "BREAK"
}

func (rv *BreakValue) InvokeMethod(method string, args ...Object) Object {
	return nil
}

func (rv *BreakValue) ToInterface() interface{} {
	return "<BREAK_VALUE>"
}

func (rv *BreakValue) Equals(other Object) Object {
	if other.Type() == BREAK_VALUE {
		return True
	}
	return False
}

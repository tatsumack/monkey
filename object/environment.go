package object

type Value struct {
	Obj Object
	IsMutable bool
}

type Environment struct {
	store map[string]*Value
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]*Value)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (*Value, bool) {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		val, ok = e.outer.Get(name)
	}
	return val, ok
}

func (e *Environment) Set(name string, obj Object, isMutable bool) *Value {
	val := Value{Obj:obj, IsMutable: isMutable}
	e.store[name] = &val
	return &val
}

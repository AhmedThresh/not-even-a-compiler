package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: outer,
	}

}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

func (e *Environment) Store(identifier string, value Object) {
	e.store[identifier] = value
}

func (e *Environment) Get(identifier string) (Object, bool) {
	obj, ok := e.store[identifier]
	if e.outer != nil && !ok {
		obj, ok = e.outer.store[identifier]
	}
	return obj, ok
}

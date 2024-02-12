package object

type Environment struct {
	store map[string]Object
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
	return obj, ok
}

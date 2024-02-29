package interpreter

type environment struct {
	enclosing *environment
	values    map[string]any
}

func newEnvironment(enclosing *environment) *environment {
	return &environment{values: make(map[string]any), enclosing: enclosing}
}

func (e *environment) get(name string) (any, bool) {
	value, ok := e.values[name]

	if ok {
		return value, ok
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return value, ok
}

func (e *environment) ancestor(depth int32) *environment {
	env := e
	for i := 0; i < int(depth); i++ {
		env = env.enclosing
	}

	return env
}

func (e *environment) getAt(name string, depth int32) (any, bool) {
	return e.ancestor(depth).get(name)
}

func (e *environment) define(name string, value any) {
	e.values[name] = value
}

func (e *environment) assign(name string, value any) bool {
	if _, ok := e.values[name]; ok {
		e.values[name] = value
		return true
	}

	if e.enclosing != nil {
		return e.enclosing.assign(name, value)
	}

	return false
}

func (e *environment) assignAt(name string, value any, depth int32) bool {
	return e.ancestor(depth).assign(name, value)
}

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

func (e *environment) define(name string, value any) {
	e.values[name] = value
}

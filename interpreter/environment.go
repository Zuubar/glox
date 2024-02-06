package interpreter

type environment struct {
	values map[string]any
}

func newEnvironment() *environment {
	return &environment{values: make(map[string]any)}
}

func (e *environment) lookup(name string) (any, bool) {
	value, ok := e.values[name]
	return value, ok
}

func (e *environment) define(name string, value any) {
	e.values[name] = value
}

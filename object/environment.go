package object

import "fmt"

// Environment is a structure that holds variable mappings.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment creates a new Environment instance.
func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

// Get retrieves an object from the environment.
func (env *Environment) Get(name string) (Object, bool) {
	obj, ok := env.store[name]
	if !ok && env.outer != nil {
		fmt.Printf("Variable not found: %s\n", name)
		return env.outer.Get(name)
	}
	return obj, ok
}

// Set assigns a value to a variable in the environment.
func (env *Environment) Set(name string, val Object) Object {
	fmt.Printf("Setting variable: %s = %s\n", name, val.Inspect())
	env.store[name] = val
	return val
}

// Extend creates a new environment with the current environment as the outer environment.
func (env *Environment) Extend() *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: env,
	}
}

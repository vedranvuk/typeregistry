# typeregistry

Package typeregistry implements a simple type registry that maps go types to, optionally, custom names.

New instances of registered types are later instantiated using reflect by type name.

## Example

```
// Registering a type by name.

func main() {
	type Test struct {
		Name string
	}

	reg := New()
	mytype := &Test{}

	reg.RegisterNamed("foo", mytype)

	intf, err := reg.GetInterface("foo")
	if err != nil {
		panic("fail")
	}
	mynewtype, ok := intf.(*Test)
	if ok {
		mynewtype.Name = "bar"
	}
}
```

Full API consists of following:

```
Register(v interface{}) error
RegisterNamed(name string, v interface{}) error
Unregister(name string) error
GetType(name string) (reflect.Type, error)
GetValue(name string) (reflect.Value, error)
GetInterface(name string) (interface{}, error)
RegisteredNames() []string
```

## License

See included LICENSE file.
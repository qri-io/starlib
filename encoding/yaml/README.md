# yaml
yaml provides functions for working with yaml data.

## Functions

#### `dumps(obj, [obj, ...]) string`
Serialize one or more objects to a yaml string.

If more than one object is provided, the returned string will use
YAML's [Multi-Document](https://yaml.org/spec/1.2.2/#example-two-documents-in-a-stream)
format.

**parameters:**

| name  | type     | description     |
| ----- | -------- | --------------- |
| `obj` | `object` | input object(s) |


#### `loads(source) object`
Read a source yaml string to a Starlark object.

**parameters:**

| name     | type     | description               |
| -------- | -------- | ------------------------- |
| `source` | `string` | input string of yaml data |

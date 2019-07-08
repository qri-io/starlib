# yaml
yaml provides functions for working with yaml data

## Functions

#### `dumps(obj) string`
serialize obj to a yaml string

**parameters:**

| name | type | description |
|------|------|-------------|
| `obj` | `object` | input object |


#### `loads(source) object`
read a source yaml string to a starlark object

**parameters:**

| name | type | description |
|------|------|-------------|
| `source` | `string` | input string of yaml data |



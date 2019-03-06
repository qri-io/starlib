# base64
base64 defines base64 encoding & decoding functions, often used to represent binary as text.

## Functions

#### `decode(src,encoding="standard") string`
parse base64 input, giving back the plain string representation

**parameters:**

| name | type | description |
|------|------|-------------|
| `src` | `string` | source string of base64-encoded text |
| `encoding` | `string` | optional. string to set decoding dialect. allowed values are: standard,standard_raw,url,url_raw |


#### `encode(src,encoding="standard") string`
return the base64 encoding of src

**parameters:**

| name | type | description |
|------|------|-------------|
| `src` | `string` | source string to encode to base64 |
| `encoding` | `string` | optional. string to set encoding dialect. allowed values are: standard,standard_raw,url,url_raw |



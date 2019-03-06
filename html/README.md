# html
html defines a jquery-like html selection & iteration functions for HTML documents

## Functions

#### `html(markup) selection`
parse an html document returing a selection at the root of the document

**parameters:**

| name | type | description |
|------|------|-------------|
| `markup` | `string` | html text to build a document from |



## Types

### `selection`
an HTML document for querying
**Methods**
#### `attr(name) string`
gets the specified attribute's value for the first element in the Selection. To get the value for each element individually, use a looping construct such as each or map method

**parameters:**

| name | type | description |
|------|------|-------------|
| `name` | `string` | attribute name to get the value of |


#### `children() selection`
gets the child elements of each element in the Selection

#### `children_filtered(selector) selection`
gets the child elements of each element in the Selection, filtered by the specified selector

**parameters:**

| name | type | description |
|------|------|-------------|
| `selector` | `string` | a query selector string to filter the current selection, returning a new selection |


#### `contents(selector) selection`
gets the children of each element in the Selection, including text and comment nodes

**parameters:**

| name | type | description |
|------|------|-------------|
| `selector` | `string` | a query selector string to filter the current selection, returning a new selection |


#### `find(selector) selection`
gets the descendants of each element in the current set of matched elements, filtered by a selector

**parameters:**

| name | type | description |
|------|------|-------------|
| `selector` | `string` | a query selector string to filter the current selection, returning a new selection |


#### `filter(selector) selection`
filter reduces the set of matched elements to those that match the selector string

**parameters:**

| name | type | description |
|------|------|-------------|
| `selector` | `string` | a query selector string to filter the current selection, returning a new selection |


#### `get(i) selection`
retrieves the underlying node at the specified index. alias: eq

**parameters:**

| name | type | description |
|------|------|-------------|
| `i` | `int` | numerical index of node to get |


#### `has(selector) selection`
reduces the set of matched elements to those that have a descendant that matches the selector

**parameters:**

| name | type | description |
|------|------|-------------|
| `selector` | `string` | a query selector string to filter the current selection, returning a new selection |


#### `parent(selector) selection`
gets the parent of each element in the Selection

**parameters:**

| name | type | description |
|------|------|-------------|
| `selector` | `string` | a query selector string to filter the current selection, returning a new selection |


#### `parents_until(selector) selection`
gets the ancestors of each element in the Selection, up to but not including the element matched by the selector

**parameters:**

| name | type | description |
|------|------|-------------|
| `selector` | `string` | a query selector string to filter the current selection, returning a new selection |


#### `siblings() selection`
gets the siblings of each element in the Selection

#### `text() string`
gets the combined text contents of each element in the set of matched elements, including descendants

#### `first(selector) selection`
gets the first element of the selection

**parameters:**

| name | type | description |
|------|------|-------------|
| `selector` | `string` | a query selector string to filter the current selection, returning a new selection |


#### `last() selection`
gets the last element of the selection

**parameters:**

| name | type | description |
|------|------|-------------|
| `selector` | `string` | a query selector string to filter the current selection, returning a new selection |


#### `len() int`
returns the number of the nodes in the selection

#### `eq(i) selection`
gets the element at index i of the selection

**parameters:**

| name | type | description |
|------|------|-------------|
| `i` | `int` | numerical index of node to get |



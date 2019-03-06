# zipfile
zipfile reads & parses zip archives

## Functions

#### `ZipFile(data)`
opens an archive for reading


## Types

### `ZipFile`
a zip archive object
**Methods**
#### `namelist() list`
return a list of files in the archive

#### `open(filename string) ZipInfo`
open a file for reading

**parameters:**

| name | type | description |
|------|------|-------------|
| `filename` | `string` | name of the file in the archive to open |


### `ZipInfo`

**Methods**
#### `read() string`
read the file, returning it's string representation


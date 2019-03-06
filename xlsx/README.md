# xlsx
xlsx implements excel file readers in starlark. currently a highly-experimental package that will definitely change at some point in the future

## Functions

#### `get_url(url string)`
fetch an excel file from a url


## Types

### `File`
an excel file
**Methods**
#### `get_sheets() dict`
return a dict of sheets in this excel file

#### `get_rows(sheetname) list`
get all populated rows / columns as a list-of-list strings


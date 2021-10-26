# DataFrame

This package is an implementation of the DataFrame class from the popular [pandas framework](https://pandas.pydata.org/). It represents a 2d columnar data structure that provides many powerful analysis and manipulation tools, similar to a spreadsheet or SQL engine. This implementation aims (eventually) to be a drop-in replacement for the pandas version. Currently it is a long way away as DataFrame has many, many methods. Only a few core methods are working, but they behave nearly the same as the original pandas methods.

## Usage

You can easily run scripts that use `DataFrame` by installing the [qri command-line tool](https://github.com/qri-io/qri) and running the `qri apply` command.

```python
load("dataframe.star", "dataframe")
df = dataframe.Dataframe([["cat", "meow", 1.7],
                          ["dog", "bark", 3.2],
                          ["eel", "zap", 0.6]],
                         columns=["name", "sound", "weight"])
print(df)
```

```
qri apply --file my_script.star
```

## Current Progress

The source file [`dataframe_all_methods.go`](https://github.com/qri-io/starlib/blob/master/dataframe/dataframe_all_methods.go) lists every method in the reference implementation. Any method defined in terms of `methNoImpl` is not yet implemented, but is intended to be. Requests for prioritizing a method implementation can be done by [filing an issue](https://github.com/qri-io/starlib/issues). Methods that are defined in terms of `methMissing` are not planned to be added, because they do not fit within the starlark environment.

## Examples

A number of example usages can be found in the [testdata files](https://github.com/qri-io/starlib/tree/master/dataframe/testdata), along with expected output. Official documentation of each method is in progress.